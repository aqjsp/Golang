# Golang入门

## 一、基本数据类型

### 1、整数类型

```go
var num3 int = 3
var num4 int8 = 4
var num5 int16 = 5
var num6 int32 = 6
var num7 int64 = 7
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

#### 4.1、定义



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

##### 5.3.2、new

##### 5.3.3、make

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

