# GMP模型

掌握golang协程调度器的原理。

## 1、调度器的设计策略

### 1.1、复用线程

Go 调度器的一个重要设计策略是 **复用线程**，即尽可能复用操作系统的线程（OS Thread）。每个线程与处理器（Processor）紧密关联，调度器通过在多个处理器之间合理分配 Goroutines 来实现高效的并发执行。

#### 1.1.1、Work Stealing 机制

**Work Stealing（偷工作）机制**是一种优化线程负载均衡的方法，目的是解决某些处理器（Processor）由于任务完成得太快而闲置的情况。具体做法是：

- 如果某个处理器上的任务执行完了，且本地队列（Local Queue）没有待处理的 Goroutine，处理器会去其他处理器的队列中“偷取”任务。
- 这种机制通过增加并发的负载均衡，避免处理器因缺乏任务而处于空闲状态。

**工作原理：**

- Go 调度器使用每个处理器（Processor）上都有自己的队列，存储待执行的 Goroutine。
- 当一个处理器完成自己的任务后，它会查看其他处理器的队列，偷取任务来执行。
- 这种机制通过减少空闲 CPU 时间，来提高资源的利用率。

**示例：**

- 假设我们有 4 个 CPU 核心，每个核心都有自己的任务队列，核心 A 完成了它的任务，它就去核心 B 或核心 C 的队列中去“偷”任务，确保不会出现空闲处理器。

#### 1.1.2、Hand-Off 机制

**Hand-Off 机制**是另一种调度策略，用于当一个 Goroutine 被执行完成时，将其传递给其他处理器（Processor）执行。通常这会发生在 Goroutine 执行过程中涉及到阻塞操作时（例如 I/O 操作），此时它可能会将处理器资源“交接”给其他 Goroutine，以确保 CPU 不会因为某个 Goroutine 阻塞而浪费资源。

**工作原理：**

- 当 Goroutine 在执行时进入阻塞状态，或者其执行完毕后，调度器会将该 Goroutine 从当前的处理器队列中移除，并将其分配到另一个可用的处理器上。
- 这种方式可以避免某个处理器在 Goroutine 阻塞时被长时间占用，从而提高并发性。

**示例：**

- 当一个 Goroutine 执行 I/O 操作（如文件读取、网络请求）时，它可能会被调度器“交接”给其他处理器继续执行其他任务。

### 1.2、利用并行

Go 调度器允许并行执行 Goroutine，即在多核 CPU 上并行运行多个 Goroutines。调度器通过将多个处理器（Processor）与操作系统的线程（OS Thread）关联，在不同核心上执行 Goroutine。

**并行执行的实现：**

- 每个处理器（Processor）实际上是一个逻辑处理单元，调度器通过将多个处理器与操作系统的线程绑定，从而实现并行执行。
- 通过合理地将 Goroutines 分配到多个处理器上，Go 调度器可以充分利用多核 CPU 的性能优势。

**并行的优势：**

- Go 调度器能够根据系统的 CPU 核心数和负载情况动态调整 Goroutines 的调度策略，最大化 CPU 的利用率。
- 通过并行执行，Go 程序能够在多个 CPU 核心上同时处理多个任务，大大提高了程序的并发能力。

### 1.3、抢占

**抢占**是一种机制，指的是调度器可以在一个 Goroutine 正在运行时中断其执行并将 CPU 资源分配给其他任务。抢占在 Go 中是通过设置定时器（ticker）来实现的，调度器会定期检查 Goroutine 是否占用 CPU 时间过长，从而防止某个 Goroutine 长时间占用处理器而导致其他 Goroutine 被饿死。

**抢占机制的工作原理：**

Go 的调度器通过 协作式抢占 和 时间片抢占 两种方式来实现：

1. **协作式抢占**：Goroutine 本身在每次执行时会主动让出控制权，允许调度器中断它的执行。
2. **时间片抢占**：调度器通过设置定时器，定期检查 Goroutine 是否运行超过一定时间，如果超过时间限制，会强制中断并切换到其他 Goroutine 执行。

**示例：**

如果有两个 Goroutine 一直在运行，而一个 Goroutine 执行了很长时间（例如，进入了一个死循环），调度器会通过抢占机制来强制停止这个 Goroutine，转而执行其他任务。

