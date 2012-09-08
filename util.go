package ggit

import (
	"hash"
	"crypto/sha1"
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
func Hash(h Hashable) (o *ObjectId) {
    sha.Reset()
    sha.Write(h.Bytes())
    return NewObjectIdFromHash(sha)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}