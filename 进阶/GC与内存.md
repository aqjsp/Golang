# GC与内存

```go
func keep() []byte {
	b := make([]byte, 10<<20) // 10 MiB
	return b[:10]
}

var sink []byte

func main() {
	sink = keep()
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println(m.HeapAlloc) // 仍是数百万字节，不是 10
}
```

没有 `free`。`keep` 返回 10 个字节的切片头，底层数组是 10 MiB。头还活着，数组就还活着。`runtime.GC()` 扫完，这块仍然黑，RSS 钉在那里。C++ 里 `vector::resize(10)` 默认不把 capacity 吐回系统；Python 里 `b = b[:10]` 对 `bytes` 会新对象，对 `bytearray` / `list` 切片仍可能共享或拷贝，取决于类型。Go 的切片头指向整块数组，这是基础篇写过的模型——本篇把它接到 GC：回收的是不可达对象，不是「你已经不用的那一段逻辑长度」。

下面按 1.21+ 的现状讲：三色、**1.8 起的混合写屏障**、pacer、栈与堆、`sync.Pool`。1.3 的全量 STW 标记清除不是现状。对照 C++ RAII 的确定析构、Python 引用计数的立即释放。

---

## 一、没有 free，活着的是可达图

### 1、根到不了的，才是垃圾

GC 从根出发：全局变量、每个 G 的栈、寄存器、把最终化器、把某些 runtime 内部结构。顺着指针走到的对象是活的。走不到的，标记结束之后可以扫掉。没有引用计数，没有析构函数按作用域自动调。

```go
func f() {
	p := &big{buf: make([]byte, 1<<20)}
	_ = p
} // p 没逃出去，函数返回后这块 1 MiB 变成不可达，下一次 GC 可回收
```

`p` 是栈上的指针，指向堆上的 `big`。返回后栈帧没了，没有别的根指着 `big`，它变白（未标记）。不是立刻还给 OS。扫的时候空间回到 Go 的堆空闲链表，可能还留在进程 RSS 里。

C++：`unique_ptr` 出作用域 `delete`，那一页有机会进 `free`，`malloc` 再决定是否 `madvise` 给内核。时机确定，在析构那一行。Python：`del p` 或名字离开，引用计数到 0 立刻跑 `tp_dealloc`，循环引用留给周期检测。Go：函数返回只保证「从这帧出发的边断了」，回收发生在之后某次 GC。

### 2、切片头把整块数组钉住

```go
func readAll() []byte {
	data, _ := os.ReadFile("big.bin") // 可能 80 MiB
	return data[:64]                  // 只要文件头
}
```

`data[:64]` 的 `ptr` 仍指向那 80 MiB 的开头。把返回值存进结构体、缓存、全局，80 MiB 跟到进程结束。基础篇讲过 `s = s[:0]` 不放数组；这里是 GC 视角的同一条： **GC 看指针，不看 len**。

```go
func headerOnly(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	out := make([]byte, 64)
	copy(out, data)
	return out, nil // 大数组没有逃出函数，可被回收
}
```

`copy` 到新小片，旧 `data` 不再被返回值引用。下一次 GC 才能拿回那 80 MiB。不是 `out := data[:64:64]`——那只是 cap 截到 64，**ptr 还在老数组上**，老数组照样活。`slices.Clip` 截 cap，不换数组。要放内存必须让老指针消失。

```go
s := data[:64:64]
s = slices.Clip(s) // cap==len==64，底层仍是 data 那块
```

只有 `copy` 到新 `make`，或让所有指向老数组的头都丢掉。

### 3、字符串、map、接口同样钉

```go
func pinString() string {
	b := make([]byte, 10<<20)
	return string(b[:8]) // 1.22 起编译器可能优化成只拷 8 字节；别赌，大缓冲转短 string 自己 copy
}

func pinMap() {
	m := make(map[int][]byte)
	m[1] = make([]byte, 10<<20)
	delete(m, 1) // 删除键，值不可达；map 结构自己还在，可能仍占着桶
}
```

