# logger
支持zap、logrus、default(常规文件写入)三种日志写入模式，日志打印
满足日常项目开发需要

## 主要配置参数
type: zap、logrus、default
path：日志路径
level：日志等级，trace, debug, info, warn, error, fatal
stdout: file、default 其中file表示日志写入文件、default表示文件控制台展示
cap：文件写入日志条数

## 使用方式
```go
func main(){
	SetupLogger(
		WithType("zap"),
		WithPath("temp/logs"),
		WithLevel("info"),
		WithStdout("file"),
		WithCap(10),
	)

	log := NewLogger(core.DefaultLogger)
	log.Info("xxxxxx")
}
```
