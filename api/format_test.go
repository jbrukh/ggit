//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
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