### 1.4、全局 G 队列

全局 G 队列（Global Goroutine Queue）是调度器中存放未被分配给处理器的 Goroutine 的一个队列。全局 G 队列中存放的任务会被调度器从中取出，分配给空闲的处理器执行。

**工作原理：**

- **全局队列**中的 Goroutines 是那些没有被分配到任何处理器的 Goroutines。每当有空闲的处理器时，它会从全局队列中取出 Goroutine 并分配执行。
- 由于全局队列是一个共享队列，因此需要使用锁来防止并发操作时发生竞态条件。

**示例：**

- 在 Go 程序启动时，所有的 Goroutines 会被加入到全局队列。随着程序的执行，调度器会根据需要从全局队列中取出 Goroutine，分配给可用的处理器。

## 2、go fun() 经历了什么过程

### `go fun()` 过程详解

1. **创建 Goroutine**: 当我们调用 `go fun()` 时，Go 程序会在内部创建一个新的 **goroutine** 来执行函数 `fun()` 中的代码。每个 goroutine 都是由 Go 的运行时调度器（scheduler）管理，而不是操作系统线程。

   这个 goroutine 与主 goroutine（通常是 main 函数所在的 goroutine）并行执行，但 Go 的运行时系统会决定它们的执行顺序。主 goroutine 和新的 goroutine 可以在同一时刻运行，也可能轮流执行，具体取决于 Go 的调度器。

2. **Go 调度器的参与**: Go 的运行时调度器（scheduler）负责管理和调度 goroutines。每个 goroutine 都会被映射到某个操作系统线程上，调度器会根据运行时的负载、CPU 核心数等因素来分配资源。

   调度器采用 **M:N 调度模型**（M个 goroutines 映射到 N个操作系统线程）。Go 程序会启动多个 **Processor**，每个 Processor 都对应着一个操作系统线程。调度器的主要任务就是将 goroutines 分配到空闲的线程上执行。

3. **栈分配**: 每个 goroutine 都会有自己的 **栈**，用于存储局部变量和函数调用的上下文。与操作系统线程相比，goroutine 的栈非常小，通常为 2KB 左右，能够在需要时动态增长和收缩。这样，可以在同一个操作系统线程上同时运行成千上万个 goroutines。

4. **创建新栈和堆**: 在 goroutine 启动时，Go 会为其创建一个小的栈。随着 goroutine 执行，栈空间可能会增长（通常是以2倍增长的策略）。而堆则用于分配那些需要跨函数调用或多次使用的内存对象。

5. **协作式调度**: Go 使用 **协作式调度**，这意味着 goroutine 会在执行过程中主动让出控制权（例如调用 `runtime.Gosched()`，或者 goroutine 自己的执行达到一定条件时，调度器就会进行上下文切换）。但是 Go 也有抢占式调度机制，它会定期检查运行中的 goroutine 是否已运行超过了预定的时间片，如果超过，调度器就会强制中断并让其他 goroutine 执行。

6. **执行函数 `fun()`**: 一旦新的 goroutine 被创建并且分配了操作系统线程，Go 调度器会开始执行 `fun()` 函数的代码。在这个过程中，`fun()` 内部的所有局部变量和函数调用都会被存储在该 goroutine 的栈中，栈空间足够大时就可以执行下去。

7. **阻塞与唤醒**: 在 goroutine 执行过程中，如果遇到阻塞（如 I/O 操作、等待 channel、等待锁等），该 goroutine 会被挂起。调度器会将 CPU 资源分配给其他 goroutines 执行，直到该 goroutine 被唤醒。例如，如果某个 goroutine 正在等待 channel 上的数据，它就会被挂起，直到另一个 goroutine 向该 channel 发送数据，它才会恢复执行。

   Go 的调度器会确保不会因为某个 goroutine 被阻塞而让整个程序停滞不前。通过调度和抢占机制，Go 调度器能够实现高度并发的执行。

8. **执行完毕与垃圾回收**: 当 goroutine 中的代码执行完毕后，该 goroutine 会终止。若此时该 goroutine 的栈空间没有被其他 goroutine 再使用，调度器会将其回收，释放资源。

   Go 的垃圾回收器（GC）会自动管理内存的分配和回收。当一个 goroutine 结束时，相关的内存和栈空间会被清理。

