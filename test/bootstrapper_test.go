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

func assertNoErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("an error occurred: %s", err.Error())
	}
}

func assertf(t *testing.T, b bool, format string, items ...interface{}) {
	if !b {
		t.Errorf(format, items)
	}
}

func Test_CreateTestRepo(t *testing.T) {
	repo, err := CreateTestRepo("cases/single_commit.sh")
	assertNoErr(t, err)
	err = repo.Destroy()
	assertNoErr(t, err)
}
