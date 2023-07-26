package logbase

import (
	"context"
	"io"
	"time"
)

type Option func(*Options)

type Options struct {
	// The logging level the logger should log at. default is `InfoLevel`
	Level Level
	// fields to always be logged
	Fields map[string]interface{}
	// It's common to set this to a file, or leave it default which is `os.Stderr`
	Out io.Writer
	// Caller skip frame count for file:line info
	CallerSkipCount int
	// Alternative options
	Context context.Context
	// Name logger name
	Name string
}

// WithFields set default fields for the logger
func WithFields(fields map[string]interface{}) Option {
	return func(args *Options) {
		args.Fields = fields
	}
}

// WithLevel set default level for the logger
func WithLevel(level Level) Option {
	return func(args *Options) {
		args.Level = level
	}
}

// WithOutput set default output writer for the logger
func WithOutput(out io.Writer) Option {
	return func(args *Options) {
		args.Out = out
	}
}

// WithCallerSkipCount set frame count to skip
func WithCallerSkipCount(c int) Option {
	return func(args *Options) {
		args.CallerSkipCount = c
	}
}

// WithName set name for logger
func WithName(name string) Option {
	return func(args *Options) {
		args.Name = name
	}
}

func SetOption(k, v interface{}) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

// ReadOptions for querying the logs
type ReadOptions struct {
	// Since what time in past to return the logs
	Since time.Time
	// Count specifies number of logs to return
	Count int
	// Stream requests continuous log stream
	Stream bool
}

// ReadOption used for reading the logs
type ReadOption func(*ReadOptions)

// Since sets the time since which to return the log records
func Since(s time.Time) ReadOption {
	return func(o *ReadOptions) {
		o.Since = s
	}
}

// Count sets the number of log records to return
func Count(c int) ReadOption {
	return func(o *ReadOptions) {
		o.Count = c
	}
}