`delete` 不把 map 缩回去。大 map 用完要 `m = nil` 或换成新 map，整张表才不可达。接口值的 data 指针指向装箱后的动态值，接口变量活着，动态值活着。基础接口篇写过 tab/data；GC 把 data 当一条边。

### 4、RSS ≠ HeapAlloc

```go
var ms runtime.MemStats
runtime.ReadMemStats(&ms)
fmt.Println("heap", ms.HeapAlloc, "sys", ms.HeapSys, "rss 看 OS")
```

`HeapAlloc` 是当前在用的堆对象字节。`HeapIdle` 是 Go 已经扫出来、还没还 OS 的。`HeapReleased` 是已经 `madvise` 出去的。进程 RSS 是内核视角：匿名页还没被要回去就还在。长时间跑的服务，堆在高水位和低水位之间锯，RSS 常贴着高水位。要看真实进程大小用 `/proc/self/statm` 或 `runtime/metrics` 的 `/memory/classes/total:bytes`，不要只看 `HeapAlloc` 掉了就以为 RSS 掉了。

---

## 二、三色标记，1.8 起混合写屏障

### 1、白灰黑

- **白**：还没扫到。标记结束时仍白，就是垃圾。
- **灰**：自己活着，子指针还没扫完。灰是工作队列。
- **黑**：自己和此刻的子指针都处理过。

![三色标记：白未扫，灰队列，黑已处理](../image/gc-tricolor.svg)

标记从根开始，根指向的对象涂灰，弹出灰对象，把它指向的白对象涂灰，自己涂黑。灰队列空，标记结束。并发标记时 mutator（用户 G）同时在改指针，三色不变式会被破坏。破坏的方式有两种经典：

- **强三色**：黑对象不能指向白对象。否则黑已经扫完，白永远不会被看到，误回收。
- **弱三色**：黑指向白也可以，但那条白必须还能从某个灰到达。

插入屏障（Dijkstra）：黑对象被写入一个白指针时，把白涂灰。保证强三色。删除屏障（Yuasa）：指针被覆盖前，把旧指向的对象涂灰。保证「丢掉的那条边」上的对象仍能被扫到，弱三色。

### 2、混合写屏障是现状

Go 1.8 起用 **hybrid write barrier**：覆盖指针时给旧值涂灰（删除），并且在特定条件下给新值涂灰（插入）。写屏障在标记阶段打开，标记结束关掉。伪代码接近 runtime 注释：

```go
// 不是公开 API。slot 是对象里的指针槽，ptr 是新值。
func writePointer(slot *unsafe.Pointer, ptr unsafe.Pointer) {
	if writeBarrierOn {
		shade(*slot) // Yuasa：旧对象变灰
		shade(ptr)   // Dijkstra：新对象变灰
	}
	*slot = ptr
}
```

真实实现还要看「当前 G 的栈扫过没有」、指针是否在堆上。效果是：并发标记期间 mutator 乱改指针，也不会把活对象藏成「只有黑指向白、灰到不了」。栈的处理因此能简化——栈上的指针写入 **不经过** 堆写屏障，1.8 的混合屏障配合「标记开始时把栈算作灰、标记期间栈不再要求实时黑」把 STW 里重扫栈的成本砍掉。细节以 `runtime/mbarrier.go` 注释为准，用户代码只需要记住：

- 现在是并发三色 + 混合写屏障，不是 1.3 停全世界再扫完；
- 写堆上指针在 GC 标记期多一次屏障开销，热路径上疯狂改指针会碰到；
- 栈上局部变量赋值不走这条堆屏障。

```go
type node struct {
	next *node
	id   int
}

func mutate(p *node, q *node) {
	p.next = q // GC 标记期间：写屏障给 p.next 旧值、q 涂灰
}
```