9. **主 goroutine 等待其他 goroutines 完成**: 由于 `go fun()` 启动的是一个并发的 goroutine，所以主 goroutine 可能会继续执行其后续的代码，而新的 goroutine 在后台并行执行。如果主 goroutine 完成执行时，其他 goroutine 还没有执行完，主 goroutine 会退出，导致整个程序结束。为了确保主 goroutine 等待其他 goroutines 完成，可以使用 `sync.WaitGroup` 或 `channel` 等同步机制来实现。

### 详细流程图：

```go
main goroutine -> go fun() -> 创建新 goroutine
                           |
                           |--- 分配栈、内存、操作系统线程
                           |
                           v
                   Go 调度器调度
                    /        \
               执行 fun()    阻塞（等待事件，如 I/O、channel 等）
                    |        |
                执行完毕  恢复执行
                    |        |
                   退出    清理资源
```

### 示例代码分析：

```go
package main

import (
	"fmt"
	"time"
)

func fun() {
	fmt.Println("Start fun")
	time.Sleep(2 * time.Second)
	fmt.Println("End fun")
}

func main() {
	go fun()  // 使用 go 关键字启动 goroutine
	fmt.Println("Main goroutine")
	time.Sleep(3 * time.Second) // 等待 goroutine 执行完毕
}
```

1. **main goroutine** 启动后会执行 `go fun()`，创建一个新的 goroutine 执行 `fun()`。
2. 在 `go fun()` 启动后，`main` goroutine 会继续执行下一行代码：`fmt.Println("Main goroutine")`，并不会等待 `fun()` 完成。
3. `fun()` 会启动后执行打印 `"Start fun"`，然后进入 `time.Sleep(2 * time.Second)`，模拟耗时操作，等待 2 秒。
4. `main` goroutine 会等待 3 秒，以确保 `fun()` 执行完成（这里使用 `time.Sleep` 只是为了演示，实际应用中可以使用更合适的同步机制，如 `sync.WaitGroup`）。
5. 2 秒后，`fun()` 继续执行并打印 `"End fun"`，并最终退出。

### 为什么 `go fun()` 会并行执行？

- **Go 调度器**会自动将新创建的 goroutine 分配到空闲的操作系统线程上，多个 goroutine 会共享操作系统的线程，调度器会根据实际情况决定什么时候执行哪个 goroutine。
- **Goroutines 是协作式的**，这意味着它们会在执行过程中主动让出 CPU 给其他 goroutine。通过这种方式，Go 可以在有限的线程上执行大量的 goroutines，充分利用系统的并发能力。

## 3、调度器的生命周期

在 Go 语言中，调度器（Scheduler）负责管理所有 goroutine 的调度和执行。调度器的设计采用了 **M:N** 调度模型，其中 **M** 表示操作系统线程，**N** 表示 goroutine。Go 运行时（runtime）通过调度器将多个 goroutine 映射到较少的操作系统线程上执行。调度器的生命周期管理了从 goroutine 创建、执行到退出的整个过程。

在 Go 调度器中，**M0** 和 **G0** 是两个核心的概念，它们分别代表调度器的初始状态和与调度器生命周期相关的特定结构。

### 1. M（Machine）和 G（Goroutine）

在 Go 中，调度器使用 **M** 和 **G** 作为核心的结构体来管理 goroutine 和操作系统线程的关系。

- **M**: 代表一个操作系统线程（Machine），即 Go 运行时所使用的线程。每个 M 对象与一个操作系统线程绑定，用于执行实际的代码。
- **G**: 代表一个 goroutine（Goroutine）。每个 G 是由调度器调度的单位，调度器将 G 映射到 M 上执行。一个 M 上可以同时执行多个 G。

### 2. M0 和 G0：调度器的初始状态

在 Go 的调度器中，**M0** 和 **G0** 是调度器启动时的关键结构，代表调度器内部的特殊状态。理解 M0 和 G0 的生命周期，能够帮助我们更好地理解调度器的初始化和执行过程。

#### M0（调度器的初始线程）

