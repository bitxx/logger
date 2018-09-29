package main

import (
	"github.com/jason-wj/study/logger"
)

func initLogger(name, logPath, logName string, level string) (err error) {
	config := make(map[string]string, 8)
	config["log_path"] = logPath
	config["log_name"] = logName
	config["log_level"] = level
	config["log_split_type"] = "size"
	err = logger.InitLogger(name, config)
	if err != nil {
		return
	}

	logger.Debug("init logger success")
	return
}

//会快速生成日志，4秒钟可以生成100mb数据，用于测试按大小切割
func Run() {
	for {
		logger.Debug("user server is running user server is running user server is running user server is running user server is running user server is running user server is running user server is running user server is running user server is running user server is running user server is running")
		//time.Sleep(time.Second)
	}

}

func main() {
	initLogger("file", "./logs", "user_server", "debug")
	Run()
	return
}
