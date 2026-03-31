# logger
支持zap、default(常规文件写入)三种日志写入模式，日志打印
满足日常项目开发需要

## 主要配置参数
* type: 日志类型，支持：zap、default
* path：日志路径
* level：日志等级，支持：trace, debug, info, warn, error, fatal
* stdout: 日志标准输出，支持：file、default `其中file表示日志写入文件、default表示文件控制台展示`
* cap：文件写入日志条数

## 使用方式
```go
func main(){
	log := NewLogger(
		WithType("zap"),
		WithPath("temp/logs"),
		WithLevel("info"),
		WithStdout("default"),
		WithCap(10),
	)

	log.Info("xxxxxx")
}
```

## 备注
可方便扩展自己需要的日志框架
