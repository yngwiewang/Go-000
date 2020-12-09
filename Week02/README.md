学习笔记

1. 在 Go 中，`error` 仅仅是一个值
2. Go 内置的 `error` 是一个接口，实现了 `Error()` 方法就是实现了这个接口
3. `New()` 方法返回内部结构体的指针是为了避免自定义的 error 与预定义的 sentinel error 混淆
4. github.com/pkg/errors 包提供了丰富的功能，也包含了 Go 1.13 引入的 Wrap/Unwrap/Is/As 函数，完全兼容标准库 errors，推荐使用。
5. 写 python 的时候，标准库或者三方库会提供针对自己的 exception， 但是实际应用的时候不会去具体甄别而是笼统的 `except Exctption as e`，Go 中的错误定义和检查感觉非常清晰，减少很多心智负担，具体应用时根据场景，对于一个错误应该仅仅处理一次，要么吞掉要么向上抛，建议在顶层打印错误日志。
6. 以 database/sql 中的error 处理为例：
```
if driverErr, ok := err.(*mysql.MySQLError); ok {
	if driverErr.Number == mysqlerr.ER_ACCESS_DENIED_ERROR {
		// Handle the permission-denied error
	}
}
```