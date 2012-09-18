package ggit

import (
    "testing"
)

func newBlobFromString(blobStr string) *Blob {
    return &Blob{
        RawObject: RawObject{
            bytes: []byte(blobStr),
        },
        repo: nil,
    }
}

var passBlob1 *Blob = newBlobFromString("blob 0\000")
var passBlob2 *Blob = newBlobFromString("blob 1\0001")
var passBlob3 *Blob = newBlobFromString("blob 10\0001234567890")
var passBlob4 *Blob = newBlobFromString("blob 10\0001 3 5 7 9 ")
var passBlob5 *Blob = newBlobFromString("blob 10\0001 3 \000 7 9 ")

var failBlob1 *Blob = newBlobFromString("")              // no data
var failBlob2 *Blob = newBlobFromString("eskimoes")      // invalid data
var failBlob3 *Blob = newBlobFromString("blob")          // no space
var failBlob4 *Blob = newBlobFromString("blob 10")       // no separator
var failBlob5 *Blob = newBlobFromString("blob0\000")     // bad space
var failBlob6 *Blob = newBlobFromString(" blob 0\000")   // extra white space
var failBlob7 *Blob = newBlobFromString("blob  0\000")   // extra white space
var failBlob8 *Blob = newBlobFromString("blood 0\000")   // wrong type
var failBlob9 *Blob = newBlobFromString("blob zero\000") // wrong size

var failBlob10 *Blob = newBlobFromString("blob 11\000hehe") // wrong size
var failBlob11 *Blob = newBlobFromString("blob\000hehe")    // bad header

func TestType(t *testing.T) {
    b := new(Blob)
    assert(t, b.Type() == OBJECT_BLOB)
}

func TestHeaderPass(t *testing.T) {
    assertHeader := func(b *Blob, otype ObjectType, osize int) {
        h, e := b.Header()
        assert(t, e == nil, "couldn't parse a valid blob: ", e)
        assert(t, otype == h.Type, "type is wrong: ", otype)
        assert(t, osize == h.Size, "size is wrong: ", osize)
    }
    assertHeader(passBlob1, OBJECT_BLOB, 0)
    assertHeader(passBlob2, OBJECT_BLOB, 1)
    assertHeader(passBlob3, OBJECT_BLOB, 10)
    assertHeader(passBlob4, OBJECT_BLOB, 10)
    assertHeader(passBlob5, OBJECT_BLOB, 10)
}

func TestPayloadPass(t *testing.T) {
    assertPayload := func(b *Blob, payload string) {
        p, e := b.Payload()
        assert(t, e == nil, "couldn't parse a payload: ", e)
        assert(t, string(p) == payload, "payload is wrong")
    }
    assertPayload(passBlob1, "")
    assertPayload(passBlob2, "1")
    assertPayload(passBlob3, "1234567890")
    assertPayload(passBlob4, "1 3 5 7 9 ")
    assertPayload(passBlob5, "1 3 \000 7 9 ")
}

func TestHeaderFail(t *testing.T) {
    assertHeader := func(b *Blob) {
        h, e := b.Header()
        assert(t, e != nil, "parsed an invalid header: ", h)
    }
    assertHeader(failBlob1)
    assertHeader(failBlob2)
    assertHeader(failBlob3)
    assertHeader(failBlob4)
    assertHeader(failBlob5)
    assertHeader(failBlob6)
    assertHeader(failBlob7)
    assertHeader(failBlob8)
    assertHeader(failBlob9)
}

func TestPayloadFail(t *testing.T) {
    assertPayload := func(b *Blob) {
        _, e := b.Payload()
        assert(t, e != nil, "parsed an invalid payload: ", e)
    }
    assertPayload(failBlob10)
    assertPayload(failBlob11)
}
