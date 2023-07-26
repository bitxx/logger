package zap

import (
	"github.com/jason-wj/logger/core"
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

	core.DefaultLogger = l
	core.Logf(core.InfoLevel, "test logf: %s", "name")
}

func TestSetLevel(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}
	core.DefaultLogger = l

	core.Init(core.WithLevel(core.DebugLevel))
	l.Logf(core.DebugLevel, "test show debug: %s", "debug msg")

	core.Init(core.WithLevel(core.InfoLevel))
	l.Logf(core.DebugLevel, "test non-show debug: %s", "debug msg")
}

func TestWithReportCaller(t *testing.T) {
	var err error
	core.DefaultLogger, err = NewLogger(WithCallerSkip(0))
	if err != nil {
		t.Fatal(err)
	}

	core.Logf(core.InfoLevel, "testing: %s", "WithReportCaller")
}

func TestFields(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}
	core.DefaultLogger = l.Fields(map[string]interface{}{
		"x-request-id": "123456abc",
	})
	core.DefaultLogger.Log(core.InfoLevel, "hello")
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
