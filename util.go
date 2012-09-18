package ggit

import (
    "bytes"
    "crypto/sha1"
    "hash"
)

const (
    NUL = '\000'
    SP  = ' '
    LF  = '\n'
)

// the hash object used to build
// hashes of our objects
var sha hash.Hash = sha1.New()

// interface for hashable objects
type Hashable interface {
    Bytes() []byte
}

// produce a hash for any object that
// can be construed as a bunch of bytes
func makeHash(h Hashable) (o *ObjectId) {
    sha.Reset()
    sha.Write(h.Bytes())
    return NewObjectIdFromHash(sha)
}

// get the first OID_SZ of the hash
func getHash(h hash.Hash) []byte {
    return h.Sum(nil)[0:OID_SZ]
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// The file mode of a tree entry implies an object type.
func deduceObjectType(mode FileMode) ObjectType {
    switch mode {
    case MODE_DLTD, MODE_FILE, MODE_EXEC:
        return OBJECT_BLOB
    case MODE_TREE:
        return OBJECT_TREE
    }
    // TODO
    panic("unknown mode")
}

func nextToken(buf *bytes.Buffer, delim byte) (tok string, err error) {
    line, e := buf.ReadBytes(delim)
    if e != nil {
        return "", e
    }
    return trimLastStr(line), nil
}
