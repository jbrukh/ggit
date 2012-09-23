package ggit

import (
    "bufio"
    "bytes"
    "runtime"
    "testing"
)

func assert(t *testing.T, b bool, items ...interface{}) {
    if !b {
        _, file, line, ok := runtime.Caller(1)
        if !ok {
            file = "(unknown file)"
        }
        t.Errorf("%s:%d: %s", file, line, items)
    }
}

func assertf(t *testing.T, b bool, format string, items ...interface{}) {
    if !b {
        t.Errorf(format, items)
    }
}

func assertPanic(t *testing.T, f func()) {
    defer func() {
        if r := recover(); r != nil {
            return
        }
        // should never get here
    }()
    f()
    // TODO: use runtime to get the line numbers of the caller
    t.Error("was expecting a panic")
}

func assertPanicFree(t *testing.T, f func()) {
    defer func() {
        if r := recover(); r != nil {
            // TODO: use runtime to get the line numbers of the caller
            t.Error("failed because it panicked")
        }
    }()
    f()
}

func objectParserForString(str string) *objectParser {
    p := new(objectParser)
    p.buf = readerForString(str)
    return p
}

func parserForBytes(b []byte) *dataParser {
    return &dataParser{
        buf: bufio.NewReader(bytes.NewBuffer(b)),
    }
}

func parserForString(str string) *dataParser {
    return parserForBytes([]byte(str))
}

func readerForString(str string) *bufio.Reader {
    return bufio.NewReader(bytes.NewBufferString(str))
}
