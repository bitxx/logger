package logger

import (
	"github.com/jason-wj/logger/logbase"
	"testing"
)

func TestLogger(t *testing.T) {
	//type: logrus zap default
	SetupLogger(
		WithType("logrus"),
		WithPath("temp/logs"),
		WithLevel("info"),
		WithStdout("file"),
		WithCap(10),
	)

	log := NewLogger(logbase.DefaultLogger)
	log.Info("xxxxxx")
}