1.5 引入并发三色，用的是插入屏障，标记结束还要 STW 重扫栈。1.8 换成混合屏障，标记终止的 STW 短了一截。之后的版本在 pacer、sweep、最终化器上继续磨，**颜色和屏障这一层从 1.8 沿用至今**。不要把 1.3 的 mark-sweep STW 当现行答案。

### 3、STW 还有，但很短

一次 GC 大约：

1. **Mark start STW**：打开写屏障，扫一部分根，把各个 G 赶到安全点（抢占，GMP 篇 1.14 信号）。
2. **并发标记**：后台 mark worker + mutator assist，用户代码继续跑。
3. **Mark termination STW**：确认灰队列空、写屏障期间的漏网补完，关掉写屏障。
4. **并发扫**：按 span 回收白对象，分配路径上顺便扫（lazy sweep）。

```go
func observe() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	fmt.Println("numGC", ms.NumGC)
	fmt.Println("pauseNs last", ms.PauseNs[(ms.NumGC+255)%256])
	fmt.Println("pauseTotal", ms.PauseTotalNs)
}
```

健康的服务里单次 pause 经常是一百微秒到一两毫秒，看堆大小和根的数量。pause 变到几十毫秒，先查是否 STW 时还有 G 抢占不下来（cgo 里死循环、大量 syscall），再查堆是不是太大、根是不是含一张巨型全局 map。

### 4、最终化器和弱引用不是析构

```go
p := &fileBuf{f: f}
runtime.SetFinalizer(p, func(x *fileBuf) { x.f.Close() })
```

最终化器在对象不可达之后、某个 GC 周期里被单独的 G 调。时机不确定，顺序不确定，对象可能被最终化器救活（再把指针存回去）。不要用它当 C++ 析构：文件句柄、锁、socket 用 `defer Close`。1.24 前后弱引用 API 在演进，生产代码处理资源仍走显式 Close。Python 的 `__del__` 同样不该当资源管理；C++ 析构是唯一确定的那一个。

---

## 三、pacer：什么时候开始扫

### 1、GOGC 是比例，不是绝对阈值

默认 `GOGC=100`：活着的堆是 10 MiB，下一次 GC 目标大约再分配 10 MiB，堆到约 20 MiB 触发。`GOGC=50` 更勤、CPU 多给 GC；`GOGC=200` 堆涨得更猛、GC 次数少。`GOGC=off` 关掉自动 GC，只剩手动 `runtime.GC()`，几乎只用于短命批处理。

```go
debug.SetGCPercent(50) // 运行时改，返回旧值
```

触发条件不是「对象个数」，是 **相对上次标记结束后的存活堆，再分配多少**。短命对象多，分配快，GC 就勤。长命对象堆在那，GOGC=100 会让堆稳定在约 2 倍存活大小附近振荡。

### 2、assist：分配太快就帮着标

并发标记期间，用户 G 还在 `new`。如果分配速度超过 pacer 的估计，mutator 必须先帮着做一点标记工作才允许这次分配。这叫 **mark assist**。症状：分配密集的请求延迟变差，CPU 图上 GC 和业务缠在一起，不是「STW 停了 10ms」那种平台期。

```go
func allocStorm() {
	for i := 0; i < 1_000_000; i++ {
		_ = make([]byte, 256) // 短命对象风暴；pacer 会把 GC 拉起来，assist 吃 CPU
	}
}
```

优化是少分配（逃逸分析篇），不是把 `GOGC` 调到 2000 假装没问题——堆顶上去，下一次标记扫描的根和存活集更大，pause 和 CPU 一起坏。

### 3、GOMEMLIMIT：1.19 起的软上限

```
GOMEMLIMIT=512MiB
```

或 `debug.SetMemoryLimit`。pacer 会在接近上限时更勤地 GC，并更积极地把空闲页还给 OS。这是软限制：极限情况下仍可能超，不是 cgroup OOM 的替代。容器里同时设 cgroup memory 和 `GOMEMLIMIT`（略低于 cgroup），让 Go 自己先收紧，而不是被内核直接杀。

