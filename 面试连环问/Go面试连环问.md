# Go面试连环问

面试官丢一句「slice 和 array 有什么区别」，要的不是清单，是你能不能被往下追：header 三个字、append 何时换数组、旧 header 还指不指旧块。下面每题按能说出口的骨架写，追问带代码，对照基础 / 进阶篇。默认 Go 1.21+。

---

### 1、slice 底层是什么？扩容怎么做？

能说出口：切片值是 header：指针、len、cap，64 位 24 字节。底层是数组。赋值拷 header，数组共享。`append` 先看 cap，够就原地写、返回 len+1 的新头；不够 runtime 另起一块，把旧 **len** 个元素拷过去，旧数组若还有 header 指着就继续活。

往下追：`b := a; b[0]=9`，`a[0]` 也是 9。再 `b=append(b,4)` 且 cap 满，之后 `b[0]=9` 改不到 `a`。函数参数同样拷头：改下标外面看见，append 换数组外面看不见，必须接返回值。

```go
a := []int{1, 2, 3} // len=3 cap=3
b := a
b[0] = 9
fmt.Println(a[0]) // 9
b = append(b, 4)  // 新数组
b[0] = 1
fmt.Println(a[0], b) // 9 [1 2 3 4]
```

再往下：1.18 起不是「<1024 翻倍否则 +25%」。大致：目标已大于 2× 旧 cap 就涨到目标；旧 cap<256 翻倍；再往上每次加 `(oldcap+3*256)/4`，从 2× 滑向约 1.25×。具体随版本和元素大小类微调，**不要断言 cap==8**。生产 `make([]T,0,n)` 估够。三索引 `s[i:j:k]` 把 cap 卡住，挡住 `append` 踩原数组后面的元素。

```go
s := make([]int, 0, 1)
for i := 0; i < 8; i++ {
	s = append(s, i)
	fmt.Println(len(s), cap(s))
}
arr := [4]int{1, 2, 3, 4}
u := arr[0:1:1]
u = append(u, 9) // 新数组，arr[1] 仍是 2
```

对照：切片篇 header 三字；C++ `vector` 赋值深拷、扩容后 iterator 失效（UB）；Go 旧头仍合法，只是指旧块，逻辑 bug 不是悬空。Python `b=a` 永远共享同一个 list。

---

### 2、map 为什么不能并发读写？

能说出口：map 变量里是 `*hmap`。赋值拷指针，表共享。并发读写 runtime **尽力**检测，命中就 `fatal error: concurrent map read and map write`，不是 panic，`recover` 接不住，进程死。没有这层检测时仍是 data race，结果未定义。保护：`sync.Mutex` / `RWMutex`，或特定场景 `sync.Map`。

往下追：nil map 读零值、写 panic。`make` 之后才能插。元素不能取地址：增长会搬桶，取到的指针会悬。range 的键序每次运行不必相同。函数里插键外面看见（指针盒子没换）；`m=make(...)` 换指针传不回去。

```go
m := map[int]int{}
go func() {
	for { m[1] = 1 }
}()
for {
	_ = m[1] // 很快 fatal
}
```

再往下：`sync.Map` 适合键稳定、写少读多、互不相交的键；不是「并发 map 的默认替换」。计数、缓存条目用 mutex 包普通 map 更直接。不要用 channel 把所有 map 操作扔进一个 G 当「线程安全」——那是把锁换成串行 G，多一次调度。

```go
type Safe struct {
	mu sync.Mutex
	m  map[string]int
}

func (s *Safe) Inc(k string) {
	s.mu.Lock()
	s.m[k]++
	s.mu.Unlock()
}
```

对照：切片篇「并发写就地处死」；C++ `unordered_map` 并发是 UB，不会专门 fatal；Python dict 在 GIL 下单线程字节码，无 GIL 解释器同样要锁。

---

### 3、GMP 是什么？goroutine 为什么能开十万个？

能说出口：G 是用户态执行体（栈、PC、状态），初始栈约 2KiB，不够连续栈拷大。M 是 OS 线程。P 是逻辑处理器，本地队列、mcache。跑用户 Go 代码必须 M+P+G。`GOMAXPROCS` 限制的是 **P 的个数**，不是 M，也不是 G。十万 G 争的是几个 P，大多数等在 channel / netpoll，不占 M。

往下追：`go f()` 不创建 OS 线程，把 G 扔进当前 P 的本地队列或 `runnext`。本地队列 256，满了一半踢全局。空闲 P work stealing。阻塞 syscall 把 M 带下去，P 被摘走，另拉 M 接手 P。`time.Sleep`、channel、网络 `Read` 是 G 阻塞，M+P 去跑别人。

