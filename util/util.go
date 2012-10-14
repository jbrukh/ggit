package util

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

const DefaultGitDir = ".git"

func InferGitDir(pth string) string {
	_, file := filepath.Split(pth)
	if file != DefaultGitDir {
		return path.Join(pth, DefaultGitDir)
	}
	return pth
}

// IsValidRepo validates a repository path to make sure it has
// the right format and that it exists.	
func IsValidRepo(pth string) bool {
	p := InferGitDir(pth)
	if _, e := os.Stat(p); e != nil {
		return false
	}
	// TODO: may want to do other checks here...
	return true
}

func UniqueHex16() string {
	return UniqueHex20()[:16]
}

func UniqueHex20() string {
	buf := make([]byte, 20)
	io.ReadFull(rand.Reader, buf)
	return fmt.Sprintf("%x", buf)
}

// ================================================================= //
// ASSERT STATEMENTS
// ================================================================= //

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

func AssertEqualString(t *testing.T, one, two string) {
	if one != two {
		t.Errorf("expecting: '%s' but got '%s'\n", one, two)
	}
}

func AssertEqualInt(t *testing.T, one, two int) {
	if one != two {
		t.Errorf("expecting: '%d' but got '%d'\n", one, two)
	}
}

func AssertNoErr(t *testing.T, err error) {
	if err != nil {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "(unknown file)"
		}
		t.Errorf("%s:%d: %s", file, line, err)
	}
}

func AssertNoErrOrDie(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("a fatal error occurred: %s", err)
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
			t.Error("failed because it panicked")
		}
	}()
	f()
}
