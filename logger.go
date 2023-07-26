package logger

import (
	"io"
	"log"
	"logger/core"
	"logger/core/plugins/logger/logrus"
	"logger/core/plugins/logger/zap"
	"logger/core/writer"
	"logger/util"
	"os"
)

// SetupLogger 日志 cap 单位为kb
func SetupLogger(opts ...Option) core.Logger {
	op := setDefault()
	for _, o := range opts {
		o(&op)
	}
	if !util.PathExist(op.path) {
		err := util.PathCreate(op.path)
		if err != nil {
			log.Fatalf("create dir error: %s", err.Error())
		}
	}
	var err error
	var output io.Writer
	switch op.stdout {
	case "file":
		output, err = writer.NewFileWriter(
			writer.WithPath(op.path),
			writer.WithCap(op.cap<<10),
		)
		if err != nil {
			log.Fatalf("logger setup error: %s", err.Error())
		}
	default:
		output = os.Stdout
	}
	var level core.Level
	level, err = core.GetLevel(op.level)
	if err != nil {
		log.Fatalf("get logger level error, %s", err.Error())
	}

	switch op.driver {
	case "zap":
		core.DefaultLogger, err = zap.NewLogger(core.WithLevel(level), zap.WithOutput(output), zap.WithCallerSkip(0))
		if err != nil {
			log.Fatalf("new zap logger error, %s", err.Error())
		}
	case "logrus":
		core.DefaultLogger = logrus.NewLogger(core.WithLevel(level), core.WithOutput(output), logrus.ReportCaller())
	default:
		core.DefaultLogger = core.NewLogger(core.WithLevel(level), core.WithOutput(output))
	}
	return core.DefaultLogger
}

func NewLogger(log core.Logger) *core.Helper {
	return core.NewHelper(log)
}