```go
func main() {
	fmt.Println(runtime.GOMAXPROCS(0)) // 默认逻辑 CPU 数
	var wg sync.WaitGroup
	wg.Add(100_000)
	for i := 0; i < 100_000; i++ {
		go func() { defer wg.Done(); time.Sleep(time.Second) }()
	}
	wg.Wait()
}
```

再往下：1.14 起信号抢占，不再只靠函数调用协作让出。sysmon 盯长时间运行的 G。cgo / 阻塞 syscall 期间 C 代码不受 Go 抢占。`LockOSThread` 把 G 钉在 M 上，贵。1.0「只在调用处切换」不是现状。网络 `Read` 走 netpoll：EAGAIN 则 G 等待、M+P 跑别人；磁盘普通文件 `Read` 往往把 M 带下去。不要说「Go 的 I/O 都不占线程」。

```go
ln, _ := net.Listen("tcp", "127.0.0.1:0")
for {
	c, _ := ln.Accept()
	go handle(c) // 一连接一 G，不是一连接一线程
}
```

对照：GMP 篇开篇十万 Sleep；C++ `std::thread` 1:1，十万条栈先把机器打穿；Python `threading` 同样 1:1 加 GIL，`asyncio` 能堆协程但不在多核上并行算。

---

### 4、channel 关闭的规则是什么？

能说出口：发送方 close，只关一次。关了再发、再关：panic。收关闭且空的 channel：立刻得零值，`ok==false`。关了但格子里还有值：先把值收完，`ok` 仍 true。`range ch` 直到 close 且空；不关，range 永远阻塞。多个发送者不要每人 close 一次。

往下追：close 语义是「不会再有新值」，可当广播「开始/停」。不要用 close 传业务数据。从已关闭 channel 收在 select 里永远就绪，会热循环吃零值——关了就 `ch=nil` 把这一支从 select 摘掉。

```go
ch := make(chan int, 1)
ch <- 1
close(ch)
v, ok := <-ch // 1 true
v, ok = <-ch  // 0 false
// ch <- 2    // panic
// close(ch)  // panic
```

再往下：nil channel 收发永不就绪。fan-in 里一支关了赋 nil，另一支继续。函数不要返回 nil channel 表示失败——调用方 `range` 会永远卡；失败用 error。零值陷阱：`chan *T` 关了收到 nil 指针，当活对象用就 panic。

```go
ch := make(chan *User)
close(ch)
u := <-ch
_ = u.Name // panic
```

对照：channel 篇会合 vs 格子；C++ 没有语言级 channel，`queue`+mutex+cv 要自己定「结束」协议；Python `queue.Queue` 无界，结束靠哨兵对象。

---

### 5、goroutine 泄漏怎么发生？怎么查？

能说出口：G 在 channel 发送无人收、`range` 无人 close、`select` 不盯 `ctx.Done()` 上永远阻塞。runtime 死锁检测只在 **没有任何 G 能跑** 时开火；main 还在服务时，泄漏的 G 检测器沉默。信号是 `runtime.NumGoroutine` 单调涨、`pprof` goroutine profile 里同一栈堆积。

往下追：HTTP handler 里 `go f()` 用 `r.Context()`，请求结束 ctx cancel，若 `f` 不看 Done 仍可能卡在别处；用 `Background()` 则请求取消完全灌不进去。无缓冲往已走掉的接收者发，发送者卡死。`time.After` 在循环 select 里造大量 timer（1.23 前更脏）。

```go
func leak(ch chan int) {
	go func() { ch <- 1 }() // 没人收，这个 G 卡到进程死
}

func leakRange() {
	ch := make(chan int)
	go func() {
		ch <- 1
		// 忘了 close
	}()
	for range ch { // 收到 1 之后永远阻塞
	}
}
```

再往下：出口必须是 `select { case ch<-x: case <-ctx.Done(): }`。`errgroup` 把取消和 Wait 焊在一起。pprof：`import _ "net/http/pprof"`，看 `goroutine?debug=2`。HTTP 实战篇「开了就不管」对照：`NumGoroutine` 不回基线。

```go
func send(ctx context.Context, ch chan<- int, v int) error {
	select {
	case ch <- v:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
```

对照：channel 篇第七节；C++ detach 的线程泄漏是线程表和栈；Python asyncio 忘记 await 的 Task 是警告。Go 的 G 便宜，泄漏可以堆到十万才被 RSS 发现。

---

### 6、GC 三色标记现在怎么工作？

能说出口：白未标记，灰待扫，黑已扫完子节点。从根（全局、各 G 栈、寄存器）涂灰，弹出灰，子白变灰，自己变黑。灰空，白是垃圾。1.5 起并发标记；1.8 起 **混合写屏障**（删+插），不是 1.3 全量 STW mark-sweep。STW 还有两小段：mark start 开屏障、mark termination 关屏障。