- **M0** 是 Go 运行时在启动时创建的一个特殊线程（线程 0），是调度器在系统启动时的起始点。它是操作系统线程的起点，用于初始化 Go 的调度器并启动第一个 goroutine。
- **M0** 负责初始化调度器的状态，包括创建 goroutine 队列、启动其他 M 线程等。M0 会调用调度器的核心函数，最终进入 Go 的调度循环。

#### G0（调度器的初始 goroutine）

- **G0** 是调度器内部专用的 goroutine，它并不执行用户代码，而是用于执行调度器的初始化任务。
- **G0** 的主要任务是启动 Go 运行时的调度工作，执行调度器的初始化任务，分配和管理调度工作。它并不参与实际的用户代码执行。
- 在运行时系统启动时，`G0` 会初始化调度器，并为每个操作系统线程创建一个对应的 `M`，并且调度和执行实际的 goroutine。

### 3. 调度器的生命周期和调度流程

调度器的生命周期涉及 M 和 G 的创建、调度、运行、挂起和销毁。我们从调度器初始化开始，详细讲解 M 和 G 的生命周期管理。

#### 调度器的初始化

- 当 Go 程序启动时，调度器首先会创建一个 `M0`（即操作系统线程）和一个 `G0`（调度器的初始化 goroutine）。`G0` 用于初始化 Go 运行时的调度器并为所有操作系统线程准备好资源。
- 随着调度器的初始化完成，多个 M 线程会被创建，并且每个 M 线程开始从全局队列中获取并执行 goroutine。

#### Goroutine 的调度

- 当调用 `go` 关键字启动一个新的 goroutine 时，调度器会将该 goroutine（G）放入某个 M 线程的队列中，等待调度。
- 每个 goroutine 都包含了其执行的上下文信息（如堆栈、指令指针等），并且可以通过调度器选择一个空闲的操作系统线程（M）来执行它。
- 由于 Go 的调度器采用了 M:N 的模型，多个 goroutine 可能会被调度到同一个操作系统线程上执行。

#### 抢占和阻塞

- Go 调度器会周期性地检查每个 goroutine 是否已经运行了足够长时间，如果一个 goroutine 在运行过程中消耗了过多时间，调度器会将其挂起，释放该线程，去调度其他 goroutine。
- 如果某个 goroutine 因为 I/O 等原因阻塞，调度器会将它从正在执行的线程中移除，并选择其他 goroutine 来执行。

#### Goroutine 完成和销毁

- 当一个 goroutine 执行完毕时，它会被从调度队列中移除，相关的栈空间也会被释放。此时，goroutine 的生命周期结束。
- 如果某个 M 上的 goroutines 都执行完毕，M 线程会进入空闲状态，并等待其他任务的调度。

#### 调度器的抢占机制

- Go 的调度器采用的是 **抢占式调度**。当一个 goroutine 长时间运行时，调度器会强制中断它，切换到其他 goroutine 执行。通过这种方式，Go 调度器避免了某个 goroutine 永远占用 CPU 的问题。
- 在 Go 1.14 版本及之后，调度器对 goroutine 的抢占性更加明显，尤其在多核处理器上能够更好地利用资源。

### 4. M 和 G 的关系及调度

Go 的调度器通过 M 和 G 管理并发执行的 goroutines。调度器确保每个 M 线程都能根据负载和需要分配合适的 G 进行执行。

- Goroutine 绑定到 M 上: 每个 G 都会被调度到一个 M 上执行。当 M 执行 G 的代码时，M 会持有该 G 的控制权。一个 M 线程在执行 goroutine 时，负责调度和执行 G 的代码。
- 调度器的队列: Go 的调度器通常会使用两种队列来管理 G 的调度：
  - 本地队列：每个 M 会有自己的本地队列，负责存储该 M 执行的 goroutine。
  - 全局队列：所有 M 共享的队列，存储所有未被执行的 goroutines。如果某个 M 的本地队列为空，调度器就会从全局队列中取出一个 goroutine 来执行。

#### M 和 G 的生命周期

- **创建 M 和 G**: M 是操作系统线程，而 G 是 goroutine。当 Go 程序启动时，会初始化一个 M0 线程和一个 G0 goroutine。随着 goroutine 的创建，M 线程会通过调度器将 G 分配到 M 上执行。
- **执行 G**: M 线程通过调度器不断获取并执行 G，直到所有 G 执行完成。
- **销毁 M 和 G**: 当所有的 G 完成时，相应的 M 线程也会退出，系统会清理所有资源，结束程序执行。