```go
debug.SetMemoryLimit(512 << 20)
```

`GOGC=100` 和 `GOMEMLIMIT` 一起用时，pacer 取更紧的那个目标。只有 limit、把 `GOGC=off`，GC 就主要被内存上限驱动。不要两个旋钮随便拧到极端再怪 runtime。

### 4、什么时候手动 GC

```go
runtime.GC()
debug.FreeOSMemory()
```

`runtime.GC()` 跑完一轮标记+扫。`FreeOSMemory` 尽力把空闲 span 还给 OS，贵，会 STW 相关路径上多干活。服务进程里定期 `FreeOSMemory` 通常是错的：下一波流量又把页要回来，抖动 RSS。批处理：算完一大波、接下来要 fork 或要给旁边进程腾内存，可以调一次。观测用 `runtime/metrics`：

```go
func heapLive() uint64 {
	samples := []metric.Sample{{Name: "/gc/heap/live:bytes"}}
	metric.Read(samples)
	return samples[0].Value.Uint64()
}
```

`MemStats` 要停一下取快照，热路径别每次请求 `ReadMemStats`。metrics 包更轻。

---

## 四、栈上还是堆上

### 1、编译器决定，不是你写 new

Go 没有 C++ 那种「`new` 在堆、局部对象在栈」的语法合同。`new(T)`、取地址 `&T{}`、`make` 出来的东西，都可以在逃逸分析证明没逃出函数时放栈。反过来，局部变量被闭包、被 `go`、被接口装箱、被返回指针，上堆。

```go
func onStack() int {
	x := 1
	return x // x 不逃逸
}

func onHeap() *int {
	x := 1
	return &x // x 逃逸到堆
}
```

GC 只扫堆（和栈上当作根的那一截指针）。栈帧随着 G 的调用生灭，帧没了，里面的值直接丢，不必标色。这是为什么逃逸分析值钱：对象留栈，分配是 SP 一减，回收是 SP 一加，GC 看不见它。

C++ 局部对象确定在栈，`new` 确定在堆。Python 对象几乎都在堆，栈上是指向 PyObject 的指针。Go 夹在中间：语法看起来都像局部变量，位置由编译器说了算。下一篇专门钉逃逸规则；这里只接 GC： **堆上的才进三色图**。

### 2、每个 P 一份 mcache

分配小对象不走全局锁。每个 P 有 mcache，里面按 size class 挂着当前的 span。`make([]byte, 17)` 对齐到某个 class（比如 24 字节），从当前 P 的 mcache 抠一块。mcache 用完向 mcentral 要，再没有向 mheap 要页。所以：

- 分配路径依赖 P，和 GMP 绑在一起；
- 没有 P 的 M（syscall 里）不能直接分配 Go 对象；
- size class 有内部碎片：17 字节对象占 24。

```go
func small() {
	a := make([]byte, 17)
	b := make([]byte, 24)
	_ = a
	_ = b
}
```

大对象（超过 32 KiB 量级，随版本）直接从 mheap 按页，不进 mcache 的小 class。大对象分配更爱锁，也更容易把堆顶抬起来。`make([]byte, 10<<20)` 就是这条路。

### 3、span 扫完才真正能复用

标记结束对象还占着 span。sweep 把死亡对象的位图清掉，span 重新挂回空闲。并发 sweep 意味着 GC「结束了」之后一小段，分配仍可能先帮着扫。`HeapReleased` 上涨发生在 sweep 之后 runtime 决定把整页还给 OS 时。所以「GC 刚跑完 RSS 立刻掉」经常不发生。

### 4、栈本身也是内存

每个 G 初始约 2 KiB 栈，可增长。十万个 runnable 之外的 G，十万份小栈。栈上的指针是 GC 的根，G 多，扫栈的工作也多。泄漏 G（channel 篇）既泄漏栈，也让每次 GC 多扫一份。`runtime.NumGoroutine` 涨和 GC CPU 涨经常一起出现。

---