往下追：并发时 mutator 改指针会破坏三色。强三色：黑不能指白。弱三色：黑可以指白，但那白还能从灰到。混合屏障覆盖指针时旧值涂灰，特定条件新值也涂灰。栈写入不走堆屏障，靠标记开始把栈当灰。写屏障只在标记期开。

```go
type N struct{ next *N }

func f(p, q *N) {
	p.next = q // 标记期间：写屏障给旧值、q 涂灰
}
```

再往下：pacer 用 `GOGC`（默认 100，堆相对存活翻倍触发）和 1.19+ `GOMEMLIMIT` 软上限。分配太快 mutator assist：用户 G 自己帮着标，延迟换 CPU。环能收（相对 `shared_ptr` 要 `weak_ptr`）。泄漏是你还握着指针：大切片头、没关的 G、全局 cache、`pprof` 文件没关。`GODEBUG=gctrace=1` 看 pause 和 heap 目标。`HeapAlloc` 是 Go 以为在用的，`HeapReleased` 才是还给 OS 的；RSS 高不一定是泄漏。

```go
func keep() []byte {
	b := make([]byte, 10<<20)
	return b[:10] // 头还活着，10 MiB 数组还黑
}
```

对照：GC 篇开篇这块 10 MiB；C++ RAII 确定析构，没有三色；Python RC 到 0 立刻走，环靠分代周期检测。不要把 1.3 的 STW 当现行答案。

---

### 7、什么会逃逸到堆？怎么看见？

能说出口：编译器证明对象不在函数返回后被引用，可以留栈；否则上堆。常见：返回指针、闭包捕获、赋进接口、`go` 起来的函数用、取地址存到外面、方法值绑接收者、太大的 `make`。GC 只扫堆和栈根；留栈是 SP 加减，三色看不见。

往下追：`go build -gcflags=-m`。`moved to heap`、`leaking param`。内联改变逃逸，`-l` 关内联看 API 边界。热循环 `fmt.Sprintf` / `...any` 把 `int` 装箱，pprof allocs 里 `convT2E`。经验随版本变，以你模块的 `-m` 为准。

```go
func ret() *int { x := 1; return &x } // x 逃逸

func box(id int) string {
	return fmt.Sprintf("id=%d", id) // id 进 any，典型逃逸
}

func better(id int) string { return strconv.Itoa(id) }
```

再往下：先 `-benchmem` 再改。切片一次 `make` 够 cap。`strconv.AppendInt` 写已有 buf。不要先上 `sync.Pool` 藏分配。`//go:noescape` 是给 runtime 的，用户代码标错就是悬空。闭包和 `go` 让捕获变量上堆；方法值 `t.Inc` 把接收者藏进函数值。`leaking param: b` 会让调用方的实参跟着逃——库 API 的逃逸写进文档，比自己猜经验强。

```go
func sink(f func()) {}

func hot(t T) {
	sink(t.Inc) // 方法值，t 逃逸
}
```

对照：逃逸篇开篇 Sprintf；C++ 局部默认栈、`new` 才堆；Python 对象默认堆。Go 看起来是局部变量，位置由编译器说了算。

---

### 8、为什么 `var err error = (*T)(nil)` 不是 nil？

能说出口：接口值两个字：tab（类型+方法表）和 data。真 nil 是两个零。把 typed nil 指针赋进去，tab 有类型，data 是 0，`err==nil` 为 false。`if err != nil` 走进去，`err.Error()` 可能对 nil 接收者 panic。打印却常是 `<nil>`，不要用打印判断。

往下追：函数内部 `var e *PathError; return e` 在成功路径返回的是 typed nil。内部变量就用 `error`，失败 `return &PathError{...}`，成功 `return nil`。

```go
func load(bad bool) error {
	var err *os.PathError
	if bad {
		err = &os.PathError{Op: "open", Path: "x", Err: os.ErrNotExist}
	}
	return err // bad==false：不是 nil error
}

func main() {
	err := load(false)
	fmt.Println(err == nil) // false
}
```

再往下：`fmt.Printf("%T %#v", err, err)`。反射：真 nil 的 `error` `ValueOf` 是 Invalid；typed nil 是 Ptr 且 `IsNil`。接口比较：两个接口 `==` 要类型和值都等；一个真 nil 一个 typed nil 不等。这是 Go 最贵的那类 bug 之一。

对照：接口篇第三节；C++ 空指针就是空，`if(p)` 一把梭，没有「带类型的空」这层；Python `None` 是单例对象，`is None` 没有 typed None。

---

