package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	//type: zap default
	log := NewLogger(
		WithType("zap"),
		WithPath("temp/logs"),
		WithLevel("info"),
		WithStdout("default"),
		WithCap(10),
	)

	log.Info("xxxxxx")
	log.Warn("######")
}
