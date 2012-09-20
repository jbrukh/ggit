package ggit

import (
    "bufio"
)

// ================================================================= //
// OBJECT HEADER PARSING
// ================================================================= //

// ObjectHeader is the deserialized (and more efficiently stored)
// version of a git object header
type objectHeader struct {
    Type ObjectType
    Size int
}

func parseObjectHeader(buf *bufio.Reader) (*objectHeader, error) {
    h := new(ObjectHeader)
    p := dataParser{buf}
    err := dataParse(func() {
        h.Type = ObjectType(p.ConsumeStrings(objectTypes))
        p.ConsumeByte(SP)
        h.Size = p.ParseInt(NUL)
    })
    return h, err
}