### 9、defer 的参数何时求值？和返回值怎么互相看见？

能说出口：执行到 `defer f(args)` 这一行，**args 立刻求完**存进记录；`f` 等到函数返回（含 panic）LIFO 调。不是离开 `{}`。命名返回值是函数入口的局部变量，`return 1` 先写入该变量，再跑 defer，defer 闭包能改到调用方拿到的值。

往下追：`defer fmt.Println(x)` 打印登记时的 x；`defer func(){ fmt.Println(x) }()` 打印返回时的 x。循环里 `defer f.Close()` 堆到函数结束才关，抽一层函数。`os.Exit` 不跑 defer。`recover` 必须写在 defer 的函数里。

```go
func deferArg() {
	x := 1
	defer fmt.Println("arg", x) // 1
	defer func() { fmt.Println("cls", x) }()
	x = 2
}

func inc() (n int) {
	defer func() { n++ }()
	return 1 // 返回 2
}
```

再往下：匿名返回值没有名字，defer 改局部不影响已拷进返回槽的值。循环百万次不要无意义 defer；`defer Unlock` 在函数级粒度，临界区要短就抽函数。现代编译器多数 defer 是栈记录，热路径 Unlock 可以写。

对照：函数篇 defer 节；C++ RAII 块结束析构；Python `with` 块结束。Go 按函数边界。

---

### 10、error wrapping 怎么做？`Is` / `As` 认什么？

能说出口：`error` 是接口，值，不是异常。`fmt.Errorf("store get %d: %w", id, err)` 包一层，`Unwrap` 链还在。`errors.Is(err, sentinel)` 沿链比 `==`（或自己的 `Is` 方法）。`errors.As` 沿链找具体类型。`%v` 只字符串，链断。不要 `err.Error()=="timeout"`。

往下追：sentinel：`var ErrNotFound = errors.New("not found")`。自定义类型实现 `Error()`，要被 `Is` 认可以写 `Is(error) bool`。`errors.Join` 多支，`Is` 对其中任一支为真。`panic` 不是业务失败。

```go
var ErrNotFound = errors.New("not found")

err := fmt.Errorf("user %d: %w", 7, ErrNotFound)
fmt.Println(errors.Is(err, ErrNotFound)) // true

err = fmt.Errorf("user %d: %v", 7, ErrNotFound)
fmt.Println(errors.Is(err, ErrNotFound)) // false
```

再往下：HTTP 边界 `Is` 成状态码，中间层只 wrap。`Canceled` 常是客户端走了，不是 500。`DeadlineExceeded` 504。对外 JSON 不漏内部路径。context 篇 `%w` 和 HTTP 实战篇 `writeStoreError` 同一条线。`errors.As(err, &e)` 要传指针的指针。自己的类型若要跨 wrap 被 `Is` 认，实现 `Is(error) bool`。sentinel 用 `errors.New` 的包级变量，不要每次 `errors.New("not found")` 新做一个——地址不同，`==` 失败。

```go
type timeout struct{ n int }

func (t timeout) Error() string { return "timeout" }
func (t timeout) Is(target error) bool {
	_, ok := target.(timeout)
	return ok
}
```

对照：C++ 异常+析构，`catch` 按类型；或 `error_code`。Python `raise` / `except`，`from e` 是 `__cause__`。Go 调用方必须看见 `if err != nil`。

---

### 11、context 为什么必须是第一个参数？能不能塞结构体？

能说出口：取消树的根是入口请求。下游每个阻塞点看同一棵树，根死枝停。签名 `func F(ctx context.Context, ...)` 是生态合同。`http.Get`、没有 `*Context` 的 SQL 把树掐断。不要把 ctx 放进可复用结构体字段——Client 活过很多请求，字段里那份要么 Background 无意义，要么是上一个请求的 deadline。

往下追：`WithTimeout` / `WithCancel` 每次 `defer cancel()`。取消向下不向上。子 deadline 不能晚于父。`WithValue` 只放请求元数据，键用未导出类型。后台任务 `WithoutCancel` 再自己加超时，cancel 在 goroutine 里 defer。

```go
func GetUser(ctx context.Context, id int) (User, error) {
	cctx, cancel := context.WithTimeout(ctx, 80*time.Millisecond)
	defer cancel()
	return store.Get(cctx, id)
}

type Client struct {
	ctx context.Context // 不要
	url string
}
```

再往下：`select { case <-ctx.Done(): return ctx.Err() }` 是手写 API 的形状。`Mutex.Lock` 不认 ctx，长锁要自己用 channel 锁或 `TryLock`。`errgroup.WithContext` 扇出。HTTP 实战把 `r.Context()` 派生超时再 `r.WithContext` 传进 Store。