## 4、可视化的GMP编程

### 1. **GMP 模型与 Trace 工具**

GMP 模型背后的核心原理是通过调度器（scheduler）来高效地管理并发任务。Go 运行时调度器通过多个 **P**（逻辑处理器）管理多个 **G**（goroutine），并且将这些 goroutine 映射到多个 **M**（操作系统线程）上执行。为了能够详细了解 Go 程序如何在并发和并行的环境中运行，Go 提供了 **trace** 工具，它可以捕获和记录 Go 程序执行的运行时事件，并将这些事件以时间轴的方式呈现出来。

#### **GMP 模型与调度过程**

- **G**（Goroutine）：代表程序中的每个 goroutine。每个 G 都会被调度到某个 M 上执行。
- **M**（Machine）：操作系统线程，负责执行实际的代码。
- **P**（Processor）：逻辑处理器，负责调度 goroutines 到空闲的 M 上执行。每个 P 都有一个本地队列来管理待执行的 goroutines。

### 2. **GODEBUG 环境变量**

Go 语言的 **GODEBUG** 环境变量允许开发者查看并调节 Go 调度器的行为。通过设置 **GODEBUG**，可以获取调度器的详细信息，包括并发调度的状态、goroutine 的调度过程、工作窃取（Work Stealing）机制等。

#### **常用的 GODEBUG 配置选项**：

- **schedtrace=1000**：每秒输出一次调度器的信息，包括 P 的数量、G 的数量、在运行的 M 的数量等。
- **gctrace=1**：输出垃圾回收的统计信息。
- **net/http/pprof=1**：开启 pprof 工具，用于分析程序的性能瓶颈。

例如，在启动 Go 程序时，你可以设置以下命令来调试并查看调度器的行为：

```go
GODEBUG=schedtrace=1000 go run main.go
```

这将每秒打印一次调度器的详细信息。

#### 示例：

```go
$ GODEBUG=schedtrace=1000 go run main.go
```

输出示例：

```go
sched: 0.02s CPU, 0.01s runtime, 0.03s idle, 0.00s poll, 0.00s sys, 0.00s gc, 0.00s block
```

这表示调度器在运行时的 CPU 时间、空闲时间、系统调用时间、垃圾回收时间等。

### 3. **Go Trace 工具**

Go Trace 工具可以生成程序执行期间的详细 trace 文件，从中你可以查看到程序执行的具体细节，包括 goroutine 的调度过程、系统调用、锁等待等事件。生成的 trace 文件以图形化的方式展现，可以通过浏览器查看，帮助开发者深入分析程序的性能。

#### **生成 trace 文件**：

你可以通过调用 `runtime/trace` 包来启用 Go 程序的 trace 功能。以下是一个简单的例子：

```go
package main

import (
    "fmt"
    "runtime/trace"
    "os"
)

func main() {
    // 创建一个 trace 文件
    f, err := os.Create("trace.out")
    if err != nil {
        fmt.Println("could not create trace file:", err)
        return
    }
    defer f.Close()

    // 启动 trace
    err = trace.Start(f)
    if err != nil {
        fmt.Println("could not start trace:", err)
        return
    }
    defer trace.Stop()

    // 程序的实际逻辑
    fmt.Println("Start Goroutine")
    go func() {
        fmt.Println("Hello from Goroutine!")
    }()
    fmt.Println("Main Goroutine")
}
```

- 通过 `trace.Start()` 开启 trace 功能，并将生成的 trace 数据写入到指定的文件（如 `trace.out`）中。
- 执行程序后，trace 文件将包含详细的调度和执行信息。

#### **分析 Trace 文件**：

1. 在命令行中运行 Go 程序，它会生成一个名为 `trace.out` 的文件。
2. 使用 `go tool trace` 工具分析这个文件：

```go
go tool trace trace.out
```

1. 这将启动一个 Web 服务，并在浏览器中展示程序的执行情况，包括 goroutine 的调度情况、内存分配、锁的等待、系统调用等。

### 4. **Trace 输出的内容**

Go trace 文件的内容包括多个部分，帮助开发者更好地理解程序在执行时的具体行为。常见的 trace 输出包括：

