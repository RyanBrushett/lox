package glox

import (
	"errors"
	"testing"
)

func TestRuntimeError(t *testing.T) {
	expected := "[Line 2] Error0: this is a test string\n"
	actual := RuntimeError(2, errors.New("this is a test string")).Error()

	if actual != expected {
		t.Errorf("Expected: %s Actual: %s", expected, actual)
	}
}

func TestRuntimeRun(t *testing.T) {
	runtime := NewRuntime()
	runtime.Run(`"unterminated string`)
	if !runtime.HadError {
		t.Errorf("Unterminated strings should be a runtime error")
	}
}
