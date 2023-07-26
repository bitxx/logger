package logger

import (
	"github.com/jason-wj/logger/logbase"
	"github.com/jason-wj/logger/logbase/plugins/logger/logrus"
	"github.com/jason-wj/logger/logbase/plugins/logger/zap"
	"github.com/jason-wj/logger/logbase/writer"
	"github.com/jason-wj/logger/util"
	"io"
	"log"
	"os"
)

// SetupLogger 日志 cap 单位为kb
func SetupLogger(opts ...Option) logbase.Logger {
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
	var level logbase.Level
	level, err = logbase.GetLevel(op.level)
	if err != nil {
		log.Fatalf("get logger level error, %s", err.Error())
	}

	switch op.driver {
	case "zap":
		logbase.DefaultLogger, err = zap.NewLogger(logbase.WithLevel(level), zap.WithOutput(output), zap.WithCallerSkip(0))
		if err != nil {
			log.Fatalf("new zap logger error, %s", err.Error())
		}
	case "logrus":
		logbase.DefaultLogger = logrus.NewLogger(logbase.WithLevel(level), logbase.WithOutput(output), logrus.ReportCaller())
	default:
		logbase.DefaultLogger = logbase.NewLogger(logbase.WithLevel(level), logbase.WithOutput(output))
	}
	return logbase.DefaultLogger
}

func NewLogger(log logbase.Logger) *logbase.Helper {
	return logbase.NewHelper(log)
}
