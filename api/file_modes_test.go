//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import (
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/util"
	"testing"
)

func Test_assertFileMode(t *testing.T) {
	m1 := uint16(0000000)
	m2 := uint16(0040000)
	m3 := uint16(0100644)
	m4 := uint16(0100755)
	m5 := uint16(0120000)
	m6 := uint16(0160000)

	test := func(m uint16, exp objects.FileMode) {
		mode, ok := assertFileMode(m)
		util.Assert(t, ok, "could not convert to file mode")
		util.Assert(t, mode == exp)
	}

	test(m1, objects.ModeNew)
	test(m2, objects.ModeTree)
	test(m3, objects.ModeBlob)
	test(m4, objects.ModeBlobExec)
	test(m5, objects.ModeLink)
	test(m6, objects.ModeCommit)
}

func Test_parseValidFileMode(t *testing.T) {
	p := objectParserForString("0000000\n0040000\n0100644\n0100755\n0120000\n0160000\n")
	util.Assert(t, p.ParseFileMode(LF) == objects.ModeNew)
	util.Assert(t, p.ParseFileMode(LF) == objects.ModeTree)
	util.Assert(t, p.ParseFileMode(LF) == objects.ModeBlob)
	util.Assert(t, p.ParseFileMode(LF) == objects.ModeBlobExec)
	util.Assert(t, p.ParseFileMode(LF) == objects.ModeLink)
	util.Assert(t, p.ParseFileMode(LF) == objects.ModeCommit)
	util.Assert(t, p.EOF())
}

func Test_parseInvalidFileMode(t *testing.T) {
	// test non-file modes
	p := objectParserForString("000200\n002000\n000644\n000755\n0120200\n01600990\n")
	for !p.EOF() {
		util.AssertPanic(t, func() {
			m := p.ParseFileMode(LF)
			m++ // for compilation, should not get here
		})
	}
}
