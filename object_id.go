package ggit

import (
    "hash"
)

const (
    OID_SZ     = 20           // bytes
    OID_HEXSZ  = OID_SZ*2     // maximum length of hex string we can translate
)

const (
    toHex = "0123456789abcdef"
)

type ObjectId struct {
    bytes []byte
    hex string
}

// create a new ObjectId from bytes; bytes are filled
// in from left to right, with no regard for the number
// of bytes in the input. Extra bytes are discarded and
// missing bytes are padded with zeros.
func NewObjectIdFromBytes(bytes []bytes) *ObjectId {
    if len(bytes) < OID_SZ {
        // TODO: decide if error
    }
    id := &ObjectId{
        bytes: make([]byte, OID_SZ),
    }
    copy(id.bytes, bytes)
    id.computeHex()
    return id;
}

func NewObjectIdFromString(hex string) *ObjectId {
    // TODO
    return nil
}

func NewObjectIdFromHash(h hash.Hash) *ObjectId {
    // TODO
    return nil
}

func (id *ObjectId) String() string {
    if id.hex == "" {
        id.hex = id.computeHex()
    }
    return id.hex
}

func computeHex(id *ObjectId) (hex string){
    out := make([]byte, 40)
    for inx, b := range id.bytes {
        // the left and right halves of the byte (8 bits)
        out[2*inx] = toHex[int(b >> 4)]
        out[2*inx+1] = toHex[int(b & 0xf)]
    }
    return string(out)
}
