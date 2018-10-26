package logger

import (
	"github.com/jason-wj/AIChain-blockchain-prime/aichain-common/logger/hook"
	"github.com/sirupsen/logrus"
	"os"
)

var logger = logrus.New()

// 封装logrus.Fields
type Fields logrus.Fields

var hk *hook.Hook

func Init() {
	SetLogFormatter(&logrus.JSONFormatter{})
	SetLogLevel(logrus.InfoLevel)
	hk = hook.NewHook("./log.log")
	logger.Out = os.Stdout
}

func SetLogLevel(level logrus.Level) {
	logger.Level = level
}
func SetLogFormatter(formatter logrus.Formatter) {
	logger.Formatter = formatter
}

// Debug
func Debug(args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		logger.WithFields(logrus.Fields{}).Debug(args)
		logger.AddHook(hk)
	}
}

// 带有field的Debug
func DebugWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.DebugLevel {
		logger.WithFields(logrus.Fields(f)).Debug(l)
		logger.AddHook(hk)
	}
}

// Info
func Info(args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		logger.WithFields(logrus.Fields{}).Info(args...)
		logger.AddHook(hk)
	}
}

// 带有field的Info
func InfoWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.InfoLevel {
		logger.WithFields(logrus.Fields(f)).Info(l)
		logger.AddHook(hk)
	}
}

// Warn
func Warn(args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		logger.WithFields(logrus.Fields{}).Warn(args...)
		logger.AddHook(hk)
	}
}

// 带有Field的Warn
func WarnWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.WarnLevel {
		logger.WithFields(logrus.Fields(f)).Warn(l)
		logger.AddHook(hk)
	}
}

// Error
func Error(args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		logger.WithFields(logrus.Fields{}).Error(args...)
		logger.AddHook(hk)
	}
}

// 带有Fields的Error
func ErrorWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.ErrorLevel {
		logger.WithFields(logrus.Fields(f)).Error(l)
		logger.AddHook(hk)
	}
}

// Fatal
func Fatal(args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		logger.WithFields(logrus.Fields{}).Fatal(args...)
		logger.AddHook(hk)
	}
}

// 带有Field的Fatal
func FatalWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.FatalLevel {
		logger.WithFields(logrus.Fields(f)).Fatal(l)
		logger.AddHook(hk)
	}
}

// Panic
func Panic(args ...interface{}) {
	if logger.Level >= logrus.PanicLevel {
		logger.WithFields(logrus.Fields{}).Panic(args...)
		logger.AddHook(hk)
	}
}

// 带有Field的Panic
func PanicWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.PanicLevel {
		logger.WithFields(logrus.Fields(f)).Panic(l)
		logger.AddHook(hk)
	}
}
