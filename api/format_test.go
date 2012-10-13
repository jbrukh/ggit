package api

import (
	"github.com/jbrukh/ggit/util"
	"testing"
)

func Test_NewStrFormat(t *testing.T) {
	f := NewStrFormat()
	f.Printf("hello %d", 10)
	f.Lf()
	util.Assert(t, f.String() == "hello 10\n")
}
