package ggit

import (
    "bufio"
    "bytes"
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

func Test_TokenBytes(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            t.Error("parser paniced with error: " + fmt.Sprintf("%v", r))
        }
    }()

    t1 := parserForString("poop\000")           // simple case
    t2 := parserForString("b")                  // empty token
    t3 := parserForString("    life\000oh\000") // more delims

    assert(t, string(t1.TokenBytes(NUL)) == "poop")
    assert(t, string(t2.TokenBytes('b')) == "")
    assert(t, string(t3.TokenBytes(NUL)) == "    life")
}

func Test_TokenBytesPanic(t *testing.T) {
    t1 := parserForString("")
    t2 := parserForString("hello\000wrong\000token")
    assertPanic(t, func() {
        t1.TokenBytes(NUL)
    })
    assertPanic(t, func() {
        t2.TokenBytes('a') // should not find 'a'
    })
}

func Test_TokenStringInt(t *testing.T) {
    t1 := parserForString("100E")
    t2 := parserForString("-100E")
    t3 := parserForString("0E")
    t4 := parserForString("1000000E")

    assert(t, t1.TokenStringInt('E') == 100)
    assert(t, t2.TokenStringInt('E') == -100)
    assert(t, t3.TokenStringInt('E') == 0)
    assert(t, t4.TokenStringInt('E') == 1000000)

}

func Test_TokenStringIntPanic(t *testing.T) {
    t1 := parserForString(".100E")
    t2 := parserForString("catE")
    t3 := parserForString("100")
    t4 := parserForString("")

    assertPanic(t, func() {
        t1.TokenStringInt('E')
    })
    assertPanic(t, func() {
        t2.TokenStringInt('E') // should not find 'a'
    })
    assertPanic(t, func() {
        t3.TokenStringInt('E')
    })
    assertPanic(t, func() {
        t4.TokenStringInt('E') // should not find 'a'
    })
}

func Test_FlushString(t *testing.T) {
    const MSG = "The quick brown fox jumped over the lazy dog."
    t1 := parserForString(MSG)
    t2 := parserForString(MSG)
    t3 := parserForString("")
    assert(t, t1.FlushString() == MSG)

    t2.buf.ReadByte()
    assert(t, t2.FlushString() == MSG[1:])
    assert(t, t3.FlushString() == "")
}

func Test_FlushString(t *testing.T) {
    const MSG = "The quick brown fox jumped over the lazy dog."
    t1 := parserForString(MSG)
    assert(t, t1.PeekString("The"))
    assert(t, t1.PeekString("The quick"))
    assert(t, t1.PeekString(MSG))

    assert(t, t1.VerifyString("The "))
    assert(t, t1.VerifyString("quick "))
    assert(t, t1.VerifyString("brown "))
    assert(t, t1.VerifyString("fox "))
    assert(t, t1.VerifyString("jumped "))
    assert(t, t1.VerifyString("over "))
    assert(t, t1.VerifyString("the "))
    assert(t, t1.VerifyString("lazy dog."))
    assert(t, t1.VerifyString(""))

    assertPanic(t, func() {
        t1.VerifyString("garbage")
    })
}
