package logger

import (
	"fmt"
	"path"
	"runtime"
	"time"
)

type LogData struct {
	Message      string
	TimeStr      string
	LevelStr     string
	FileName     string
	FuncName     string
	LineNo       int
	WarnAndFatal bool
}

func getLineInfo() (filename string, funcName string, lineNo int) {
	//0表示获取下面这行代码执行时候的信息，比如这行代码返回line=7,file=util.go,funcName=getLineInfo()
	//1表示调用getLineInfo()这个方法所在地方的信息
	//2为调用包含getLineInfo()那个文件所在的方法的信息
	//3。。依次类推
	//可以看出这是一个栈，runtime.Caller
	//要根据层数来判断
	pc, file, line, ok := runtime.Caller(4) //pc表示当前程序执行的指令的位置

	if ok {
		filename = file
		funcName = runtime.FuncForPC(pc).Name() //根据pc指令位置找到对应的方法名
		lineNo = line
	}
	return
}

/*
当业务调用打日志的方法时，我们把数据写入到channel中（队列）
然后我们有一个后台的线程不断的从channel中获取这些日志数据，最终写入到文件
*/
func writeLog(level int, format string, args ...interface{}) *LogData {
	now := time.Now()
	nowStr := now.Format("2006-01-02 15:04:05:999")
	levelStr := getLevelText(level)
	filename, funcName, lineNo := getLineInfo()

	filename = path.Base(filename) //返回末尾路径
	funcName = path.Base(funcName)
	msg := fmt.Sprintf(format, args...) //args即使不输入也默认不为空。。。
	logData := &LogData{
		Message:      msg,
		TimeStr:      nowStr,
		LevelStr:     levelStr,
		FileName:     filename,
		FuncName:     funcName,
		LineNo:       lineNo,
		WarnAndFatal: false,
	}
	if level == LogLevelError || level == LogLevelWarn || level == LogLevelFatal {
		logData.WarnAndFatal = true
	}
	return logData

	//标准信息
	//fmt.Fprintf(file, "%s %s (%s:%s:%d) %s", nowStr, levelStr, filename, funcName, lineNo, msg)

	//fmt.Fprintln(file) //换行
}