#### **Goroutine 分析**

- **Goroutine 状态**：可以看到 goroutine 的创建、调度和阻塞状态，帮助分析程序的并发性能。
- **调度事件**：每个 goroutine 被调度到某个 M 上执行的时间节点，显示了调度的延迟、切换等信息。
- **阻塞事件**：如等待锁、IO 操作等，帮助开发者识别程序中的性能瓶颈。

#### **CPU 和内存使用分析**

- **CPU 使用情况**：展示了程序在不同时间段的 CPU 时间使用情况，帮助开发者了解程序的计算密集型或 IO 密集型操作。
- **内存分配**：显示了程序的内存分配情况，帮助分析内存泄漏和垃圾回收的性能。

#### **Lock 分析**

- **锁竞争**：显示了程序中的锁竞争情况，包括等待锁的 goroutine 和锁的释放时间。

#### **其他调度相关信息**

- **全局队列和本地队列**：显示了各个 P 的本地队列，以及全局队列的状态，帮助开发者分析工作窃取和负载均衡问题。
- **系统调用和调度延迟**：展示了系统调用的延迟和调度延迟，帮助识别性能瓶颈。

## 5、场景

5.1、场景1

P拥有G1，M1获取P后开始运行G1，G1使用go fun()创建了G2，为了局部性G2有线加入到P1的本地队列。

### 1. 简单的 Goroutine 调度：多核 CPU 上的负载均衡

在 Go 中，Goroutines 是并发执行的单位，而调度器的工作是将这些 Goroutines 分配到可用的 M（操作系统线程）上执行。在多核 CPU 上，Go 的调度器会创建多个 P（处理器），每个 P 可以拥有自己的本地队列，用于存放 Goroutines。

#### 场景：

假设有一个任务需要并发执行，Go 会创建多个 Goroutines，然后将这些 Goroutines 分配到不同的 P 上执行。如果有多个 P（假设是 4 个 P），而 M 的数量足够多（比如 8 个 M），调度器会动态调整每个 P 上的 Goroutines 数量，并且通过调度机制确保每个 CPU 核心都在高效地工作。

**代码示例**：

```go
package main

import (
    "fmt"
    "time"
)

func work(id int) {
    fmt.Printf("Goroutine %d is working\n", id)
    time.Sleep(time.Second)
}

func main() {
    for i := 0; i < 10; i++ {
        go work(i)
    }

    time.Sleep(5 * time.Second) // 等待所有 Goroutines 完成
}
```

#### GMP 模型中的作用：

- **G**（Goroutines）是执行的基本单位，每个 Goroutine 会通过调度器分配给一个 **M**（操作系统线程）来执行。
- **P**（Processor）是一个逻辑处理器，负责调度 Goroutines 到空闲的 **M** 上执行。在多核机器上，Go 会根据可用的 CPU 核心数量自动创建多个 P。

在这个简单的例子中，假设有 10 个 Goroutines，调度器会将它们分配到 4 个 P 上执行，这样可以最大化地利用 4 个 CPU 核心。

### 2. Goroutine 阻塞与调度的切换

在 Go 中，当某个 Goroutine 发生阻塞时（例如等待 I/O 操作完成，或者在锁上等待），调度器会将其从当前 P 上移除，转而执行队列中的其他 Goroutines。为了确保高效的并发执行，Go 的调度器设计了工作窃取（Work Stealing）机制，使得空闲的 P 可以从其他 P 中窃取任务来继续执行，从而避免了某些 P 长时间空闲。

#### 场景：

假设程序中有多个 goroutine，在执行过程中某些 goroutine 发生了阻塞，调度器会把它们移除当前的处理器队列，然后让其他空闲的处理器（P）继续从全局队列或其他处理器的队列中窃取任务，保证了 CPU 不会空闲。

