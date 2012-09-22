package ggit

import (
	"bytes"
	"crypto/sha1"
	"hash"
)

const (
	NUL   = '\000'
	SP    = ' '
	LF    = '\n'
	LT    = '<'
	GT    = '>'
	PLUS  = "+"
	MINUS = "-"
)

var signs []string = []string{
	PLUS,
	MINUS,
}

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

func sgn(a int) int {
	if a < 0 {
		return -1
	}
	return 1
}

// The file mode of a tree entry implies an object type.
func deduceObjectType(mode FileMode) ObjectType {
	switch mode {
	case ModeDeleted, ModeFile, ModeExecutable:
		return ObjectBlob
	case ModeTree:
		return ObjectTree
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
