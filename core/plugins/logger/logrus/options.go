package logrus

import (
	"github.com/jason-wj/logger/core"
	"github.com/sirupsen/logrus"
)

type Options struct {
	core.Options
	Formatter logrus.Formatter
	Hooks     logrus.LevelHooks
	// Flag for whether to log caller info (off by default)
	ReportCaller bool
	// Exit Function to call when FatalLevel log
	ExitFunc func(int)
}

type formatterKey struct{}

func WithTextTextFormatter(formatter *logrus.TextFormatter) core.Option {
	return core.SetOption(formatterKey{}, formatter)
}
func WithJSONFormatter(formatter *logrus.JSONFormatter) core.Option {
	return core.SetOption(formatterKey{}, formatter)
}

type hooksKey struct{}

func WithLevelHooks(hooks logrus.LevelHooks) core.Option {
	return core.SetOption(hooksKey{}, hooks)
}

type reportCallerKey struct{}

// warning to use this option. because logrus doest not open CallerDepth option
// this will only print this package
func ReportCaller() core.Option {
	return core.SetOption(reportCallerKey{}, true)
}

type exitKey struct{}

func WithExitFunc(exit func(int)) core.Option {
	return core.SetOption(exitKey{}, exit)
}

type logrusLoggerKey struct{}

func WithLogger(l logrus.StdLogger) core.Option {
	return core.SetOption(logrusLoggerKey{}, l)
}
