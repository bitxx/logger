package logrus

import (
	"errors"
	"github.com/jason-wj/logger/core"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestName(t *testing.T) {
	l := NewLogger()

	if l.String() != "logrus" {
		t.Errorf("error: name expected 'logrus' actual: %s", l.String())
	}

	t.Logf("testing logger name: %s", l.String())
}

func TestWithFields(t *testing.T) {
	l := NewLogger(core.WithOutput(os.Stdout)).Fields(map[string]interface{}{
		"k1": "v1",
		"k2": 123456,
	})

	core.DefaultLogger = l

	core.Log(core.InfoLevel, "testing: Info")
	core.Logf(core.InfoLevel, "testing: %s", "Infof")
}

func TestWithError(t *testing.T) {
	l := NewLogger().Fields(map[string]interface{}{"error": errors.New("boom!")})
	core.DefaultLogger = l

	core.Log(core.InfoLevel, "testing: error")
}

func TestWithLogger(t *testing.T) {
	// with *logrus.Logger
	l := NewLogger(WithLogger(logrus.StandardLogger())).Fields(map[string]interface{}{
		"k1": "v1",
		"k2": 123456,
	})
	core.DefaultLogger = l
	core.Log(core.InfoLevel, "testing: with *logrus.Logger")

	// with *logrus.Entry
	el := NewLogger(WithLogger(logrus.NewEntry(logrus.StandardLogger()))).Fields(map[string]interface{}{
		"k3": 3.456,
		"k4": true,
	})
	core.DefaultLogger = el
	core.Log(core.InfoLevel, "testing: with *logrus.Entry")
}

func TestJSON(t *testing.T) {
	core.DefaultLogger = NewLogger(WithJSONFormatter(&logrus.JSONFormatter{}))

	core.Logf(core.InfoLevel, "test logf: %s", "name")
}

func TestSetLevel(t *testing.T) {
	core.DefaultLogger = NewLogger()

	core.Init(core.WithLevel(core.DebugLevel))
	core.Logf(core.DebugLevel, "test show debug: %s", "debug msg")

	core.Init(core.WithLevel(core.InfoLevel))
	core.Logf(core.DebugLevel, "test non-show debug: %s", "debug msg")
}

func TestWithReportCaller(t *testing.T) {
	core.DefaultLogger = NewLogger(ReportCaller())

	core.Logf(core.InfoLevel, "testing: %s", "WithReportCaller")
}
