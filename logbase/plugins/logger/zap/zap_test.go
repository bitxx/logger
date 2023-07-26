package zap

import (
	"github.com/jason-wj/logger/logbase"
	"testing"
)

func TestName(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}

	if l.String() != "zap" {
		t.Errorf("name is error %s", l.String())
	}

	t.Logf("test logger name: %s", l.String())
}

func TestLogf(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}

	logbase.DefaultLogger = l
	logbase.Logf(logbase.InfoLevel, "test logf: %s", "name")
}

func TestSetLevel(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}
	logbase.DefaultLogger = l

	logbase.Init(logbase.WithLevel(logbase.DebugLevel))
	l.Logf(logbase.DebugLevel, "test show debug: %s", "debug msg")

	logbase.Init(logbase.WithLevel(logbase.InfoLevel))
	l.Logf(logbase.DebugLevel, "test non-show debug: %s", "debug msg")
}

func TestWithReportCaller(t *testing.T) {
	var err error
	logbase.DefaultLogger, err = NewLogger(WithCallerSkip(0))
	if err != nil {
		t.Fatal(err)
	}

	logbase.Logf(logbase.InfoLevel, "testing: %s", "WithReportCaller")
}

func TestFields(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}
	logbase.DefaultLogger = l.Fields(map[string]interface{}{
		"x-request-id": "123456abc",
	})
	logbase.DefaultLogger.Log(logbase.InfoLevel, "hello")
}

func TestFile(t *testing.T) {
	/*output, err := writer.NewFileWriter("testdata", "log")
	if err != nil {
		t.Errorf("logger setup error: %s", err.Error())
	}
	//var err error
	core.DefaultLogger, err = NewLogger(core.WithLevel(core.TraceLevel), WithOutput(output))
	if err != nil {
		t.Errorf("logger setup error: %s", err.Error())
	}
	core.DefaultLogger = core.DefaultLogger.Fields(map[string]interface{}{
		"x-request-id": "123456abc",
	})
	fmt.Println(core.DefaultLogger)
	core.DefaultLogger.Log(core.InfoLevel, "hello")*/
}
