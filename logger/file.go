package logger

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type FileLogger struct {
	level         int
	logPath       string
	logName       string
	file          *os.File      //存放info debug trace 级别日志
	warnFile      *os.File      //存放error warn 级别日志
	logDataChan   chan *LogData //管道
	logSplitType  int
	logSplitSize  int64 //int表示的值不够用，在此是kb级别，不是mb
	lastSplitHour int   //最新切割时间
}

func NewFileLogger(config map[string]string) (logger LogInterface, err error) {
	logPath, ok := config["log_path"]
	if !ok {
		err = fmt.Errorf("not found 'log_path'")
		return
	}

	logName, ok := config["log_name"]
	if !ok {
		err = fmt.Errorf("not found 'log_name'")
		return
	}

	logLevel, ok := config["log_level"]
	if !ok {
		err = fmt.Errorf("not found 'log_level'")
		return
	}

	logChanSize, ok := config["log_chain_size"]
	if !ok {
		logChanSize = "50000"
	}

	var logSplitType = LogSplitTypeHour //默认按小时切分
	var logSplitSize int64              //按大小切分，设置大小
	logSplitStr, ok := config["log_split_type"]
	if !ok {
		logSplitStr = "hour"
	} else {
		if logSplitStr == "size" {
			logSplitSizeStr, ok := config["log_split_size"]
			if !ok {
				logSplitSizeStr = "104857600" //100mb
			}
			logSplitSize, err = strconv.ParseInt(logSplitSizeStr, 10, 64)
			if err != nil {
				logSplitSize = 104857600
			}
			logSplitType = LogSplitTypeSize
		} else {
			logSplitType = LogSplitTypeHour
		}
	}

	chanSize, err := strconv.Atoi(logChanSize)
	if err != nil {
		chanSize = 50000
	}

	level := getLogLevel(logLevel)

	//level int, logPath string, logName string
	logger = &FileLogger{
		level:         level,
		logPath:       logPath,
		logName:       logName,
		logDataChan:   make(chan *LogData, chanSize),
		logSplitType:  logSplitType,
		logSplitSize:  logSplitSize,
		lastSplitHour: time.Now().Hour(),
	}
	logger.Init()
	return logger, nil
}

func (f *FileLogger) Init() {
	//写info、debug等的日志
	filename := fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open faile %s failed,err:%v", filename, err))
	}
	f.file = file

	//写错误日志和failt日志的文件
	filename = fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
	file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open faile %s failed,err:%v", filename, err))
	}
	f.warnFile = file

	go f.writeLogBackground()
}

//splitFileHour 按时间切分日志具体逻辑
func (f *FileLogger) splitFileHour(warnFile bool) {
	now := time.Now()
	hour := now.Hour()
	if hour == f.lastSplitHour {
		return
	}
	f.lastSplitHour = hour
	var backupFileName string
	var fileName string
	if warnFile {
		backupFileName = fmt.Sprintf("%s/%s.log.wf_%04d_%02d_%02d_%02d", f.logPath,
			f.logName, now.Year(), now.Month(), now.Day(), f.lastSplitHour)
		fileName = fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
	} else {
		backupFileName = fmt.Sprintf("%s/%s.log_%04d_%02d_%02d_%02d", f.logPath,
			f.logName, now.Year(), now.Month(), now.Day(), f.lastSplitHour)
		fileName = fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	}
	file := f.file
	if warnFile {
		file = f.warnFile
	}
	file.Close()
	os.Rename(fileName, backupFileName) //复制备份日志

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return
	}
	if warnFile {
		f.warnFile = file
	} else {
		f.file = file
	}
}

//splitFileSize 按大小切分日志具体逻辑
func (f *FileLogger) splitFileSize(warnFile bool) {
	file := f.file
	if warnFile {
		file = f.warnFile
	}
	info, err := file.Stat()
	if err != nil {
		return
	}
	fileSize := info.Size()
	//没达到切割空间大小，则停止
	if fileSize < f.logSplitSize {
		return
	}

	now := time.Now()
	var backupFileName string
	var fileName string
	if warnFile {
		backupFileName = fmt.Sprintf("%s/%s.log.wf_%04d_%02d_%02d_%02d_%02d_%02d", f.logPath,
			f.logName, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		fileName = fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
	} else {
		backupFileName = fmt.Sprintf("%s/%s.log_%04d_%02d_%02d_%02d_%02d_%02d", f.logPath,
			f.logName, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		fileName = fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	}

	file.Close()
	os.Rename(fileName, backupFileName) //复制备份日志，重命名

	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return
	}
	if warnFile {
		f.warnFile = file
	} else {
		f.file = file
	}
}

//检测切分逻辑，并进行切分
func (f *FileLogger) checkSplitFile(warnFile bool) {
	if f.logSplitType == LogSplitTypeHour {
		f.splitFileHour(warnFile)
		return
	}
	/*if f.logSplitType == LogSplitTypeSize {
		f.splitFileSize(warnFile)
		return
	}*/
	f.splitFileSize(warnFile)

}

//后台写日志
func (f *FileLogger) writeLogBackground() {
	for logData := range f.logDataChan {
		var file *os.File = f.file
		if logData.WarnAndFatal {
			file = f.warnFile
		}
		f.checkSplitFile(logData.WarnAndFatal)
		fmt.Fprintf(file, "%s %s (%s:%s:%d) %s\n", logData.TimeStr,
			logData.LevelStr, logData.FileName, logData.FuncName,
			logData.LineNo, logData.Message)
	}
}

func (f *FileLogger) SetLevel(level int) {
	if level <= LogLevelDebug || level > LogLevelFatal {
		level = LogLevelDebug
	}
	f.level = level
}

func (f *FileLogger) Debug(format string, args ...interface{}) {
	if f.level > LogLevelDebug {
		return
	}
	logData := writeLog(LogLevelDebug, format, args...)
	//判断channel中，是否满了，若没满，则写入，若是满了，则跳过，丢弃部分日志，不影响业务
	//使用select判断
	select {
	case f.logDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Trace(format string, args ...interface{}) {
	if f.level > LogLevelTrace {
		return
	}
	logData := writeLog(LogLevelTrace, format, args...)
	//判断channel中，是否满了，若没满，则写入，若是满了，则跳过，丢弃部分日志，不影响业务
	//使用select判断
	select {
	case f.logDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Info(format string, args ...interface{}) {
	if f.level > LogLevelInfo {
		return
	}
	logData := writeLog(LogLevelInfo, format, args...)
	//判断channel中，是否满了，若没满，则写入，若是满了，则跳过，丢弃部分日志，不影响业务
	//使用select判断
	select {
	case f.logDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Warn(format string, args ...interface{}) {
	if f.level > LogLevelWarn {
		return
	}
	logData := writeLog(LogLevelWarn, format, args...)
	//判断channel中，是否满了，若没满，则写入，若是满了，则跳过，丢弃部分日志，不影响业务
	//使用select判断
	select {
	case f.logDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Error(format string, args ...interface{}) {
	if f.level > LogLevelError {
		return
	}
	logData := writeLog(LogLevelError, format, args...)
	//判断channel中，是否满了，若没满，则写入，若是满了，则跳过，丢弃部分日志，不影响业务
	//使用select判断
	select {
	case f.logDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Fatal(format string, args ...interface{}) {
	if f.level > LogLevelFatal {
		return
	}
	logData := writeLog(LogLevelFatal, format, args...)
	//判断channel中，是否满了，若没满，则写入，若是满了，则跳过，丢弃部分日志，不影响业务
	//使用select判断
	select {
	case f.logDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Close() {
	f.file.Close()
	f.warnFile.Close()
}
