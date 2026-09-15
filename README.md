# Golang

面向 C++ 后端转 Go、或第一次认真用 Go 写服务的笔记。goroutine 不是「更轻的线程」，interface 不是虚函数表，GC 也不是「不用管内存」。

```
基础/         入门、切片与 map、函数与方法、接口
进阶/         GMP、channel、GC、context、逃逸分析
实战项目/     把前面的点串进一个能跑的 HTTP 服务
面试连环问/   校招 / 社招里真实出现过的 Go 追问
```

默认 Go 1.21+（泛型、`slog`、`loopvar` 语义）。GMP / GC 按当前官方实现讲，不把 1.3 的故事当现状。

## 基础

- [Go入门](基础/Go入门.md)
- [切片与Map](基础/切片与Map.md)
- [函数与方法](基础/函数与方法.md)
- [接口](基础/接口.md)

## 进阶

- [GMP调度](进阶/GMP调度.md)
- [Channel与并发](进阶/Channel与并发.md)
- [GC与内存](进阶/GC与内存.md)
- [Context与错误处理](进阶/Context与错误处理.md)
- [逃逸分析与性能](进阶/逃逸分析与性能.md)

## 实战项目

- [HTTP服务](实战项目/HTTP服务.md)

## 面试

- [Go面试连环问](面试连环问/Go面试连环问.md)
