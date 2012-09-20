package ggit

import (
    "bufio"
    "bytes"
    "errors"
    "fmt"
    "testing"
)

func parserForBytes(b []byte) *dataParser {
    return &dataParser{
        bufio.NewReader(bytes.NewBuffer(b)),
    }
}

func parserForString(str string) *dataParser {
    return parserForBytes([]byte(str))
}

func Test_ReadBytes(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            t.Error("parser paniced with error: " + fmt.Sprintf("%v", r))
        }
    }()

    t1 := parserForString("poop\000")           // simple case
    t2 := parserForString("b")                  // empty token
    t3 := parserForString("    life\000oh\000") // more delims

    assert(t, string(t1.ReadBytes(NUL)) == "poop")
    assert(t, string(t2.ReadBytes('b')) == "")
    assert(t, string(t3.ReadBytes(NUL)) == "    life")
}

func Test_ReadBytesPanic(t *testing.T) {
    t1 := parserForString("")
    t2 := parserForString("hello\000wrong\000token")
    assertPanic(t, func() {
        t1.ReadBytes(NUL)
    })
    assertPanic(t, func() {
        t2.ReadBytes('a') // should not find 'a'
    })
}

func Test_String(t *testing.T) {
    const MSG = "The quick brown fox jumped over the lazy dog."
    t1 := parserForString(MSG)
    t2 := parserForString(MSG)
    t3 := parserForString("")
    assert(t, t1.String() == MSG)

    t2.buf.ReadByte()
    assert(t, t2.String() == MSG[1:])
    assert(t, t3.String() == "")
}

func Test_ConsumePeekString(t *testing.T) {
    const MSG = "The quick brown fox jumped over the lazy dog."
    t1 := parserForString(MSG)
    assert(t, t1.PeekString(3) == "The")
    assert(t, t1.PeekString(9) == "The quick")
    assert(t, t1.PeekString(len(MSG)) == MSG)

    assertPanicFree(t, func() {
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

    assertPanic(t, func() {
        t1.ConsumeString("garbage")
    })
}

func Test_dataParse(t *testing.T) {
    err := dataParse(func() {
        // we only care about parseErr's
        panic(errors.New("not a parseErr"))
    })
    if err != nil {
        t.Error("threw an error when not supposed to")
    }

    err = dataParse(func() {
        panic(parseErr("this is a parse error"))
    })
    if err == nil {
        t.Error("didn't throw an error when supposed to")
    }
}

func Test_ParseObjectId(t *testing.T) {
    var oid *ObjectId
    t1 := parserForString(testOidCrazy)
    oid = t1.ParseObjectId()
    assert(t, oid.String() == testOidCrazy)
}

func Test_ParseAtoi(t *testing.T) {
    t1 := parserForString("-100\000")
    t2 := parserForString("101\000")
    t3 := parserForString("0\000")
    t4 := parserForString("dog\000")
    t5 := parserForString("eleven\000")
    t6 := parserForString("\000")
    t7 := parserForString("14.3\000")

    assertPanicFree(t, func() {
        assert(t, t1.ParseAtoi(NUL) == -100)
        assert(t, t2.ParseAtoi(NUL) == 101)
        assert(t, t3.ParseAtoi(NUL) == 0)
    })

    assertPanic(t, func() {
        t4.ParseAtoi(NUL)
    })
    assertPanic(t, func() {
        t5.ParseAtoi(NUL)
    })
    assertPanic(t, func() {
        t6.ParseAtoi(NUL)
    })
    assertPanic(t, func() {
        t7.ParseAtoi(NUL)
    })
}
