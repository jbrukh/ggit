package api

import (
	"testing"
)

func Test_NewStrFormat(t *testing.T) {
	f := NewStrFormat()
	f.Printf("hello %d", 10)
	f.Lf()
	assert(t, f.String() == "hello 10\n")
}
