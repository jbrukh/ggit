//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package util

import (
	"errors"
	"testing"
)

func Test_ReadBytes(t *testing.T) {
	t1 := ParserForString("poop\000")           // simple case
	t2 := ParserForString("b")                  // empty token
	t3 := ParserForString("    life\000oh\000") // more delims

	AssertPanicFree(t, func() {
		Assert(t, string(t1.ReadBytes('\000')) == "poop")
		Assert(t, string(t2.ReadBytes('b')) == "")
		Assert(t, string(t3.ReadBytes('\000')) == "    life")
	})
}

func Test_ReadBytesPanic(t *testing.T) {
	t1 := ParserForString("")
	t2 := ParserForString("hello\000wrong\000token")
	AssertPanic(t, func() {
		t1.ReadBytes('\000')
	})
	AssertPanic(t, func() {
		t2.ReadBytes('a') // should not find 'a'
	})
}

func Test_String(t *testing.T) {
	const MSG = "The quick brown fox jumped over the lazy dog."
	t1 := ParserForString(MSG)
	t2 := ParserForString(MSG)
	t3 := ParserForString("")

	AssertPanicFree(t, func() {
		Assert(t, t1.String() == MSG)
		t2.buf.ReadByte()
		Assert(t, t2.String() == MSG[1:])
		Assert(t, t3.String() == "")
	})
}

func Test_ConsumePeekString(t *testing.T) {
	const MSG = "The quick brown fox jumped over the lazy dog."
	t1 := ParserForString(MSG)

	AssertPanicFree(t, func() {
		Assert(t, t1.PeekString(3) == "The")
		Assert(t, t1.PeekString(9) == "The quick")
		Assert(t, t1.PeekString(len(MSG)) == MSG)
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

	AssertPanic(t, func() {
		t1.ConsumeString("garbage")
	})
}

func Test_dataParse(t *testing.T) {
	err := SafeParse(func() {
		// we only care about parseErr's
		panic(errors.New("not a parseErr"))
	})
	if err != nil {
		t.Error("threw an error when not supposed to")
	}

	err = SafeParse(func() {
		PanicErr("this is a parse error")
	})
	if err == nil {
		t.Error("didn't throw an error when supposed to")
	}
}

func Test_ParseAtoi(t *testing.T) {
	t1 := ParserForString("-100\000")
	t2 := ParserForString("101\000")
	t3 := ParserForString("0\000")
	t4 := ParserForString("dog\000")
	t5 := ParserForString("eleven\000")
	t6 := ParserForString("\000")
	t7 := ParserForString("14.3\000")

	AssertPanicFree(t, func() {
		Assert(t, t1.ParseAtoi('\000') == -100)
		Assert(t, t2.ParseAtoi('\000') == 101)
		Assert(t, t3.ParseAtoi('\000') == 0)
	})

	AssertPanic(t, func() {
		t4.ParseAtoi('\000')
	})
	AssertPanic(t, func() {
		t5.ParseAtoi('\000')
	})
	AssertPanic(t, func() {
		t6.ParseAtoi('\000')
	})
	AssertPanic(t, func() {
		t7.ParseAtoi('\000')
	})
}

func Test_ParseIntN(t *testing.T) {
	t1 := ParserForString("-100\000")
	t2 := ParserForString("101\000")
	t3 := ParserForString("0\000")
	t4 := ParserForString("+11")

	AssertPanicFree(t, func() {
		Assert(t, t1.ParseIntN(4, 10, 0) == -100)
		Assert(t, t2.ParseIntN(3, 10, 0) == 101)
		Assert(t, t3.ParseIntN(1, 10, 0) == 0)
		Assert(t, t4.ParseIntN(3, 10, 0) == 11)
	})
}

var animals []string = []string{
	"dog",
	"cat",
	"doggie",
}

func Test_ConsumeStrings(t *testing.T) {
	t1 := ParserForString("dogcat")
	t2 := ParserForString("doggie")

	AssertPanicFree(t, func() {
		Assert(t, t1.ConsumeStrings(animals) == "dog")
		Assert(t, t1.ConsumeStrings(animals) == "cat")
		Assert(t, t2.ConsumeStrings(animals) == "dog") // only first match is returned
	})

	t3 := ParserForString("dogcat")
	t4 := ParserForString("")

	AssertPanic(t, func() {
		t3.ConsumeStrings([]string{})
	})
	AssertPanic(t, func() {
		t3.ConsumeStrings(nil)
	})
	AssertPanic(t, func() {
		t3.ConsumeStrings([]string{"blob", "tree", "commit", "tag"})
	})
	AssertPanic(t, func() {
		t4.ConsumeStrings(animals)
	})
	AssertPanic(t, func() {
		t4.ConsumeStrings([]string{""})
	})
}

func Test_Count(t *testing.T) {
	t1 := ParserForString("tree 4\000lalala")
	AssertPanicFree(t, func() {
		t1.ReadString('\000')
	})
	Assert(t, t1.Count() == 7)
	t1.ResetCount()
	Assert(t, t1.Count() == 0)
	AssertPanicFree(t, func() {
		t1.String()
	})
	Assert(t, t1.Count() == 6)
}