对照：context 篇；C++ `stop_token` / asio cancellation_slot 要自己绑每个 op；Python `CancelledError` 抛进协程，`contextvars` 隐式。Go 把 ctx 写在参数上，看得见。

---

### 12、什么时候用 mutex，什么时候用 channel？

能说出口：保护同一块内存的读写，mutex。移交所有权、等待事件、把流水线串起来，channel。计数器、map、缓存条目用锁更直接；把计数器扔进一个 G 用消息 +1，是多一次拷贝和调度。持锁时不要做会阻塞的 channel 收发，临界区变死锁窗口。

往下追：无缓冲 channel 是会合，happens-before 在发送完成和对应接收之间。锁的 happens-before 在 Unlock 和下一把 Lock 之间。两种都能同步，选语义清楚的那个，不要「看起来并发」。

```go
type counter struct {
	mu sync.Mutex
	n  int
}

func (c *counter) Inc() { c.mu.Lock(); c.n++; c.mu.Unlock() }

// 所有权：算完的结果交给唯一消费者
results := make(chan Result, n)
```

再往下：`RWMutex` 读多写少。`sync.Map` 见 map 题。channel 当锁：容量 1 的缓冲，拿 token 放 token，能 `select` Done，mutex 做不到取消。这是「锁也要能超时」时 channel 赢的地方。持锁调下游 HTTP 是把别人的延迟焊进你的临界区，别的 G 全堵。Go 没有可重入锁；同 G 再 `Lock` 自己那把是死锁，不是「同一线程可以进」。

```go
func trySend(ctx context.Context, sem chan struct{}) error {
	select {
	case sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
```

对照：channel 篇第五节；C++ 默认 mutex，queue 要自己写；Python `Lock` vs `asyncio.Queue`。Go 把 channel 做成语言级，不表示处处替代锁。

---

### 13、select 多个 case 都就绪时走哪条？

能说出口：进入 select 时所有 case 一起看。多个就绪 **均匀随机** 挑一个，不是从上到下。没有 default 且都没就绪，G 阻塞。有 default，立刻走 default。空 `select{}` 永远停这个 G。

往下追：已关闭 channel 的接收永远就绪。nil channel 永不就绪，用来摘支。循环里 `time.After` 每次新 timer。`default` 做非阻塞试探：`select { case ch<-v: default: 丢掉或报忙 }`。

```go
select {
case v := <-a:
	fmt.Println("a", v)
case v := <-b:
	fmt.Println("b", v)
default:
	fmt.Println("none")
}
```

再往下：依赖 case 书写顺序是 bug。测试里两边同时关，哪条先跑不确定，断言顺序会抖。fan-in 关了赋 nil。超时：`NewTimer` + `Stop`，或 ctx。HTTP 超时是 ctx，不是每个叶子手写 `After`。`select{}` 在 `main` 里等于停主 G、进程不退，其它 G 继续——别当「等信号」用，信号用 `signal.NotifyContext`。

```go
func merge(a, b <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil
					continue
				}
				out <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				out <- v
			}
		}
	}()
	return out
}
```

对照：channel 篇第三节；C++ 没有等价物，`poll`/`select(2)` 是 fd；Python `asyncio.wait` 可 FIRST_COMPLETED，不是随机一条语言级 case。

---

### 14、make 和 new 有什么区别？

能说出口：`new(T)` 分配一个 T 的零值，返回 `*T`。`make` 只给 slice / map / channel：分配并初始化内部结构，返回**类型本身**不是指针。`new([]int)` 是 `*[]int` 指向 nil 切片，几乎没人用。切片、map、channel 用 `make` 或字面量。

往下追：`make([]T, n)` len=cap=n 元素零值；`make([]T, 0, n)` 随后 append；`make(map[K]V, hint)` hint 不是容量保证；`make(chan T, 0)` 无缓冲，`make(chan T, n)` n 个格子。`new` 任何类型都能用，包括结构体，等价 `&T{}` 对零值。

```go
p := new(int)          // *int，*p==0
s := make([]int, 2, 5) // [0 0] len2 cap5
m := make(map[string]int)
ch := make(chan int, 1)
```

再往下：`make` 不是 malloc 的 Go 名。堆还是栈仍由逃逸分析决定，`new` 的对象可以留栈。C++ `new T` 是堆+构造；Go `new` 不保证堆。Python 没有这层区分。复合字面量 `&T{...}` 常替代 `new`，能写字段。零值可用：`var s []int` 是 nil 切片，`append` 仍合法；`var m map[K]V` 写会 panic，必须 `make`。channel 零值是 nil，收发阻塞，这是 select 摘支的根，不是 `make` 忘了那么简单。

