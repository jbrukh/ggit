//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
assert.go provides convenience methods for generic testing.
*/
package util

import (
	"runtime"
	"testing"
)

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
		t.Errorf(format, items...)
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
			t.Fatal("failed because it panicked")
		}
	}()
	f()
}
