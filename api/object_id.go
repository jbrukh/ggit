//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import (
	"errors"
	"fmt"
	"hash"
)

const (
	OidSize    = 20          // bytes
	OidHexSize = OidSize * 2 // maximum length of hex string we can translate
)

// ================================================================= //
// OBJECT ID
// ================================================================= //

// ObjectId represents a git SHA1 hash that is 
// used to represent objects and allows conversion
// between the binary and string versions of the
// id. ObjectIds are known colloquially as "oids".
type ObjectId struct {
	bytes []byte
	repr  string
}

// OidFromBytes creates a new ObjectId from a byte slice. 
// Bytes are filled in from left to right, with no regard
// for the number of bytes in the input. Extra bytes are
// discarded and missing bytes are padded with zeros.
func OidFromBytes(bytes []byte) (id *ObjectId, err error) {
	if len(bytes) < OidSize {
		return nil, errors.New("not enough bytes for oid")
	}
	id = &ObjectId{
		bytes: make([]byte, OidSize),
	}
	copy(id.bytes, bytes)
	return
}

// OidFromArray convers an array of bytes into an ObjectId
// stored in binary form. Because array size is fixed at 
// compile time, this method does not throw an error.
func OidFromArray(bytes [OidSize]byte) (id *ObjectId) {
	oid, _ := OidFromBytes(bytes[:]) // no error can happen
	return oid
}

// OidFromString creates an ObjectId from a string representation
// of the hash. The length of the string should be OidHexSize, and
// must consist of the characters [a-zA-Z0-9] or else an error is
// returned.
func OidFromString(hex string) (id *ObjectId, err error) {
	id = &ObjectId{
		bytes: make([]byte, OidSize),
	}
	_, err = fmt.Sscanf(hex, "%x", &id.bytes)
	return
}

func OidFromHash(h hash.Hash) (id *ObjectId) {
	hsh := h.Sum(nil)
	id = &ObjectId{
		bytes: hsh[0:OidSize], // TODO: what if size exceeds hash?
	}
	return
}

func OidNow(correctHex string) *ObjectId {
	oid, err := OidFromString(correctHex)
	if err != nil {
		panic("provide a correct oid")
	}
	return oid
}

// String returns the hex string that represents
// the ObjectId bytes
func (id *ObjectId) String() string {
	if id.repr == "" {
		id.repr = fmt.Sprintf("%x", id.bytes)
	}
	return id.repr
}

// ================================================================= //
// PARSING
// ================================================================= //

// objectIdParser is a dataparser that supports parsing of oids.
type objectIdParser struct {
	dataParser
}

// ParseObjectId reads the next OidHexSize bytes from the
// Reader and places the resulting object id in oid.
func (p *objectIdParser) ParseObjectId() *ObjectId {
	hex := string(p.consume(OidHexSize))
	oid, e := OidFromString(hex)
	if e != nil {
		panicErrf("expected: hex string of size %d", OidHexSize)
	}
	return oid
}

func (p *objectIdParser) ParseObjectIdBytes() *ObjectId {
	b := p.consume(OidSize)
	oid, e := OidFromBytes(b)
	if e != nil {
		panicErrf("expected: hash bytes %d long", OidSize)
	}
	return oid
}

// ================================================================= //
// FORMATTING
// ================================================================= //

func (f *Format) ObjectId(oid *ObjectId) (int, error) {
	return fmt.Fprint(f.Writer, oid.String())
}
