package util

import (
	"runtime"
	"testing"
)

func Assert(t *testing.T, b bool, items ...interface{}) {
	if !b {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "(unknown file)"
		}
		t.Errorf("%s:%d: %s", file, line, items)
	}
}

func Assertf(t *testing.T, b bool, format string, items ...interface{}) {
	if !b {
		t.Errorf(format, items)
	}
}

func AssertNoErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("an error occurred: %s", err)
	}
}

func AssertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
		// should never get here
	}()
	f()
	// TODO: use runtime to get the line numbers of the caller
	t.Error("was expecting a panic")
}

func AssertPanicFree(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r != nil {
			// TODO: use runtime to get the line numbers of the caller
			t.Error("failed because it panicked")
		}
	}()
	f()
}