```go
var s []int
s = append(s, 1) // ok，从 nil 长出来

var m map[int]int
// m[1] = 1 // panic
m = make(map[int]int)
m[1] = 1
```

对照：切片篇 `new([]int)` 那句；入门篇零值。面试把 `make` 说成「堆分配」不准确。

---

### 15、string 和 `[]byte` 转换贵在哪？

能说出口：string 内部 `ptr+len`，只读，无 cap。`[]byte(s)` / `string(b)` 语义上拷贝字节，保证改一边不影响另一边。`s[i]` 是 byte；`range s` 是 rune。`len("中")==3`。编译器在「转换后原值不再改」时可能优化掉拷贝，那是优化不是保证。

往下追：热路径反复转是分配器客户。`strings.Builder`、`strconv.Append*`、直接写 `[]byte`。`unsafe` 零拷贝能让 string 的只读合同被 `[]byte` 改掉，别在生产协议栈上赌。checksum 可以，前提是这段 bytes 转换后不再写。

```go
s := "Go中"
b := []byte(s)
b[0] = 'g'
fmt.Println(s, string(b)) // Go中 go中
```

再往下：JSON、`fmt`、map 的 key 是 string 时，从 buffer 来的 key 会拷一次。`sync.Pool` 的 `[]byte` Put 前若曾转成 string 且 string 还活着，不能再改这块 buf——除非你确定做了拷贝。`s[1:3]` 按字节切，可能切在 rune 中间，类型仍是 string，只是不再合法 UTF-8。改 UTF-8 用 `[]rune` 再转回，又一次拷。

```go
r := []rune("Go中")
r[2] = '文'
fmt.Println(string(r))
```

对照：切片篇 string 节；C++ `string` 可变（或 `string_view` 不拥有）；Python `str` 不可变，`bytes` / `bytearray` 分开。Go 用拷贝保住不可变。

---

### 16、方法集怎么算？为什么值赋不进带指针方法的接口？

能说出口：`T` 的方法集 = `(T)` 接收者。`*T` 的方法集 = `(T)` 加 `(*T)`。调用处 `s.Ptr()` 编译器能取地址；**赋给接口不行**——接口存动态值的拷贝，不会回头取原变量地址来补方法。`var _ I = S{}` 失败，`var _ I = &S{}` 成功。

往下追：值方法改的是副本；指针方法改原对象。`sync.Mutex` 不能拷，必须指针接收者或嵌入指针。方法值 `t.Inc` 把接收者藏进闭包，`t` 常逃逸。

```go
type S struct{ n int }
func (S) Val()     {}
func (*S) Ptr()    {}

var s S
s.Ptr()                 // 语法糖 (&s).Ptr
// var _ interface{ Ptr() } = s  // 编译失败
var _ interface{ Ptr() } = &s    // ok
```

再往下：接口调用走 itab，对象本身无 vptr。小接口（一两个方法）是风格也是方法集合同：实现方不 import 接口包。HTTP 实战 `http.Handler` 就一个 `ServeHTTP`。

对照：函数篇方法集；接口篇 itab。C++ 成员函数在类里声明，虚的进 vtable，没有「值 / 指针方法集不对称赋接口」。Python 绑定 `self` 在运行时。

---

### 17、嵌入是继承吗？

能说出口：不是。嵌入提升字段名和方法名，外层不是内层的子类型，不能把 `Car` 赋给 `Engine` 参数。没有虚分派到外层「重写」——外层同名方法只是挡住提升项，内层的还在 `car.Engine.Start`。没有 is-a。

往下追：`Car` 值方法集含 `Engine` 的值方法；`*Car` 含 Engine 和 `*Engine`。嵌入 `*Engine` 时提升指针方法集，零值 nil，提升调用会 panic，构造时填指针。提升让外层碰巧满足内层满足的接口。

```go
type Engine struct{}
func (Engine) Start() { fmt.Println("eng") }

type Car struct{ Engine }
func (Car) Start() { fmt.Println("car") }

c := Car{}
c.Start()        // car
c.Engine.Start() // eng
```

再往下：用来组合 mutex：`type T struct{ mu sync.Mutex; n int }` 不要提升成 `t.Lock()` 对外——把锁导出等于让调用方参与你的临界区。嵌入 `http.ResponseWriter` 做 `statusWriter` 是装饰，HTTP 实战中间件那样。同名提升冲突：外层定义挡住；两个匿名字段提升同一名字，外层必须自己写一个，否则调用歧义编译失败。没有虚析构问题——没有继承、没有 `delete` 基类指针。

```go
type A struct{}
func (A) F() {}
type B struct{}
func (B) F() {}
type C struct {
	A
	B
}
// c.F() // 编译失败：ambiguous
```