## 五、sync.Pool：不是自己的堆，是抗 GC 的暂存

### 1、Get / Put

```go
var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 4096)
		return &b
	},
}

func handle() {
	bp := bufPool.Get().(*[]byte)
	b := (*bp)[:0]
	b = append(b, ...)
	*bp = b
	bufPool.Put(bp)
}
```

`Get` 从当前 P 的本地池拿，没有就 `New()`。`Put` 放回本地池。池里的对象 **每次 GC 会被清掉**：Pool 不保证 Put 进去的东西下次还在。它是「抗住两次 GC 之间的分配风暴」，不是缓存。把数据库连接、必须关掉的句柄放进 Pool，GC 一来对象被丢掉且没 Close，泄漏。Pool 只放无状态、可重置、丢了能再 `make` 的。

### 2、和 GMP 的关系

Pool 的私有池按 P 分。G 从 P0 的 Pool.Get，Put 时可能已经在 P3，对象进了 P3 的池。空的 P 会从别的 P 偷，和 work stealing 同一思想。`GOMAXPROCS` 很大时，对象散在很多 P 上，局部性差一些。不要为了 Pool 去调 `GOMAXPROCS`。

### 3、Put 之前清引用

```go
type req struct {
	Body []byte
	next *req
}

func put(p *sync.Pool, r *req) {
	r.Body = nil // 否则 Pool 活着，10 MiB Body 被钉到下一次 GC 清池
	r.next = nil
	p.Put(r)
}
```

Pool 清是在 GC 时。两次 GC 之间，Put 进去的大切片照样占 HeapAlloc。重置再 Put，是用 Pool 时必做的一步。`fmt`、`encoding/json` 内部大量 Pool，自己写中间件拷那些模式时把重置带上。

### 4、什么时候不要 Pool

- 对象很大、很少分配：Pool 的命中率低，还干扰 GC。
- 对象持有必须释放的资源。
- 大小每次都不同，拿出来还要 `make` 到新长度：收益被吃光。
- 还没 profile 就先上 Pool：多半是在藏分配，而不是减少分配。先看 pprof `allocs`，确认是这块 `make` 热，再 Pool。

```go
var p sync.Pool

func no() {
	p.Put(make([]byte, 10<<20)) // 下一次 GC 前这 10 MiB 一直在
}
```

---

## 六、对照 C++ RAII、Python RC

### 1、C++：寿命钉在类型上

```cpp
{
    std::vector<char> v(10 << 20);
    v.resize(10);                 // capacity 仍是 10 MiB
}                                 // 析构，free 整块
```

`resize` 不缩 capacity，这一点和 Go 切片一样。不一样的是 `}` 那一行：析构确定发生，`operator delete` 确定调。Go 的 `sink = b[:10]` 把寿命交给可达图，没有 `}` 能把数组弄死，除非 `sink` 不再指向它。C++ 也可以钉住：把 `vector` 移到全局，同样 10 MiB 不走。RAII 保证的是 **最后一个所有者没了就跑析构**，不是保证你不会把所有者藏起来。

`unique_ptr` / `shared_ptr` 是所有权写在类型里。Go 没有。活着 = 有指针从根够得到。环状引用在 Go 里不是泄漏（三色能处理环），在 `shared_ptr` 里是泄漏，要 `weak_ptr`。Go 的泄漏是 **还以为自己不用了，其实有个全局 / 长寿命切片头 / 没关的 G 抓着**。

### 2、Python：计数到 0 立刻走

```python
def keep():
    b = bytearray(10 * 1024 * 1024)
    return b[:10]  # bytearray 切片是新对象，拷 10 字节；原 b 计数归零，立刻释放
```

`bytes` / `bytearray` 切片常拷贝，`list` 切片浅拷贝一份指针表。Python 程序员较少踩「子切片钉死整块 buffer」——那是 `memoryview` 或某些 C 扩展 buffer 的坑。Go 切片默认共享底层数组，这条坑每天都有。

