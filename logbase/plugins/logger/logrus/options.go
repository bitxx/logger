package logrus

import (
	"github.com/jason-wj/logger/logbase"
	"github.com/sirupsen/logrus"
)

type Options struct {
	logbase.Options
	Formatter logrus.Formatter
	Hooks     logrus.LevelHooks
	// Flag for whether to log caller info (off by default)
	ReportCaller bool
	// Exit Function to call when FatalLevel log
	ExitFunc func(int)
}

type formatterKey struct{}

func WithTextTextFormatter(formatter *logrus.TextFormatter) logbase.Option {
	return logbase.SetOption(formatterKey{}, formatter)
}
func WithJSONFormatter(formatter *logrus.JSONFormatter) logbase.Option {
	return logbase.SetOption(formatterKey{}, formatter)
}

type hooksKey struct{}

func WithLevelHooks(hooks logrus.LevelHooks) logbase.Option {
	return logbase.SetOption(hooksKey{}, hooks)
}

type reportCallerKey struct{}

// warning to use this option. because logrus doest not open CallerDepth option
// this will only print this package
func ReportCaller() logbase.Option {
	return logbase.SetOption(reportCallerKey{}, true)
}

type exitKey struct{}

func WithExitFunc(exit func(int)) logbase.Option {
	return logbase.SetOption(exitKey{}, exit)
}

type logrusLoggerKey struct{}

func WithLogger(l logrus.StdLogger) logbase.Option {
	return logbase.SetOption(logrusLoggerKey{}, l)
}
