package logrus

import (
	"errors"
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
	l := NewLogger(logbase.WithOutput(os.Stdout)).Fields(map[string]interface{}{
		"k1": "v1",
		"k2": 123456,
	})

	logbase.DefaultLogger = l

	logbase.Log(logbase.InfoLevel, "testing: Info")
	logbase.Logf(logbase.InfoLevel, "testing: %s", "Infof")
}

func TestWithError(t *testing.T) {
	l := NewLogger().Fields(map[string]interface{}{"error": errors.New("boom!")})
	logbase.DefaultLogger = l

	logbase.Log(logbase.InfoLevel, "testing: error")
}

func TestWithLogger(t *testing.T) {
	// with *logrus.Logger
	l := NewLogger(WithLogger(logrus.StandardLogger())).Fields(map[string]interface{}{
		"k1": "v1",
		"k2": 123456,
	})
	logbase.DefaultLogger = l
	logbase.Log(logbase.InfoLevel, "testing: with *logrus.Logger")

	// with *logrus.Entry
	el := NewLogger(WithLogger(logrus.NewEntry(logrus.StandardLogger()))).Fields(map[string]interface{}{
		"k3": 3.456,
		"k4": true,
	})
	logbase.DefaultLogger = el
	logbase.Log(logbase.InfoLevel, "testing: with *logrus.Entry")
}

func TestJSON(t *testing.T) {
	logbase.DefaultLogger = NewLogger(WithJSONFormatter(&logrus.JSONFormatter{}))

	logbase.Logf(logbase.InfoLevel, "test logf: %s", "name")
}

func TestSetLevel(t *testing.T) {
	logbase.DefaultLogger = NewLogger()

	logbase.Init(logbase.WithLevel(logbase.DebugLevel))
	logbase.Logf(logbase.DebugLevel, "test show debug: %s", "debug msg")

	logbase.Init(logbase.WithLevel(logbase.InfoLevel))
	logbase.Logf(logbase.DebugLevel, "test non-show debug: %s", "debug msg")
}

func TestWithReportCaller(t *testing.T) {
	logbase.DefaultLogger = NewLogger(ReportCaller())

	logbase.Logf(logbase.InfoLevel, "testing: %s", "WithReportCaller")
}