Python 循环引用：`a.x = a`，计数永远 ≥1，靠 `gc` 模块的分代周期检测。Go 的环：A 指 B、B 指 A，根到不了它们，三色结束后都是白，一起收。不必弱引用破环。Go 里 `__del__` 式最终化器反而麻烦，因为要处理「复活」。

GIL 让 CPython 的 RC 增减好做；无 GIL 的解释器才把 RC 变成原子或改走别的 GC。Go 从一开始就按无 GIL、多 P 并发 mutator 来设计写屏障。

### 3、三种语言各会把 RSS 钉在哪

- C++：自己没 `delete`、`vector` 容量、内存池、fragmentation、仍活着的 `shared_ptr` 环。
- Python：全局 dict 缓存、`lru_cache`、循环引用加 `__del__`、C 扩展 `malloc` 不走 Python 堆。
- Go：切片/string/map 底层、全局变量、`sync.Pool` 两次 GC 之间、泄漏的 G、`pprof` 文件没关、长寿命 cache。

没有一种「有 GC 就不用管内存」。Go 管的是 **不用写 free**，不管 **逻辑上已经瘦下来的视图还指着胖子**。

---

## 七、写代码时怎么少跟 GC 打架

### 1、能预分配就预分配

```go
out := make([]int, 0, len(in))
for _, v := range in {
	out = append(out, v*2)
}
```

反复 `append` 从 0 长，中间那些被换掉的底层数组都是短命垃圾，GC 要扫要收。一次 `make` 够 cap，中间零次。下一篇从逃逸和 pprof 再钉一遍；这里只说对 GC 的影响：分配次数降下去，pacer 少醒。

### 2、别让短命对象进接口

```go
func f(v any) {}

for i := 0; i < 1e6; i++ {
	f(i) // int 装箱上堆，一百万小对象，GC 的活
}
```

接口装箱是分配源。热循环里 `fmt.Sprintf("%d", i)` 同样。具体类型参数、或自己 `strconv.Itoa`。

### 3、大缓冲复用要自己盯寿命

```go
var buf []byte

func handler(w http.ResponseWriter, r *http.Request) {
	buf = buf[:0]
	buf = append(buf, ...) // 多个请求并发时这是 race；每个 G 自己的 buf 或 Pool
	w.Write(buf)
}
```

全局 `buf` 复用：并发 race，而且一旦某次请求把 cap 撑到 32 MiB，后面所有请求那条切片头都钉着 32 MiB。Pool 里 Put 前若不清 cap，同样。常见写法：

```go
if cap(b) > 64<<10 {
	return // 太大的不放回 Pool，让 GC 收
}
p.Put(&b)
```

### 4、指针多的结构体，GC 扫得贵

```go
type A struct {
	p *int
	q *int
	r *int
}

type B struct {
	p, q, r int
}
```

堆上一百万个 `A`，每个三个指针，标记阶段要跟三条边。`B` 没有指针，span 可以标成 noscan，GC 扫到这块几乎跳过。热点数据能做成无指针的就做（`int` id 代替 `*T`，数组代替链表）。这是 GC CPU 的结构性原因，比调 `GOGC` 更本。

---

## 八、观测

### 1、`GODEBUG=gctrace=1`

```
GODEBUG=gctrace=1 ./app
```

每轮一行：墙钟、CPU 占用、堆前后大小、是否 assist。看「是不是每 50ms 一轮」和「heap_goal 是不是在涨」。

### 2、pprof heap

```go
import _ "net/http/pprof"
```

`/debug/pprof/heap` 看当前活着的对象按分配栈聚合。`inuse_space` 是此刻占的，`alloc_space` 是累计分配。RSS 钉死先看 `inuse_space`：谁还抓着那块大数组。短命风暴看 `allocs` profile。下一篇把命令行跑法写全。

### 3、`runtime.ReadMemStats` 字段

