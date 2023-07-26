package zap

import (
	"io"
	"logger/core"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Options struct {
	core.Options
}

type callerSkipKey struct{}

func WithCallerSkip(i int) core.Option {
	return core.SetOption(callerSkipKey{}, i)
}

type configKey struct{}

// WithConfig pass zap.Config to logger
func WithConfig(c zap.Config) core.Option {
	return core.SetOption(configKey{}, c)
}

type encoderConfigKey struct{}

// WithEncoderConfig pass zapcore.EncoderConfig to logger
func WithEncoderConfig(c zapcore.EncoderConfig) core.Option {
	return core.SetOption(encoderConfigKey{}, c)
}

type namespaceKey struct{}

func WithNamespace(namespace string) core.Option {
	return core.SetOption(namespaceKey{}, namespace)
}

type writerKey struct{}

func WithOutput(out io.Writer) core.Option {
	return core.SetOption(writerKey{}, out)
}
