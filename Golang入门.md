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

