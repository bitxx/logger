package zap

import (
	"github.com/bitxx/logger/logbase"
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Options struct {
	logbase.Options
}

type callerSkipKey struct{}

func WithCallerSkip(i int) logbase.Option {
	return logbase.SetOption(callerSkipKey{}, i)
}

type configKey struct{}

// WithConfig pass zap.Config to logger
func WithConfig(c zap.Config) logbase.Option {
	return logbase.SetOption(configKey{}, c)
}

type encoderConfigKey struct{}

// WithEncoderConfig pass zapcore.EncoderConfig to logger
func WithEncoderConfig(c zapcore.EncoderConfig) logbase.Option {
	return logbase.SetOption(encoderConfigKey{}, c)
}

type namespaceKey struct{}

func WithNamespace(namespace string) logbase.Option {
	return logbase.SetOption(namespaceKey{}, namespace)
}

type writerKey struct{}

func WithOutput(out io.Writer) logbase.Option {
	return logbase.SetOption(writerKey{}, out)
}
