# Golang入门

## 一、基本数据类型

### 1、整数类型

```go
var num3 int = 3
var num4 int8 = 4      // 占1个字节   -2^7  ~ 2^7 - 1
var num5 int16 = 5     // 占2个字节   -2^15 ~ 2^15 - 1
var num6 int32 = 6     // 占4个字节   -2^31 ~ 2^31 - 1
var num7 int64 = 7     // 占8个字节   -2^63 ~ 2^63 - 1
```

### 2、浮点类型

```go
var num1 float32 = 1
var num2 float64 = 2
```

### 3、字符类型

#### 转义字符

```go
fmt.Println("aaa\nbbb")     // \n 换行
fmt.Println("aaa\bbbb")     // \b 退格
fmt.Println("aaa\rbbbb")    // \r 光标回到本行的开头，后续输入就会替换原有的字符
fmt.Println("aaaaa\tbbbbb") // -t 制表符 8位

输出：
aaa
bbb
aabbb
bbbb
aaaaa   bbbbb
```

### 4、布尔类型

### 5、字符串类型

```go
var s1 string = "hello golang"
fmt.Println(s1)
```

字符串是不可变的：指的是字符串一旦定义好，其中的字符的值不能改变。

### 6、基本数据类型的默认值

```go
var num1 float32
var num2 float64
var num3 int
var num4 int8
var num5 int16
var num6 int32
var num7 int64
var num8 string
var num9 byte
var num10 bool
fmt.Println(num1)
fmt.Println("float32: ", unsafe.Sizeof(num1))
fmt.Println(num2)
fmt.Println("float64: ", unsafe.Sizeof(num2))
fmt.Println(num3)
fmt.Println("int: ", unsafe.Sizeof(num3))
fmt.Println(num4)
fmt.Println("int8: ", unsafe.Sizeof(num4))
fmt.Println(num5)
fmt.Println("int16: ", unsafe.Sizeof(num5))
fmt.Println(num6)
fmt.Println("int32: ", unsafe.Sizeof(num6))
fmt.Println(num7)
fmt.Println("int64: ", unsafe.Sizeof(num7))
fmt.Println(num8)
fmt.Println("string: ", unsafe.Sizeof(num8))
fmt.Println(num9)
fmt.Println("byte: ", unsafe.Sizeof(num9))
fmt.Println(num10)
fmt.Println("bool: ", unsafe.Sizeof(num10))


输出结果：
0
float32:  4
0
float64:  8
0
int:  8
0
int8:  1
0
int16:  2
0
int32:  4
0
int64:  8

string:  16
0
byte:  1
false
bool:  1
```

### 7、数据类型之间的转换

#### 7.1、基本数据类型之间的转换

```go
var n1 int = 100
// var n2 float32 = n1  // 在这里自动转换不好使，比如显式转换
var n2 float32 = float32(n1)
fmt.Println(n1)
fmt.Println(n2)
fmt.Printf("%T\n", n1) // int  n1的类型还是int类型，只是将n1的值100转为了float32而已，n1还是int类型

// 将int64转为int8的时候，编译不会出错，但是会数据溢出
var n3 int64 = 888888
var n4 int8 = int8(n3)
fmt.Println(n4) //56

var n5 int32 = 12
var n6 int64 = int64(n5) + 30 // 一定要匹配 = 左右的数据类型
fmt.Println(n5)
fmt.Println(n6)

var n7 int64 = 12
var n8 int8 = int8(n7) + 127 // 编译通过，但是结果可能会溢出
// var n9 int8 = int8(n7) + 128 // 编译不会通过
fmt.Println(n8)
// fmt.Println(n9)
```

#### 7.2、基本数据类型和字符串类型之间的转换

##### 方式一、Sprintf("", )