- `HeapAlloc`：在用对象；
- `HeapIdle`：空闲、未还 OS；
- `HeapReleased`：已还 OS；
- `NextGC`：下一次触发的堆目标；
- `NumGC`、`PauseTotalNs`；
- `GCCPUFraction`：进程历史上 GC 占的 CPU 比例。

`GCCPUFraction` 是累计值，服务刚启动那几次 GC 会把它拉得很难看，跑稳了再读。

---

## 九、一段可跑的对照

```go
package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

func pin() []byte {
	b := make([]byte, 4<<20) // 4 MiB
	for i := range b {
		b[i] = 1
	}
	return b[:16]
}

func copyOut() []byte {
	b := make([]byte, 4<<20)
	for i := range b {
		b[i] = 1
	}
	out := make([]byte, 16)
	copy(out, b)
	return out
}

func heapStats() (alloc, sys uint64) {
	runtime.GC()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ms.HeapAlloc, ms.HeapSys
}

var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 1024)
		return &b
	},
}

func withPool(n int) {
	for i := 0; i < n; i++ {
		bp := bufPool.Get().(*[]byte)
		b := (*bp)[:0]
		for k := 0; k < 64; k++ {
			b = append(b, byte(k))
		}
		*bp = b
		bufPool.Put(bp)
	}
}

func withoutPool(n int) {
	for i := 0; i < n; i++ {
		b := make([]byte, 0, 1024)
		for k := 0; k < 64; k++ {
			b = append(b, byte(k))
		}
		_ = b
	}
}

func main() {
	debug.SetGCPercent(100)

	sink := pin()
	a1, _ := heapStats()
	fmt.Println("pinned heapAlloc", a1, "len", len(sink))

	sink = copyOut()
	a2, _ := heapStats()
	fmt.Println("copied heapAlloc", a2, "len", len(sink))

	sink = nil
	a3, _ := heapStats()
	fmt.Println("released heapAlloc", a3)

	const n = 20000
	start := time.Now()
	withoutPool(n)
	t1 := time.Since(start)
	runtime.GC()
	start = time.Now()
	withPool(n)
	t2 := time.Since(start)
	fmt.Println("alloc loop", t1, "pool loop", t2)

	p := &struct{ n int }{7}
	runtime.SetFinalizer(p, func(x *struct{ n int }) {
		fmt.Println("finalizer", x.n)
	})
	p = nil
	runtime.GC()
	time.Sleep(20 * time.Millisecond)
	fmt.Println("GOGC", debug.SetGCPercent(-1))
	debug.SetGCPercent(100)
}
```

跑完对照：`pin` 之后 `HeapAlloc` 仍是数百万；`copyOut` 之后掉到接近小对象水平；`sink=nil` 再 GC 更低。Pool 那圈分配次数少，耗时通常低于每次 `make`（幅度看机器）。最终化器可能打出 `finalizer 7`，也可能这次调度还没轮到——这就是「不能当析构」的现场。把 `return b[:16]` 改成长期放在全局 `var cache []byte = pin()`，RSS 会一直贴着那 4 MiB。

---

## 十、清单

没有 `free`。活着的是根够得到的对象。切片 len 缩短不砍底层数组，`Clip` 也不换数组，要放内存就 `copy` 到新片或丢掉所有头。三色白灰黑；1.8 起混合写屏障（删+插），并发标记，STW 短，不是 1.3 全停。pacer 用 `GOGC` 比例和下次堆目标，分配过快 assist；1.19+ `GOMEMLIMIT` 是软上限。对象在栈还是堆由逃逸分析决定，GC 扫堆和栈根。`sync.Pool` 按 P 分池，每次 GC 清空，只放可丢可重置的无资源对象，Put 前把大引用抹掉。C++ 析构确定发生；Python RC 到 0 立刻走、环靠周期检测；Go 环能收，钉死 RSS 的是你还握着的指针。观测用 `gctrace`、heap profile、`HeapAlloc`/`HeapReleased` 分开看。

下一篇：下游 HTTP 已经超时，父请求还在等。context 取消树，和 `error` 的 `%w`。
