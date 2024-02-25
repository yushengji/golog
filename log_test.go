package fslog

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestLevel(t *testing.T) {
	defer Flush()
	SetLv(DebugLv)
	buffer := &bytes.Buffer{}
	NewOutput(buffer)
	if Lv() != DebugLv {
		t.Error("level not match")
	}

	Info("info message")
	out := buffer.String()
	if !strings.Contains(out, "info message") {
		t.Error("lower level cannot log higher level message")
	}

	SetLv(ErrorLv)
	defer SetLv(DebugLv)
	if Lv() != ErrorLv {
		t.Error("level not match")
	}

	Debug("debug message")
	out = buffer.String()
	if strings.Contains(out, "debug message") {
		t.Error("higher level can log lower level message")
	}
}

func TestDebug(t *testing.T) {
	defer Flush()
	SetLv(DebugLv)
	buffer := &bytes.Buffer{}
	NewOutput(buffer)
	Debug("test debug %s", "fmt")
	out := buffer.String()
	if !strings.Contains(out, `"level":"debug"`) {
		t.Error("debug level not match")
	}

	if !strings.Contains(out, "test debug") {
		t.Error("debug message not match")
	}

	if !strings.Contains(out, "fmt") {
		t.Error("debug fmt message not match")
	}

	if !strings.Contains(out, "test debug fmt") {
		t.Error("debug all message not match")
	}

	if !strings.Contains(out, "fslog/log_test.go") {
		t.Error("debug caller not match")
	}
}

func TestInfo(t *testing.T) {
	defer Flush()
	SetLv(DebugLv)
	buffer := &bytes.Buffer{}
	NewOutput(buffer)
	Info("test info %s", "fmt")
	out := buffer.String()
	if !strings.Contains(out, `"level":"info"`) {
		t.Error("info level not match")
	}

	if !strings.Contains(out, "test info") {
		t.Error("info message not match")
	}

	if !strings.Contains(out, "fmt") {
		t.Error("info fmt message not match")
	}

	if !strings.Contains(out, "test info fmt") {
		t.Error("info fmt all message not match")
	}

	if !strings.Contains(out, "fslog/log_test.go") {
		t.Error("info caller not match")
	}
}

func TestWarn(t *testing.T) {
	defer Flush()
	SetLv(DebugLv)
	buffer := &bytes.Buffer{}
	NewOutput(buffer)
	Warn("test warn %s", "fmt")
	out := buffer.String()
	if !strings.Contains(out, `"level":"warn"`) {
		t.Error("warn level not match")
	}

	if !strings.Contains(out, "test warn") {
		t.Error("warn message not match")
	}

	if !strings.Contains(out, "fmt") {
		t.Error("warn fmt message not match")
	}

	if !strings.Contains(out, "test warn fmt") {
		t.Error("warn fmt all message not match")
	}

	if !strings.Contains(out, "fslog/log_test.go") {
		t.Error("warn caller not match")
	}
}

func TestError(t *testing.T) {
	defer Flush()
	SetLv(DebugLv)
	buffer := &bytes.Buffer{}
	NewOutput(buffer)
	Error("test error %s", "fmt")
	out := buffer.String()
	if !strings.Contains(out, `"level":"error"`) {
		t.Error("error level not match")
	}

	if !strings.Contains(out, "test error") {
		t.Error("error message not match")
	}

	if !strings.Contains(out, "fmt") {
		t.Error("error fmt message not match")
	}

	if !strings.Contains(out, "test error fmt") {
		t.Error("error fmt all message not match")
	}

	if !strings.Contains(out, "fslog/log_test.go") {
		t.Error("error caller not match")
	}
}

func TestWith(t *testing.T) {
	defer Flush()
	buffer := &bytes.Buffer{}
	NewOutput(buffer)
	With("key", "value").Info("test with")
	out := strings.Trim(buffer.String(), " ")
	if !strings.Contains(out, `"key":"value"`) {
		t.Error("with key value not match")
	}
}

func TestNewFileOutput(t *testing.T) {
	defer Flush()
	SetLv(DebugLv)
	NewFileOutput(WithFilename("logs/test"))
	Debug("test new file output")
	contentBytes, err := os.ReadFile("logs/test")
	if err != nil {
		t.Error("read file error", err)
	}

	content := string(contentBytes)
	if !strings.Contains(content, "test new file output") {
		t.Error("file output content not match")
	}
}

func TestFileExtend(t *testing.T) {
	for i := 0; i < 100000; i++ {
		Info(strconv.Itoa(i))
	}
}

func TestPrintLocation(t *testing.T) {
	defer Flush()
	buffer := &bytes.Buffer{}
	NewOutput(buffer)
	Info("test location")
	if !strings.Contains(buffer.String(), `"caller":"fslog/log_test.go:`) {
		t.Error("caller output have deviation")
	}
}
