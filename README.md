# fslog
Golang日志包，基于zap二次封装，开箱即用
## 特性
- 支持多级别打印：Debug、Info、Warn、Error
- 支持日志文件切割
- 支持用户自定义zap二次开发
- 快捷的gin框架日志输出接入
- 运行时动态修改配置特性
- 日志多输出源支持
- 属性快照打印支持
## 打印
默认的可以直接打印字符串，如果需要模板字符串支持直接拼接参数即可（用法与fmt.Sprintf相同）。
如果需要查看当前环境下变量信息，则可以使用With方法，打印时信息中会包含With方法中的字段信息。
Error支持字符串、格式化字符串、error类型，如果时error类型，会自动打印堆栈信息。
```go
package main

import (
	"errors"
    "github.com/yushengji/fslog"
)

func main() {
    fslog.Debug("this is a debug log")
    fslog.Info("support %s log", "info")
    fslog.With("field", 1).Warn("this is a warn log")
    fslog.Error(errors.New("this is a error"))
}
```

## 性能
11th i7 16G Golang 1.22版本下，输出源仅为控制台： 32.5 ns/op 1 allocs/op