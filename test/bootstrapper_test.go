package test

import (
	"runtime"
	"testing"
)

func assert(t *testing.T, b bool, items ...interface{}) {
	if !b {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "(unknown file)"
		}
		t.Errorf("%s:%d: %s", file, line, items)
	}
}

func assertf(t *testing.T, b bool, format string, items ...interface{}) {
	if !b {
		t.Errorf(format, items)
	}
}

func Test_CreateTestRepo(t *testing.T) {
	repo, err := CreateTestRepo("cases/empty_repo.sh")
	assert(t, err == nil)
	err = repo.Destroy()
	assert(t, err == nil)
}