**代码示例**：

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func work(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Printf("Goroutine %d starting\n", id)
    time.Sleep(2 * time.Second)  // 模拟阻塞
    fmt.Printf("Goroutine %d done\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 0; i < 5; i++ {
        wg.Add(1)
        go work(i, &wg)
    }

    wg.Wait() // 等待所有 Goroutines 完成
}
```

#### GMP 模型中的作用：

- 由于每个 P 有一个本地队列来存放 Goroutines，当某个 Goroutine 阻塞时，调度器会将其移出队列并选择另一个可执行的任务。
- 如果某个 P 上的队列空闲了，它会尝试从其他 P 的队列中窃取任务。

### 3. 工作窃取（Work Stealing）机制

Go 的调度器采用了工作窃取机制，即当某个 P 的本地队列为空时，它可以从其他 P 的队列中窃取任务。这一机制保证了即使在某些 P 阻塞的情况下，其他空闲的 P 也能够继续执行任务，从而提高并发性能。

#### 场景：

当多个 Goroutines 在不同的 CPU 核心上运行时，如果某个核心的 Goroutines 完成了它的任务，调度器会让这个核心去“窃取”其他核心的任务，以便保持 CPU 的高效利用。

**代码示例**：

```go
package main

import (
    "fmt"
    "sync"
)

func work(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Printf("Goroutine %d working\n", id)
}

func main() {
    var wg sync.WaitGroup

    // 启动多个 Goroutines
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go work(i, &wg)
    }

    // 等待所有 Goroutines 完成
    wg.Wait()
}
```

#### GMP 模型中的作用：

- 每个 P 有一个本地队列来存放待处理的 Goroutines。
- 如果某个 P 完成了它的任务且队列为空，它将尝试从其他 P 的队列中窃取任务。
- 这种工作窃取机制保证了 Goroutines 的负载均衡，减少了某些处理器长时间空闲的情况。

### 4. 调度器的并发与并行控制

Go 通过调度器控制程序的并发性和并行性。在多核系统上，Go 默认会根据可用的 CPU 核心数量创建多个 P，从而实现并行执行。你也可以通过设置 **GOMAXPROCS** 来控制调度器使用的最大核心数。

#### 场景：

假设程序有多个 Goroutines，并且运行在多核机器上，Go 的调度器会根据系统的核心数来分配 P 数量，从而实现多核并行执行。通过设置 **GOMAXPROCS**，你可以控制程序使用多少个处理器。

**代码示例**：

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
)

func work(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Printf("Goroutine %d working\n", id)
}

func main() {
    runtime.GOMAXPROCS(4) // 设置使用 4 个 CPU 核心
    var wg sync.WaitGroup

    for i := 0; i < 10; i++ {
        wg.Add(1)
        go work(i, &wg)
    }

    wg.Wait() // 等待所有 Goroutines 完成
}
```

#### GMP 模型中的作用：

- **P**（逻辑处理器）控制了并发调度的能力，多个 P 可以并行执行多个 Goroutines。
- **GOMAXPROCS** 用于控制 Go 程序使用的核心数，从而影响并行程度。

### 5. Goroutine 执行与系统调用（IO 密集型）

在 IO 密集型操作中，Go 会将 Goroutine 调度到其他可用的 M（操作系统线程）上执行，避免单一 M 阻塞而导致其他任务无法执行。Go 调度器会智能地调度 Goroutine，在需要进行系统调用（如网络请求、文件操作等）时，将其挂起并允许其他任务继续执行。

#### 场景：

假设你的程序包含大量的网络请求操作，Go 的调度器会把每个等待响应的 Goroutine 挂起，而将其他不依赖于 I/O 的 Goroutines 调度到可用的 M 上执行，从而最大化 CPU 的利用率。

**代码示例**：

```go
package main

import (
    "fmt"
    "net/http"
    "sync"
)

func fetchURL(url string, wg *sync.WaitGroup) {
    defer wg.Done()
    _, err := http.Get(url)
    if err != nil {
        fmt.Println("Error fetching", url)
        return
    }
    fmt.Println("Fetched:", url)
}

func main() {
    var wg sync.WaitGroup
    urls := []string{
        "http://example.com",
        "http://google.com",
        "http://golang.org",
    }

    for _, url := range urls {
        wg.Add(1)
        go fetchURL(url, &wg)
    }

    wg.Wait() // 等待所有的 Goroutines 完成
}
```

#### GMP 模型中的作用：

- 在 IO 密集型任务中，调度器会使用更多的 M 来处理更多的 Goroutines，使得程序能够有效利用 CPU。
- 当 Goroutine 被阻塞（如等待 HTTP 响应）时，它不会占用 M，而是将 M 释放出来去运行其他任务。