//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import (
	"fmt"
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/util"
	"math/rand"
	"testing"
)

func makeTestBlobWithSize(size int, contents string) string {
	return fmt.Sprintf("blob %d\000%s", size, contents)
}

func makeTestBlob(contents string) string {
	return makeTestBlobWithSize(len(contents), contents)
}

// TODO add more cases
var testCasesBlobs = []string{
	"a",
	"hello",
	`Package os provides a platform-independent interface to operating system functionality. The design is Unix-like, although the error handling is Go-like; failing calls return values of type error rather than error numbers. Often, more information is available within the error. For example, if a call that takes a file name fails, such as Open or Stat, the error will include the failing file name when printed and will be of type *PathError, which may be unpacked for more information. The os interface is intended to be uniform across all operating systems. Features not generally available appear in the system-specific package syscall.`,
	"",
}

func Test_parseValidBlob(t *testing.T) {
	for _, v := range testCasesBlobs {
		tb := makeTestBlob(v)
		p := objectParserForString(tb)
		hdr, err := p.ParseHeader()
		util.AssertNoErr(t, err)
		util.Assert(t, hdr.Type() == objects.ObjectBlob)
		util.Assert(t, hdr.Size() == int64(len(v)))

		o, err := p.ParsePayload()
		util.AssertNoErr(t, err)
		util.Assert(t, o.Header().Type() == objects.ObjectBlob)
		util.Assert(t, o.Header().Size() == int64(len(v)))
		util.Assert(t, o.ObjectId() == nil) // wasn't set in the test scenario

		var b *Blob
		util.AssertPanicFree(t, func() {
			b = o.(*Blob)
		})
		util.AssertEqualString(t, string(b.Data()), v)
	}
}

func Test_parseInvalidBlobHeader(t *testing.T) {
	test := func(badBlob string) {
		p := objectParserForString(badBlob)
		hdr, e := p.ParseHeader()
		util.Assert(t, e != nil)
		util.Assert(t, hdr == nil)
	}
	test("bad 100\000dsfsdfsdfsdf")
	test("bl0b 10\000est")
	test("100\000hehe")
	test("bad \000hoho")
	test("no header at all")
}

func Test_parseBlobBadSize(t *testing.T) {
	test := func(contents string) {
		l := len(contents)
		size := rand.Intn(l+1) - 1 + 2*rand.Intn(2)*l // never equal to l
		tb := makeTestBlobWithSize(size, contents)
		p := objectParserForString(tb)
		hdr, e := p.ParseHeader()
		util.AssertNoErr(t, e)
		util.Assert(t, hdr.Type() == objects.ObjectBlob)
		util.Assert(t, hdr.Size() == int64(size))

		// should not be able to parse
		_, e = p.ParsePayload()
		util.Assert(t, e != nil)
	}
	for _, v := range testCasesBlobs {
		test(v)
	}
}