![Sprintf](https://cdn.jsdelivr.net/gh/aqjsp/Pictures/image-20241107225238132.png)

```go
var n1 int = 19
var s1 string = fmt.Sprintf("%d", n1)
fmt.Printf("s1对应的类型是： %T, s1 = %q \n", s1, s1)

var n2 float32 = 1.23
var s2 string = fmt.Sprintf("%f", n2)
fmt.Printf("s2对应的类型是： %T, s2 = %q \n", s2, s2)

var n3 bool = false
var s3 string = fmt.Sprintf("%t", n3)
fmt.Printf("s3对应的类型是： %T, s3 = %q \n", s3, s3)

var n4 byte = 'a'
var s4 string = fmt.Sprintf("%c", n4)
fmt.Printf("s4对应的类型是： %T, s4 = %q \n", s4, s4)
```

##### 方式二、strconv

## 二、复杂数据类型

### 1、指针

```go
var num int = 18
// &符号+变量，就可以获取这个变量内存的地址
fmt.Println(&num)

// 定义一个指针变量：
// var代表要声明一个变量
// ptr 指针变量的名字
// ptr 对应的类型是：*int 是一个指针类型 ，可以理解为 指向int类型的指针
// &num 就是一个地址，是ptr变量的具体的值
var ptr *int = &num
fmt.Println(ptr)
fmt.Println("ptr本身的存储空间的值：", &ptr)

// 想获取ptr这个指针或者这个地址指向的那个数据
fmt.Printf("ptr指向的数值：%v \n", *ptr)

输出：
0xc000084048
0xc000084048
ptr本身的存储空间的值： 0xc00004c008
ptr指向的数值：18 
```

总结：

1. &  取内存地址
2. *根据地址取值

#### 1.1、可以通过指针改变指向值

```go
var num int = 10
fmt.Println(num)

var ptr *int = &num
*ptr = 20
fmt.Println(num)
```

#### 1.2、指针变量接收的一定是地址值

```go
var ptr *int = num // 错误
```

#### 1.3、指针变量的地址不可以不匹配

```go
var ptr *float32 = &num
```

#### 1.4、基本数据类型（又叫值类型）

都有对应的指针类型，形式为 * 数据类型，比如int对应的指针就是* int，float32对应的指针类型就是* float32，依次类推。

### 2、

## 三、运算符

## 四、流程控制

### 关键字

#### `break`

作用：

- `break` 用于立即终止循环（`for`、`switch` 或 `select`），并跳出当前的循环块或 `switch` 语句。
- 在多层嵌套循环中，`break` 仅作用于它所在的最近一层循环。

使用场景：提前退出循环，如在查找某个满足条件的元素时。

例子：

```go
for i := 0; i < 10; i++ {
    if i == 5 {
        break // 终止循环，当 i == 5 时跳出 for 循环
    }
    fmt.Println(i)
}

输出：
0
1
2
3
4
```

#### `continue`

作用：

- `continue` 用于跳过当前循环的剩余部分，并立即进入下一次循环迭代。
- 在 `for` 循环中，当执行 `continue` 时，控制权会返回到循环的增量表达式（如果有的话）。

使用场景：跳过当前不符合条件的迭代，继续下一次迭代，例如跳过奇数或某些条件下的值。

例子：

```go
for i := 0; i < 10; i++ {
    if i%2 == 0 {
        continue // 跳过本次循环的剩余部分，直接进入下一次循环
    }
    fmt.Println(i)
}

输出：
1
3
5
7
9
```

#### `goto`

作用：

- `goto` 是一种无条件跳转语句，用于将控制转移到同一函数或方法中的某个标签位置。
- 标签是由一个标识符加冒号 (`:`) 构成的，可以放置在代码的任意位置。

使用场景：在 Go 中不推荐频繁使用 `goto`，因为它可能会让代码难以理解。但在某些情况下（如在嵌套循环中需要跳出多层循环），`goto` 可以作为一种简便的控制方式。

例子：

```go
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if i == 1 && j == 1 {
            goto end // 跳转到 end 标签，直接结束嵌套循环
        }
        fmt.Println(i, j)
    }
}
end: // 标签位置
fmt.Println("跳出了嵌套循环")

输出：
0 0
0 1
0 2
1 0
跳出了嵌套循环
```

#### `return`

作用：

- `return` 用于从函数中返回值，并结束该函数的执行。执行 `return` 后，函数不会再执行 `return` 语句后的任何代码。
- 如果函数有返回值，`return` 语句后跟随返回的值。

使用场景：返回函数的计算结果或状态，或在特定条件下提前退出函数。

例子：

```go
func add(a int, b int) int {
    return a + b // 返回 a + b 的结果，并结束函数
}

func main() {
    sum := add(3, 5)
    fmt.Println(sum) // 输出: 8
}

输出：
8
```

## 五、函数

### 1、简单函数

```go
func 函数名 (形参列表) (返回值类型列表) {
	执行语句...
	return + 返回值列表
}
```

例子：

```go
func cal(a int, b int) int {
	return a + b
}
func main() {
	sum := cal(10, 20)
	fmt.Println(sum)
}
```

### 2、函数详解

```go
func cal(a int, b int) (int, int) {
	return a + b, a - b
}
func main() {
	sum, res := cal(10, 20)
	fmt.Println(sum, res)
}
```

### 3、闭包

#### 3.1、定义

闭包是指内部函数捕获并引用外部函数作用域中的变量。使用闭包时，这些外部变量的生命周期会延长，直到闭包不再引用它们为止。Go 的闭包主要通过匿名函数（`func`）实现。

#### 3.2、例子

```go
func getSum() func(int) int {
	var sum int = 0
	return func(num int) int {
		sum += num
		return sum
	}
}

func main() {
	f := getSum()
	fmt.Println(f(1))
	fmt.Println(f(2))
	fmt.Println(f(3))
}
```

#### 3.3、原理

闭包能引用外部变量的原因在于，Go 语言会将这些变量的引用包含在闭包的环境中。闭包在创建时将变量的引用“捕获”在自身的上下文中，而不是将值复制进去。这种特性允许闭包函数在后续的调用中继续访问这些变量的最新值。

#### 3.4、使用场景



### 4、defer关键字

在 Golang 中，`defer` 关键字用于在函数结束时执行一些代码。`defer` 语句的核心特点是无论函数如何结束——正常返回、遇到错误或是发生 `panic`——被 `defer` 的代码都会在函数退出前执行。因此，`defer` 常用于释放资源、关闭文件、解锁互斥锁等清理操作。

#### `defer` 的基本语法和用法

`defer` 语句在函数体内定义，需要跟随一个函数调用。当执行到 `defer` 语句时，Go 会记录该语句并推迟到当前函数返回时再执行。

**语法**：

```go
defer functionName(args)
```

**示例**：

```go
package main

import "fmt"

func main() {
    fmt.Println("Start")
    defer fmt.Println("This is deferred")
    fmt.Println("End")
}
```

**输出**：

```sql
Start
End
This is deferred
```

在这个例子中，`defer fmt.Println("This is deferred")` 被推迟到 `main()` 函数即将结束时执行，所以它最后输出。

#### `defer` 的执行顺序

如果在一个函数中有多个 `defer` 语句，这些 `defer` 调用会按照后进先出的顺序（LIFO）执行。这类似于栈结构的后入先出。

**示例**：

```go
package main

import "fmt"

func main() {
    defer fmt.Println("First")
    defer fmt.Println("Second")
    defer fmt.Println("Third")
}
```

**输出**：

```sql
Third
Second
First
```

`defer` 会按顺序将所有的调用入栈，在函数退出时逆序执行这些 `defer` 语句。

#### `defer` 的常见用途

1. **释放资源**：在操作文件、数据库连接等资源时，`defer` 可以确保它们在函数结束前正确释放。

   ```go
   package main
   
   import (
       "fmt"
       "os"
   )
   
   func main() {
       file, err := os.Open("test.txt")
       if err != nil {
           fmt.Println("Error opening file:", err)
           return
       }
       defer file.Close() // 确保在函数结束时关闭文件
   
       // 读取文件内容
       // ...
   }
   ```

2. **解锁互斥锁**：在多线程程序中使用 `sync.Mutex` 或 `sync.RWMutex` 时，可以使用 `defer` 来确保锁在操作完成后正确解锁。

   ```go
   var mu sync.Mutex
   
   func criticalSection() {
       mu.Lock()
       defer mu.Unlock() // 确保在函数结束前解锁
       
       // 执行一些需要锁保护的操作
   }
   ```

3. **捕获并处理异常**：在 Go 中，`defer` 配合 `recover` 可以捕获 `panic` 并处理，避免程序崩溃。

   ```go
   func safeDivision(a, b int) {
       defer func() {
           if r := recover(); r != nil {
               fmt.Println("Recovered from panic:", r)
           }
       }()
       
       fmt.Println("Result:", a/b) // 若 b 为 0，会引发 panic
   }
   ```

#### `defer` 的工作原理与参数评估

`defer` 语句会在它定义时计算参数值，但会推迟执行实际的函数调用。这意味着 `defer` 语句捕获的是当前的参数值。

**示例**：

```go
package main

import "fmt"

func printValue(val int) {
    fmt.Println("Deferred value:", val)
}

func main() {
    x := 10
    defer printValue(x) // 此时 x 的值是 10
    x = 20
    fmt.Println("Updated value:", x)
}
```

**输出**：

```yaml
Updated value: 20
Deferred value: 10
```

尽管 `x` 在 `defer` 语句之后被更新为 20，但 `defer` 捕获的是当时的值（10），所以 `printValue(10)` 被推迟执行。

#### `defer` 在返回值中的应用

在有命名返回值的函数中，`defer` 可以修改返回值，因为 `defer` 语句会在返回前执行。

**示例**：

```go
package main

import "fmt"

func modifyReturnValue() (result int) {
    defer func() {
        result += 10 // 修改命名返回值
    }()
    return 5
}

func main() {
    fmt.Println(modifyReturnValue()) // 输出: 15
}
```

在此例中，`result` 是命名返回值，函数返回值会在 `defer` 执行时被修改。

#### 注意事项

1. **性能**：`defer` 虽然方便，但在高性能需求的代码中需谨慎使用。因为 `defer` 的实现相对开销较高，尤其在大量循环中使用时可能会影响性能。
2. **不适合实时性要求的操作**：由于 `defer` 的延迟执行特点，不适合对实时性要求高的操作场景。
3. **参数捕获机制**：`defer` 在定义时捕获的参数可能与最终执行时的上下文不同，需注意避免误用。

#### 总结

- `defer` 用于推迟代码执行至函数退出前，适合资源管理、异常处理等场景。
- `defer` 语句遵循后进先出顺序执行。
- 其参数在定义时评估，实际调用在函数退出前。

### 5、系统函数

#### 5.1、字符串相关函数

##### 5.1.1、字符串遍历

1. `for - range`

```go
func main() {
	str := "hello golang你好"
	for i, value := range str {
		fmt.Printf("下标：%d, 值：%c \n", i, value)
	}
}
```

2. 切片：`r := []rune(str)`

```go
func main() {
	str := "你好 Golang"
	r := []rune(str)
	for i := 0; i < len(r); i++ {
		fmt.Printf("%c \n", r[i])
	}
}
```

##### 5.1.2、字符串转整数

```go
n, err := strconv.Atoi("66")
```

##### 5.1.3、整数转字符串

```go
str = strconv.Itoa(8888)
```

##### 5.1.4、查找子串是否在指定的字符串中

```go
strings.Contains("fadshjkfhakj", "fa")

例子：
flags := strings.Contains("fadshjkfhakj", "fdaa")
fmt.Println(flags)
```

##### 5.1.5、统计一个字符串有几个指定的子串

```go
strings.Count("fdafadgadf","dfa")

例子：
count := string.Count("fdafag","a")
fmt.Println(count)
```

##### 5.1.6、不区分大小写的字符串比较

```go
strings.EqualFold("hello", "HELLO")

例子：
flags := strings.EqualFold("hello", "HELLO")
fmt.Println(flags)
```

##### 5.1.7、返回子串在字符串第一次出现的索引值，如果没有返回-1

```go
strings.Index("fdagagagag","ag")

例子：
index := strings.Index("fdagagagag","ag")
fmt.Println(index)
```

##### 5.1.8、字符串的替换

```go
strings.Replace("fzdfadgadgaga", "ga", "qint", n)

例子：
str := strings.Replace("fzdfadgadgaga", "ga", "qint", -1) // 全部替换
str := strings.Replace("fzdfadgadgaga", "ga", "qint", 2) // 替换2个
fmt.Println(str)
```

##### 5.1.9、按照指定的某个字符，为分割标识，将一个字符串拆分成字符串数组

```go
strings.Split("go-python-golang", "-")
例子：
arr := strings.Split("go-python-golang", "-")
```

##### 5.1.10、大小写切换

```go
strings.ToLower("GO")  // to 小写
strings.ToUpper("go")  // to 大写
```

##### 5.1.11、将字符串左右两边的空格去掉

```go
strings.TrimSpace("   go java cpp    ")

例子：
str := strings.TrimSpace("   go java cpp    ")  // 不会去掉字符串中间的空格
fmt.Println(str)
```

##### 5.1.12、将左右指定的字符去掉

```go
strings.Trim("--golang--", "-")

例子：
str := strings.Trim("--golang--", "-")
fmt.Println(str)
```

##### 5.1.13、判断字符串是否以指定的字符串开头

```go
strings.HasPrefix("https://www.aqjszz.com/docs", "https")

例子
flags := strings.HasPrefix("https://www.aqjszz.com/docs", "https")
fmt.Println(flags)
```

##### 5.1.14、判断字符串是否以指定的字符串结尾

```go
strings.HasSuffix("fadf.png", "jpg")
```

#### 5.2、日期和时间相关函数

```go
now := time.Now()
fmt.Println(now.Format("2006-01-02 15:04:05"))
```

#### 5.3、内置函数

##### 5.3.1、len

`len` 用于获取容器类型的长度或大小。它适用于多种数据结构，包括数组、切片、字符串、映射（map）和通道（channel）。

**语法**：

```go
len(v interface{}) int
```

- **数组**：返回数组的长度。
- **切片**：返回切片当前包含的元素个数，而不是底层数组的容量。
- **字符串**：返回字符串的字节数（注意：对于包含非 ASCII 字符的字符串，字节数可能不等于字符数）。
- **映射（map）**：返回映射中的键值对数量。
- **通道（channel）**：返回通道中排队等待接收的元素数量（适用于缓冲通道）。

**示例**：

```go
package main

import "fmt"

func main() {
    arr := [5]int{1, 2, 3, 4, 5}
    slice := []int{1, 2, 3}
    str := "hello"
    m := map[string]int{"one": 1, "two": 2}
    ch := make(chan int, 5)
    ch <- 1
    ch <- 2

    fmt.Println("数组长度:", len(arr))   // 输出: 5
    fmt.Println("切片长度:", len(slice)) // 输出: 3
    fmt.Println("字符串长度:", len(str)) // 输出: 5
    fmt.Println("映射长度:", len(m))     // 输出: 2
    fmt.Println("通道长度:", len(ch))    // 输出: 2
}
```

##### 5.3.2、new

`new` 用于分配内存，但只会返回类型的指针，并不会进行初始化。它适合用于一些简单的类型，例如数值类型、结构体等。`new` 返回的是指向类型的零值的指针。

**语法**：

```go
new(Type) *Type
```

- **用法**：`new(T)` 会分配一个 T 类型的内存空间并返回指针，指针指向的值是类型 T 的零值。
- `new` 主要用于需要一个指针类型而不需要额外初始化的场景，适合少量分配和直接访问的情况。

**示例**：

```go
package main

import "fmt"

func main() {
    p := new(int)     // 分配一个 int 类型的内存，初始值为 0
    *p = 10           // 修改指针指向的值
    fmt.Println(*p)   // 输出: 10

    s := new([]int)   // 分配一个指向空切片的指针
    fmt.Println(s)    // 输出: &[]
}
```

**注意**：`new` 分配的只是内存，没有初始化。对于结构体类型 `new(T)` 创建的指针等价于 `&T{}`。

##### 5.3.3、make

`make` 用于初始化和分配引用类型（即：切片、映射和通道）的内存空间。`make` 会返回一个初始化后的值，而不是指针。

**语法**：

```go
make(t Type, size ...IntegerType) Type
```

- **切片（slice）**：`make([]T, len, cap)` 创建一个具有指定长度和容量的切片。
- **映射（map）**：`make(map[K]V, hint)` 创建一个指定键值对数量（容量提示）的映射。
- **通道（channel）**：`make(chan T, buffer)` 创建一个带缓冲区的通道。

**示例**：

```go
package main

import "fmt"

func main() {
    // 创建一个长度为 5 的整数切片
    slice := make([]int, 5, 10)
    fmt.Println("切片长度:", len(slice))  // 输出: 5
    fmt.Println("切片容量:", cap(slice))  // 输出: 10

    // 创建一个映射
    m := make(map[string]int)
    m["foo"] = 42
    fmt.Println("映射:", m)               // 输出: map[foo:42]

    // 创建一个缓冲区大小为 3 的通道
    ch := make(chan int, 3)
    ch <- 1
    ch <- 2
    fmt.Println("通道长度:", len(ch))     // 输出: 2
}
```

**注意**：`make` 仅用于这三种类型，因为它不仅仅是分配内存，还初始化底层数据结构的相关信息。

##### 总结对比

| 函数   | 主要用途                           | 返回值            | 适用类型                       |
| ------ | ---------------------------------- | ----------------- | ------------------------------ |
| `len`  | 获取长度或大小                     | int               | 数组、切片、字符串、映射、通道 |
| `new`  | 分配内存并返回指向零值的指针       | *T（类型的指针）  | 值类型（基本类型、结构体等）   |
| `make` | 分配并初始化内存，返回初始化后的值 | T（初始化后的值） | 切片、映射、通道               |

## 六、错误处理

### 6.1、错误捕获机制

```go
defer + recover
```

![error](https://cdn.jsdelivr.net/gh/aqjsp/Pictures/image-20241108235411227.png)

例子：

```go
func main() {
	test()
	fmt.Println("main end")
}

func test() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("err = ", err)
		}
	}()
	num1 := 10
	num2 := 0
	res := num1 / num2
	fmt.Println(res)
}

输出：
err =  runtime error: integer divide by zero
main end
```

### 6.2、自定义错误

需要调用errors包下的New函数：

func [New](https://github.com/golang/go/blob/master/src/errors/errors.go?name=release#9)

```go
func New(text string) error
```

使用字符串创建一个错误,请类比fmt包的Errorf方法，差不多可以认为是New(fmt.Sprintf(...))。

例子：

```go
func main() {
	err := test()
	if err != nil {
		fmt.Println("err = ", err)
	}
	fmt.Println("main end")
}

func test() (err error) {
	num1 := 10
	num2 := 0
	if num2 == 0 {
		return errors.New("除数不能为0")
	} else {
		res := num1 / num2
		println("num1/num2 = ", res)
		fmt.Println(res)
		return nil
	}
}

输出：
err =  除数不能为0
main end
```

## 七、数组

### 7.1、命名

```go
var 数组名 [数组大小]数组类型
```

### 7.2、初始化

```go
第一种:
var arr1 [3]int = [3]int{1,2,3}
第二种：
var arr2 = [3]int{4,5,6}
第三种：
var arr3 = [...]int{12,34,45,65,67}
第四种：
var arr4 = [...]int{2:22, 1:11, 3:33, 0:55}
```

### 7.3、类型

数组长度属于数组类型。

### 7.4、二维数组

遍历

方式一：for循环遍历

方式二：for-range遍历

```go
func main() {
	var arr [3][3]int = [3][3]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	fmt.Println(arr)

	// 方式一：for循环遍历
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			fmt.Printf("arr[%d][%d] = %d \t", i, j, arr[i][j])
		}
		fmt.Println()
	}

	fmt.Println("-----------------")

	// 方式二：range循环遍历
	for i, v := range arr {
		for j, v2 := range v {
			fmt.Printf("arr[%d][%d] = %d \t", i, j, v2)
		}
		fmt.Println()
	}
}
```

## 八、切片---左闭右开

切片是建立在数组之上

```go
func main() {
	var arr [6]int = [6]int{3, 5, 6, 7, 8, 9}
	// 切片是构建在数组之上
	var slice []int = arr[2:5]
	fmt.Println(slice)
}
```

**切片（slice）** 是一种动态数组。切片提供了便捷、灵活的数组操作方式，底层结构可以动态增长或缩小。

### 切片的底层结构

在 Go 中，切片的底层是基于数组实现的，其底层结构由一个 `sliceHeader` 管理，定义在 `reflect` 包中，包含以下三部分：

1. 指针（Pointer）：指向切片底层数组的起始位置。
2. 长度（Length）：表示切片当前的长度，即切片实际存储的元素个数。
3. 容量（Capacity）：表示从切片起始位置到底层数组末尾的最大可用空间。

因此，切片可以表示为以下结构：

```go
type sliceHeader struct {
    ptr    *elementType // 指向底层数组的指针
    length int          // 当前切片长度
    cap    int          // 当前切片的容量
}
```

### 内存布局示意图

假设有以下代码：

```go
arr := [5]int{1, 2, 3, 4, 5}
s := arr[1:4]
```

在这段代码中，`arr` 是一个长度为 5 的数组，而 `s` 是一个切片，表示 `arr` 中从索引 1 到 3 的部分。底层内存布局如下：

```go
底层数组： [1, 2, 3, 4, 5]
              ↑     ↑
           起始位置   长度为3
切片 s:   指向 arr[1]（值为2），长度为3，容量为4（从arr[1]到arr[4]）
```

- `s.ptr` 指向 `arr[1]` 的位置。
- `s.length` 为 `3`，表示 `s` 包含的元素数量 `[2, 3, 4]`。
- `s.cap` 为 `4`，因为从 `arr[1]` 开始到数组末尾一共有 4 个元素。

### 切片扩容

Go 的切片具有动态扩容能力，当切片使用 `append` 操作超出其容量时，Go 会分配一个新的、更大容量的底层数组，并将原切片的数据复制到新数组中。扩容机制一般遵循以下原则：

1. **小切片（小于 1024 元素）**：每次扩容时，容量翻倍。
2. **大切片（大于等于 1024 元素）**：每次扩容增加原来容量的 1/4 左右。

举例：

```go
s := make([]int, 2, 2) // 创建一个长度为2，容量为2的切片
s = append(s, 1)       // 超出容量，触发扩容
```

在 `append` 操作中，由于超出容量，Go 会分配一个新的切片，长度为 3，容量为 4，并将数据 `[0, 0, 1]` 复制到新切片。

### 切片的共享与独立

在使用切片时要注意切片可能会共享同一个底层数组。例如：

```go
arr := []int{1, 2, 3, 4, 5}
s1 := arr[1:4]
s2 := arr[2:5]
s1[1] = 10
fmt.Println(arr) // 输出：[1, 2, 10, 4, 5]
```

在这里，`s1` 和 `s2` 都指向了 `arr` 的一部分，因此修改 `s1` 的值会影响 `s2` 和 `arr`。要独立地使用切片，可以使用 `copy` 函数创建切片的副本。

### 切片的拷贝

通过 `copy` 函数来实现对切片的拷贝，`copy` 会复制切片的内容，但不会共享底层数组。这意味着拷贝后的切片与原切片独立存储，修改其中一个切片不会影响另一个。

#### `copy` 函数的用法

`copy` 函数的语法为：

```go
copy(dst, src []Type) int
```

- `dst` 是目标切片（拷贝到的切片），`src` 是源切片（被拷贝的切片）。
- `copy` 返回复制的元素个数。它只会复制 `dst` 和 `src` 之间较小长度的部分。

#### 示例代码

```go
package main

import "fmt"

func main() {
    src := []int{1, 2, 3, 4, 5}
    dst := make([]int, len(src))

    // 使用 copy 函数将 src 拷贝到 dst
    count := copy(dst, src)

    fmt.Println("拷贝的元素个数:", count) // 输出：拷贝的元素个数: 5
    fmt.Println("源切片:", src)        // 输出：源切片: [1 2 3 4 5]
    fmt.Println("目标切片:", dst)       // 输出：目标切片: [1 2 3 4 5]

    // 修改目标切片
    dst[0] = 10
    fmt.Println("修改后的源切片:", src) // 输出：修改后的源切片: [1 2 3 4 5]
    fmt.Println("修改后的目标切片:", dst) // 输出：修改后的目标切片: [10 2 3 4 5]
}
```

#### 结果说明

1. `copy(dst, src)` 将 `src` 切片的内容逐一复制到 `dst` 切片中，两个切片的底层数据完全独立。
2. 修改 `dst` 的元素不会影响到 `src`，因为 `dst` 中的数据是新拷贝的副本。

## 九、映射map

Go 语言中，`map` 是一种内置的数据结构，用于实现键值对（key-value）存储。它可以高效地通过键快速查找、插入、和删除数据。`map` 是无序的，键的顺序在迭代时不保证固定。

### 基本语法与使用

1. **声明和初始化**

   使用 `make` 函数初始化一个 `map`：

   ```go
   m := make(map[string]int) // 键是字符串类型，值是整数类型
   ```

   或者在声明的同时初始化：

   ```go
   m := map[string]int{"apple": 1, "banana": 2, "cherry": 3}
   ```

2. **基本操作**

   - **添加或更新键值对**：

     ```go
     m["apple"] = 5 // 更新或添加键 "apple" 的值为 5
     ```

   - **访问元素**：

     ```go
     value := m["apple"] // 获取键 "apple" 对应的值，若不存在则返回类型的零值
     ```

   - **删除键值对**：

     使用 `delete` 函数从 `map` 中移除指定的键：

     ```go
     delete(m, "banana") // 删除键 "banana" 对应的键值对
     ```

   - **检测键是否存在**：

     在访问 `map` 中的某个键时，可以使用双重赋值的形式来检查键是否存在：

     ```go
     value, ok := m["apple"]
     if ok {
         fmt.Println("键存在，值为:", value)
     } else {
         fmt.Println("键不存在")
     }
     ```

### `map` 的底层实现原理

Go 的 `map` 是基于哈希表实现的，通过哈希函数将键映射到特定的位置。每个位置可以存放一个或多个键值对，以应对哈希冲突。

1. **哈希桶**：

   `map` 将键值对分组存储在多个**桶**（bucket）中，每个桶中可以容纳若干个键值对。

2. **哈希冲突**：

   由于不同的键可能会被映射到相同的位置，Go 使用链表或开放寻址来解决冲突。在某些情况下，Go 会将冲突较多的桶再分成小桶进行分散存储。

3. **动态扩容**：

   当 `map` 中的键值对增加较多时，哈希表会自动扩容并重新分布桶，以保证查询和插入的性能。

### 特性与限制

1. **无序性**：

   `map` 中的元素是无序的，迭代时每次的顺序可能会不同，Go 不保证 `map` 中的键的顺序。

2. **并发安全性**：

   `map` 不是并发安全的，在多线程环境下，多个 goroutine 同时读写 `map` 会导致竞态条件，通常需要使用 `sync.RWMutex` 进行加锁操作来保证安全。

3. **键的类型**：

   `map` 的键可以是任何可比较的类型，比如 `int`、`string`、`struct`（无切片、映射字段），但不能使用切片、映射、函数等不可比较的类型作为键。

### 示例代码

以下代码展示了 `map` 的基本使用：

```go
package main

import "fmt"

func main() {
    // 初始化一个 map
    m := map[string]int{"apple": 1, "banana": 2, "cherry": 3}

    // 添加或更新元素
    m["date"] = 4
    m["apple"] = 5

    // 检查某个键是否存在
    value, ok := m["banana"]
    if ok {
        fmt.Println("banana 存在，值为:", value)
    } else {
        fmt.Println("banana 不存在")
    }

    // 删除元素
    delete(m, "cherry")

    // 遍历 map
    for key, value := range m {
        fmt.Printf("%s -> %d\n", key, value)
    }
}
```

### `map` 的零值和空 map

- 一个未初始化的 `map` 的零值是 `nil`，操作一个 `nil` `map` 不会引发错误，但插入操作无效。

  ```go
  var m map[string]int
  fmt.Println(m == nil) // 输出: true
  ```

### `map` 的使用场景

`map` 非常适合需要快速查找的数据结构，如：

- 统计字符或单词的出现频率
- 记录对象的映射关系（如用户 ID 到用户信息）
- 缓存查询结果以提高效率

## 十、面向对象

### 1、结构体

#### 1. 结构体的定义与使用

结构体是用于将不同类型的数据字段组合在一起的一种数据类型，类似于其他语言的类。每个字段可以是不同的数据类型，结构体的实例可以包含所有字段的数据。

##### 定义结构体

```go
type Person struct {
    Name string
    Age  int
    Job  string
}
```

在上面的例子中，`Person` 是一个结构体类型，它有 `Name`、`Age` 和 `Job` 三个字段。

##### 创建结构体实例

可以使用以下方法创建结构体的实例：

```go
// 方式一：字段赋值
p1 := Person{Name: "Alice", Age: 30, Job: "Engineer"}

// 方式二：创建指针
p2 := &Person{Name: "Bob", Age: 25, Job: "Designer"}

// 方式三：省略字段名（不推荐）
// 这种方式取决于字段的顺序，容易出错
p3 := Person{"Charlie", 35, "Manager"}
```

#### 2. 访问和修改结构体字段

结构体的字段通过点运算符来访问和修改：

```go
fmt.Println(p1.Name) // 输出：Alice
p1.Age = 31          // 修改年龄
fmt.Println(p1.Age)  // 输出：31
```

#### 3. 方法绑定与接收者

Go 没有类的概念，因此将方法绑定到结构体上来实现对象的行为。方法的接收者（receiver）定义了该方法的所属结构体。

##### 定义方法

可以将一个方法绑定到结构体上，类似于类中的方法：

```go
func (p Person) Greet() string {
    return "Hello, my name is " + p.Name
}
```

在这里，`(p Person)` 是接收者，表示该方法是 `Person` 结构体的一个方法。`Greet` 方法返回包含姓名的问候语。

##### 使用方法

```go
fmt.Println(p1.Greet()) // 输出: Hello, my name is Alice
```

##### 指针接收者 vs 值接收者

- **值接收者**：方法接收结构体的副本，对副本的修改不会影响原始结构体。
- **指针接收者**：方法接收结构体的指针，指针接收者可以修改结构体的原始数据。

通常，在需要修改结构体内部数据或提高性能时使用指针接收者：

```go
func (p *Person) UpdateJob(newJob string) {
    p.Job = newJob
}
p1.UpdateJob("Artist")
fmt.Println(p1.Job) // 输出: Artist
```

#### 4. 构造函数

Go 中没有类构造函数，但可以定义一个函数来初始化和返回结构体实例，作为构造函数的替代。

```go
func NewPerson(name string, age int, job string) Person {
    return Person{Name: name, Age: age, Job: job}
}

p4 := NewPerson("David", 28, "Developer")
```

这里 `NewPerson` 函数就相当于一个构造函数，用于创建 `Person` 结构体实例。

#### 5. 组合与继承

Go 不支持传统的继承，但通过**组合**（Composition）来实现类似的功能。结构体可以包含其他结构体，以共享字段和方法。

##### 结构体嵌套

```go
type Address struct {
    City    string
    Country string
}

type Employee struct {
    Person
    Address
    Position string
}
```

在这里，`Employee` 结构体包含 `Person` 和 `Address`，可以访问这两个结构体的所有字段。

##### 使用嵌套结构体

```go
e := Employee{
    Person:   Person{Name: "Emma", Age: 26, Job: "Analyst"},
    Address:  Address{City: "New York", Country: "USA"},
    Position: "Data Analyst",
}
fmt.Println(e.Name) // 输出: Emma
fmt.Println(e.City) // 输出: New York
```

#### 6. 面向对象特性总结

Go 实现面向对象的特性主要依靠结构体和方法组合实现：

- **封装**：通过结构体定义属性，方法定义行为，将数据和操作封装在一个结构体中。
- **继承**：使用组合而非继承来复用代码，结构体可以嵌套其他结构体。
- **多态**：通过接口（`interface`）来实现多态，允许不同的结构体实现相同的接口方法。

#### 7. 内存布局与结构体

在内存中，结构体是一个连续的数据块，字段顺序与定义顺序一致。对于嵌套的结构体，内存布局也是按照组合的顺序来排列。

##### 示例

一个内存示例如下：

```go
type Person struct {
    Name string // 8字节 (假设)
    Age  int    // 8字节
}

type Employee struct {
    Person         // 嵌套的 Person，占用16字节
    Position string // 8字节
}
```

结构体 `Employee` 在内存中会以字段排列的顺序存储，`Person` 的字段在 `Employee` 的前半部分存储，紧接着存储 `Position` 字段。

### 2、方法

方法和函数的区别

### 3、跨包访问

大写字母开头无问题，解决小写字母结构体：工厂模式

### 4、封装

### 5、继承

注意事项：

### 6、接口

### 7、多态

细节：多态数组

### 8、断言

## 十一、文件操作

## 十二、协程和管道

### 1、协程

#### 1.1、主死从随

#### 1.2、锁

##### 1.2.1、互斥锁

##### 1.2.2、读写锁

### 2、管道

#### 2.1、基础

#### 2.2、遍历
