package ggit

import (
    "hash"
)

const (
    OID_SZ     = 20           // bytes
    OID_HEXSZ  = OID_SZ*2     // maximum length of hex string we can translate
)


type ObjectId struct {
    bytes []byte
    repr string
}

// create a new ObjectId from bytes; bytes are filled
// in from left to right, with no regard for the number
// of bytes in the input. Extra bytes are discarded and
// missing bytes are padded with zeros.
func NewObjectIdFromBytes(bytes []byte) (id *ObjectId) {
    if len(bytes) < OID_SZ {
        // TODO: decide if error
    }
    id = &ObjectId{
        bytes: make([]byte, OID_SZ),
    }
    copy(id.bytes, bytes)
    return
}

func NewObjectIdFromString(hex string) (id *ObjectId) {
    // TODO
    return nil
}

func NewObjectIdFromHash(h hash.Hash) (id *ObjectId) {
    id = &ObjectId{
        bytes: make([]byte, 20),
    }
    h.Sum(id.bytes)
    return
}

// String returns the hex string that represents
// the ObjectId bytes
func (id *ObjectId) String() string {
    if id.repr == "" {
        id.repr = computeRepr(id)
    }
    return id.repr
}

func computeRepr(id *ObjectId) (hex string){
    const toHex = "0123456789abcdef"
    out := make([]byte, OID_HEXSZ)
    for inx, b := range id.bytes {
        // the left and right halves of the byte (8 bits)
        out[2*inx] = toHex[int(b >> 4)]
        out[2*inx+1] = toHex[int(b & 0xf)]
    }
    return string(out)
}
