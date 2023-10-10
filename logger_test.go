package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	//type: logrus zap default
	log := NewLogger(
		WithType("logrus"),
		WithPath("temp/logs"),
		WithLevel("info"),
		WithStdout("default"),
		WithCap(10),
	)

	log.Info("xxxxxx")
	log.Warn("######")
}
