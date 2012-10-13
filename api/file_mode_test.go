package api

import (
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

	test := func(m uint16, exp FileMode) {
		mode, ok := assertFileMode(m)
		util.Assert(t, ok, "could not convert to file mode")
		util.Assert(t, mode == exp)
	}

	test(m1, ModeNew)
	test(m2, ModeTree)
	test(m3, ModeBlob)
	test(m4, ModeBlobExec)
	test(m5, ModeLink)
	test(m6, ModeCommit)
}