对照：函数篇嵌入；C++ 公开继承 is-a + 虚函数；Python 类继承 MRO。Go 组合优先，嵌入是语法糖不是继承声明。

---

### 18、GOMAXPROCS 改的是什么？调大一定更快吗？

能说出口：改的是 P 的个数，即同时跑 **Go 代码** 的并行度。默认 `NumCPU()`。M 可以更多：阻塞 syscall 会另雇 M。G 可以很多。调到 1，Go 代码串行交错，channel 仍工作，data race 窗口变小但不是消失——单 P 仍抢占切换 G。

往下追：CPU 密集且真并行，P 对齐逻辑核。IO 密集，G 等在 netpoll，再加 P 不帮忙，反而 cache 和调度开销。容器里 CPU quota 小于 `NumCPU()` 时要显式调，否则 P 过多，GC 和调度以为有那些核。

```go
runtime.GOMAXPROCS(4)
fmt.Println(runtime.GOMAXPROCS(0))
fmt.Println(runtime.NumGoroutine())
fmt.Println(runtime.NumCPU())
```

再往下：`debug.SetGCPercent`、`GOMEMLIMIT` 和 P 数量一起影响 pacer。`sync.Pool` 按 P 分池，P 很多对象散。不要为了 Pool 去调 `GOMAXPROCS`。cgo 密集会制造很多 M，那是 syscall 路径，不是把 P 调大能消化的。容器 `cfs_quota` 2 核、镜像里看见 64 个 sibling，默认 64 个 P 会把自己调度打穿——显式 `GOMAXPROCS=2` 或 `automaxprocs`。

```go
func hog() {
	for {
		x := 0
		for i := 0; i < 1e8; i++ {
			x += i
		}
		_ = x
		// 1.14 前这种循环协作点少，能饿死同 P 上别的 G
		// 现在信号抢占会进来，但仍占满一个 P
	}
}
```

对照：GMP 篇 P 一节；C++ 线程池 size；Python GIL 下调线程数几乎不帮 CPU 密集。Go 调 P 是调「允许多少并行 Go 代码」。

---

### 19、什么是 data race？Go 怎么查？

能说出口：两个 G 碰同一变量，至少一个写，没有 happens-before。未定义行为，不是「偶尔脏读」那么温和。同步：channel 收发、mutex、WaitGroup、Once、atomic、`go` 启动 happens-before 那个 G 开始。没有 GIL。`go test -race` / `-race` 编译当常规。

往下追：闭包抓循环变量（1.22 前同一份 `i`；1.22 起每次迭代新变量，旧代码坑还在老版本）。`http.Request` 在 handler 返回后被 `go` 再用。slice header 拷走后底层数组仍共享，一边 append 一边读是 race。map 并发是 fatal 或 race 两条路。

```go
var a int
go func() { a = 1 }()
fmt.Println(a) // race
```

再往下：`-race` 有内存和 CPU 税，生产默认关，CI 开。它查的是执行到的路径。`atomic.Int64` 单变量计数；多字段仍要锁。假共享：相邻 `int64` 被不同 P 写。race 报告的栈是线索，根因常是「谁以为自己独享这块数组」。1.22 前：

```go
for i := 0; i < 3; i++ {
	go func() { fmt.Println(i) }() // 常打 3 3 3，且是 race
}
```

1.22 起每次迭代新 `i`。老模块语言版本仍是旧语义。`go test -race ./...` 写进 CI，比面试当场背 happens-before 清单有用。

对照：channel 篇内存模型；C++ 同款 data race UB，TSan；Python GIL 遮住一部分，无 GIL 就回来了。Go 把检测做成一等工具链。

---

### 20、init 按什么顺序跑？为什么别在里面连网？

能说出口：先按导入图深度优先初始化依赖包（每包一次）；包内按文件名排序，先求值包级变量，再跑所有 `init`（可多个）；最后 `main.main`。你不能调用 `init`。空白导入 `_ "pkg"` 就是为了跑它的 init（driver 注册）。

往下追：`init` 没有参数、没有 `error` 返回值，失败只能 panic。连网、读会失败的配置、听端口，放到 `main` 里显式调用。包级 `var x = f()` 若 `f` 依赖别的包状态，顺序坑比看起来少（有导入图），但仍别搞跨包副作用图。

```go
var defaultPort int

func init() { defaultPort = 8080 }
func init() { /* 同包多个，按文件名再按出现顺序 */ }
```

再往下：测试 `TestMain` 在 init 之后。循环导入编译失败，所以 init 图是 DAG。C++ 跨翻译单元静态初始化顺序未定义（fiasco）；Go 确定，但确定不等于「适合做 IO」。Python 模块顶层代码≈整个文件是 init。同包多个文件的 `init`：先文件名，再出现顺序；不要靠「我写在上面所以先跑」——文件名一改顺序就变。包级 `var _ = register()` 和 `init` 等价副作用，更难搜。

