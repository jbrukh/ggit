//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import (
	"errors"
	"github.com/jbrukh/ggit/util"
	"testing"
)

func Test_ReadBytes(t *testing.T) {
	t1 := parserForString("poop\000")           // simple case
	t2 := parserForString("b")                  // empty token
	t3 := parserForString("    life\000oh\000") // more delims

	util.AssertPanicFree(t, func() {
		util.Assert(t, string(t1.ReadBytes(NUL)) == "poop")
		util.Assert(t, string(t2.ReadBytes('b')) == "")
		util.Assert(t, string(t3.ReadBytes(NUL)) == "    life")
	})
}

func Test_ReadBytesPanic(t *testing.T) {
	t1 := parserForString("")
	t2 := parserForString("hello\000wrong\000token")
	util.AssertPanic(t, func() {
		t1.ReadBytes(NUL)
	})
	util.AssertPanic(t, func() {
		t2.ReadBytes('a') // should not find 'a'
	})
}

func Test_String(t *testing.T) {
	const MSG = "The quick brown fox jumped over the lazy dog."
	t1 := parserForString(MSG)
	t2 := parserForString(MSG)
	t3 := parserForString("")

	util.AssertPanicFree(t, func() {
		util.Assert(t, t1.String() == MSG)
		t2.buf.ReadByte()
		util.Assert(t, t2.String() == MSG[1:])
		util.Assert(t, t3.String() == "")
	})
}

func Test_ConsumePeekString(t *testing.T) {
	const MSG = "The quick brown fox jumped over the lazy dog."
	t1 := parserForString(MSG)

	util.AssertPanicFree(t, func() {
		util.Assert(t, t1.PeekString(3) == "The")
		util.Assert(t, t1.PeekString(9) == "The quick")
		util.Assert(t, t1.PeekString(len(MSG)) == MSG)
		t1.ConsumeString("The ")
		t1.ConsumeString("quick ")
		t1.ConsumeString("brown ")
		t1.ConsumeString("fox ")
		t1.ConsumeString("jumped ")
		t1.ConsumeString("over ")
		t1.ConsumeString("the ")
		t1.ConsumeString("lazy dog.")
		t1.ConsumeString("")
	})

	util.AssertPanic(t, func() {
		t1.ConsumeString("garbage")
	})
}

func Test_dataParse(t *testing.T) {
	err := safeParse(func() {
		// we only care about parseErr's
		panic(errors.New("not a parseErr"))
	})
	if err != nil {
		t.Error("threw an error when not supposed to")
	}

	err = safeParse(func() {
		panicErr("this is a parse error")
	})
	if err == nil {
		t.Error("didn't throw an error when supposed to")
	}
}

func Test_ParseAtoi(t *testing.T) {
	t1 := parserForString("-100\000")
	t2 := parserForString("101\000")
	t3 := parserForString("0\000")
	t4 := parserForString("dog\000")
	t5 := parserForString("eleven\000")
	t6 := parserForString("\000")
	t7 := parserForString("14.3\000")

	util.AssertPanicFree(t, func() {
		util.Assert(t, t1.ParseAtoi(NUL) == -100)
		util.Assert(t, t2.ParseAtoi(NUL) == 101)
		util.Assert(t, t3.ParseAtoi(NUL) == 0)
	})

	util.AssertPanic(t, func() {
		t4.ParseAtoi(NUL)
	})
	util.AssertPanic(t, func() {
		t5.ParseAtoi(NUL)
	})
	util.AssertPanic(t, func() {
		t6.ParseAtoi(NUL)
	})
	util.AssertPanic(t, func() {
		t7.ParseAtoi(NUL)
	})
}

func Test_ParseIntN(t *testing.T) {
	t1 := parserForString("-100\000")
	t2 := parserForString("101\000")
	t3 := parserForString("0\000")
	t4 := parserForString("+11")

	util.AssertPanicFree(t, func() {
		util.Assert(t, t1.ParseIntN(4, 10, 0) == -100)
		util.Assert(t, t2.ParseIntN(3, 10, 0) == 101)
		util.Assert(t, t3.ParseIntN(1, 10, 0) == 0)
		util.Assert(t, t4.ParseIntN(3, 10, 0) == 11)
	})
}

var animals []string = []string{
	"dog",
	"cat",
	"doggie",
}

func Test_ConsumeStrings(t *testing.T) {
	t1 := parserForString("dogcat")
	t2 := parserForString("doggie")

	util.AssertPanicFree(t, func() {
		util.Assert(t, t1.ConsumeStrings(animals) == "dog")
		util.Assert(t, t1.ConsumeStrings(animals) == "cat")
		util.Assert(t, t2.ConsumeStrings(animals) == "dog") // only first match is returned
	})

	t3 := parserForString("dogcat")
	t4 := parserForString("")

	util.AssertPanic(t, func() {
		t3.ConsumeStrings([]string{})
	})
	util.AssertPanic(t, func() {
		t3.ConsumeStrings(nil)
	})
	util.AssertPanic(t, func() {
		t3.ConsumeStrings(objectTypes)
	})
	util.AssertPanic(t, func() {
		t4.ConsumeStrings(animals)
	})
	util.AssertPanic(t, func() {
		t4.ConsumeStrings([]string{""})
	})
}

func Test_Count(t *testing.T) {
	t1 := parserForString("tree 4\000lalala")
	util.AssertPanicFree(t, func() {
		t1.ReadString(NUL)
	})
	util.Assert(t, t1.Count() == 7)
	t1.ResetCount()
	util.Assert(t, t1.Count() == 0)
	util.AssertPanicFree(t, func() {
		t1.String()
	})
	util.Assert(t, t1.Count() == 6)
}
