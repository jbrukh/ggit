package api

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"hash"
	"runtime"
	"testing"
)

const (
	NUL = '\000'
	SP  = ' '
	LF  = '\n'
	LT  = '<'
	GT  = '>'
)

const (
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

// The file mode of a tree entry implies an object type.
func deduceObjectType(mode FileMode) ObjectType {
	switch mode {
	case ModeNew, ModeBlob, ModeBlobExec:
		return ObjectBlob
	case ModeTree:
		return ObjectTree
	}
	// TODO
	panic("unknown mode")
}

// ================================================================= //
// METHODS FOR TESTING
// ================================================================= //

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
