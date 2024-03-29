package logrus

import (
	"context"
	"fmt"
	"github.com/bitxx/logger/logbase"
	"os"

	"github.com/sirupsen/logrus"
)

type entryLogger interface {
	WithFields(fields logrus.Fields) *logrus.Entry
	WithError(err error) *logrus.Entry

	Log(level logrus.Level, args ...interface{})
	Logf(level logrus.Level, format string, args ...interface{})
}

type logrusLogger struct {
	Logger entryLogger
	opts   Options
}

func (l *logrusLogger) Init(opts ...logbase.Option) error {
	for _, o := range opts {
		o(&l.opts.Options)
	}

	if formatter, ok := l.opts.Context.Value(formatterKey{}).(logrus.Formatter); ok {
		l.opts.Formatter = formatter
	}
	if hs, ok := l.opts.Context.Value(hooksKey{}).(logrus.LevelHooks); ok {
		l.opts.Hooks = hs
	}
	if caller, ok := l.opts.Context.Value(reportCallerKey{}).(bool); ok && caller {
		l.opts.ReportCaller = caller
	}
	if exitFunction, ok := l.opts.Context.Value(exitKey{}).(func(int)); ok {
		l.opts.ExitFunc = exitFunction
	}

	switch ll := l.opts.Context.Value(logrusLoggerKey{}).(type) {
	case *logrus.Logger:
		// overwrite default options
		l.opts.Level = logrusToLoggerLevel(ll.GetLevel())
		l.opts.Out = ll.Out
		l.opts.Formatter = ll.Formatter
		l.opts.Hooks = ll.Hooks
		l.opts.ReportCaller = ll.ReportCaller
		l.opts.ExitFunc = ll.ExitFunc
		l.Logger = ll
	case *logrus.Entry:
		// overwrite default options
		el := ll.Logger
		l.opts.Level = logrusToLoggerLevel(el.GetLevel())
		l.opts.Out = el.Out
		l.opts.Formatter = el.Formatter
		l.opts.Hooks = el.Hooks
		l.opts.ReportCaller = el.ReportCaller
		l.opts.ExitFunc = el.ExitFunc
		l.Logger = ll
	case nil:
		log := logrus.New() // defaults
		log.SetLevel(loggerToLogrusLevel(l.opts.Level))
		log.SetOutput(l.opts.Out)
		log.SetFormatter(l.opts.Formatter)
		log.ReplaceHooks(l.opts.Hooks)
		log.SetReportCaller(l.opts.ReportCaller)
		log.ExitFunc = l.opts.ExitFunc
		l.Logger = log
	default:
		return fmt.Errorf("invalid logrus type: %T", ll)
	}

	return nil
}

func (l *logrusLogger) String() string {
	return "logrus"
}

func (l *logrusLogger) Fields(fields map[string]interface{}) logbase.Logger {
	return &logrusLogger{l.Logger.WithFields(fields), l.opts}
}

func (l *logrusLogger) Log(level logbase.Level, args ...interface{}) {
	l.Logger.Log(loggerToLogrusLevel(level), args...)
}

func (l *logrusLogger) Logf(level logbase.Level, format string, args ...interface{}) {
	l.Logger.Logf(loggerToLogrusLevel(level), format, args...)
}

func (l *logrusLogger) Options() logbase.Options {
	// FIXME: How to return full opts?
	return l.opts.Options
}

// New builds a new logger based on options
func NewLogger(opts ...logbase.Option) logbase.Logger {
	// Default options
	options := Options{
		Options: logbase.Options{
			Level:   logbase.InfoLevel,
			Fields:  make(map[string]interface{}),
			Out:     os.Stderr,
			Context: context.Background(),
		},
		Formatter:    new(logrus.TextFormatter),
		Hooks:        make(logrus.LevelHooks),
		ReportCaller: false,
		ExitFunc:     os.Exit,
	}
	l := &logrusLogger{opts: options}
	_ = l.Init(opts...)
	return l
}

func loggerToLogrusLevel(level logbase.Level) logrus.Level {
	switch level {
	case logbase.TraceLevel:
		return logrus.TraceLevel
	case logbase.DebugLevel:
		return logrus.DebugLevel
	case logbase.InfoLevel:
		return logrus.InfoLevel
	case logbase.WarnLevel:
		return logrus.WarnLevel
	case logbase.ErrorLevel:
		return logrus.ErrorLevel
	case logbase.FatalLevel:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

func logrusToLoggerLevel(level logrus.Level) logbase.Level {
	switch level {
	case logrus.TraceLevel:
		return logbase.TraceLevel
	case logrus.DebugLevel:
		return logbase.DebugLevel
	case logrus.InfoLevel:
		return logbase.InfoLevel
	case logrus.WarnLevel:
		return logbase.WarnLevel
	case logrus.ErrorLevel:
		return logbase.ErrorLevel
	case logrus.FatalLevel:
		return logbase.FatalLevel
	default:
		return logbase.InfoLevel
	}
}
