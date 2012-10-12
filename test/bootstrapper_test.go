package test

import (
	"runtime"
	"testing"
)

// TODO: these replicate the functionality in util.go, but
// there is no good way to make them common and private.

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

func Test_SanityTest(t *testing.T) {
	repo, err := Repo("cases/single_commit.sh")
	assertNoErr(t, err)
	err = repo.Destroy()
	assertNoErr(t, err)
}