```go
// driver.go
func init() { sql.Register("foo", &drv{}) }

// main
import _ "example.com/foo" // 只为跑上面那行
```

对照：入门篇 init；HTTP 实战听端口在 `main`，Store 构造在 `main`，不藏进 init。

---

### 21、Go 为什么很久没有泛型，1.18 才有？

能说出口：早期刻意不做：要编译快、要类型简单、接口已经覆盖「运行时异构」（`io.Reader`）。缺的是「编译期一族代码」：`Sum([]int)` 和 `Sum([]float64)` 不能一份不装箱的实现。拖了十年，社区用 `any`+断言、`go generate`、`interface{ ~int }` 的草案吵了好几轮。**1.18** 落地类型参数和约束。

往下追：接口不是泛型。`func f(x any)` 进去要断言，运行时。泛型 `func Sum[T int|float64](xs []T) T` 编译期特化（或字典），热路径为了**不装箱**，不是为了少写重载。`io.Reader` 仍不该写成泛型——实现集合开放。

```go
func Sum[T int | float64](xs []T) T {
	var s T
	for _, x := range xs {
		s += x
	}
	return s
}
```

再往下：约束是接口。`comparable` 才能当 map 键。方法上的类型参数 1.18 有限制，后来版本在放松。容器类（泛型 Set）1.18 之后才自然；之前 `map[T]struct{}` 或接口。面试别答「Go 没有泛型」——那是 1.17 及更早。也别把所有 `any` 改成泛型：开放集合用小接口。`any` 约束的泛型方法调用仍可能走字典、仍可能装箱，热路径用 `~int` 这类近似类型或干脆具体类型。GC 形状：`[]any` 是 N 次盒子；`[]T` 是连续 T。

```go
func Map[S ~[]E, E any](s S, f func(E) E) S {
	out := make(S, len(s))
	for i, v := range s {
		out[i] = f(v)
	}
	return out
}
```

对照：接口篇「接口不是泛型」；C++ 模板是编译期、错误信息地狱、头文件膨胀，Go 晚做是为了避开那套成本；Python 注解不生成多份机器码。

---

### 22、sync.Pool 是对象池吗？Get 到的一定是上次 Put 的吗？

能说出口：不是 C++ 那种「还回去一定还在」的池。按 P 分本地池，抗的是**两次 GC 之间**的分配风暴。每次 GC 会清 Pool。Put 进去的东西下次 Get 可能 `New()` 全新一只。合同：可丢、可重置、无必须 Close 的资源。数据库连接放进去，GC 一来对象丢了没关，泄漏。

往下追：`Get` 当前 P 没有就 `New`。Put 时 G 可能已在另一个 P，对象进那一侧。Put 前把大切片、指针字段抹掉，否则两次 GC 之间 10MiB 仍钉在 HeapAlloc。先 pprof 证明这块 `make` 热，再 Pool。

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

再往下：对象很大很少分配、每次大小都不同、持有 fd，不要 Pool。`fmt` / `json` 内部大量 Pool，抄模式时把重置带上。C++ 对象池篇那句「不是 allocator」这里换成「不是缓存、不是所有权」：Go 的 Pool 随时可以扔给你的对象。Put 一个 10MiB 的 `[]byte`，下一次 GC 前这块一直在；Get 到的 cap 可能比你这次需要的大很多，自己 `[:0]` 再用，不要假定 cap==4096。Pool 不是跨请求的 session 存储。

```go
func put(p *sync.Pool, r *req) {
	r.Body = nil // 否则 Pool 活着，Body 钉到下一次 GC
	r.next = nil
	p.Put(r)
}
```

对照：GC 篇第五节；C++ 实战对象池 mutex+链表，析构才真正 `delete`；Python 几乎不写池。Go 有 GC 兜底，所以合同更弱，也更不容易把连接池写进去。

---

检查清单：每题先一句合同（header 共享还是拷、谁 close、tab/data 哪一字非零、P 不是 M），再一句实现（growslice、itab、混合写屏障、本地队列 256），再一句事故（fatal map、typed nil、G 泄漏、Pool 清了没 Close）。追问往「append 之后旧头还指哪」「ctx 为什么不进结构体」「1.18 泛型解决的是装箱不是 Reader」走，不往背诵关键字走。基础四篇钉盒子和接口，进阶五篇钉调度、channel、GC、ctx、逃逸；HTTP 实战把 ctx、Handler 接口、recover、`-race` 收进一个能跑的进程。本篇把同一根轴收成能说出口的句子。
