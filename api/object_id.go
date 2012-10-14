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
func OidFromArray(bytes [20]byte) (id *ObjectId) {
	oid, _ := OidFromBytes(bytes[:]) // no error can happen
	return oid
}

// OidFromString creates an ObjectId from a string representation
// of the hash. The length of the string should be OidHexSize, and
// must consist of the characters [a-zA-Z0-9] or else an error is
// returned.
func OidFromString(hex string) (id *ObjectId, err error) {
	bytes, e := computeBytes(hex)
	if e == nil {
		id = &ObjectId{
			bytes: bytes,
		}
	}
	return id, e
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
		id.repr = computeRepr(id.bytes)
	}
	return id.repr
}

// computes the hex string representation of
// the object id
func computeRepr(bytes []byte) (hex string) {
	const byte2hex = "0123456789abcdef"
	out := make([]byte, OidHexSize)
	for inx, b := range bytes {
		// the left and right halves of the byte (8 bits)
		i := 2 * inx
		out[i] = byte2hex[int(b>>4)]
		out[i+1] = byte2hex[int(b&0xf)]
	}
	return string(out)
}

func computeBytes(hex string) (bytes []byte, err error) {
	if len(hex) < OidHexSize {
		err = errors.New("hex is too short")
		return
	}
	bytes = make([]byte, OidSize)
	for inx, _ := range bytes[0:OidSize] {
		i := 2 * inx
		left, err := hex2byte(hex[i])
		right, err := hex2byte(hex[i+1])
		if err != nil {
			return nil, err
		}
		bytes[inx] = left<<4 | right
	}
	return
}

func hex2byte(ch byte) (byte, error) {
	// TODO: find out if go map is faster (almost certainly not);
	// testing should be done on a high level, though
	switch ch {
	case '0':
		return 0x0, nil
	case '1':
		return 0x1, nil
	case '2':
		return 0x2, nil
	case '3':
		return 0x3, nil
	case '4':
		return 0x4, nil
	case '5':
		return 0x5, nil
	case '6':
		return 0x6, nil
	case '7':
		return 0x7, nil
	case '8':
		return 0x8, nil
	case '9':
		return 0x9, nil
	case 'a', 'A':
		return 0xA, nil
	case 'b', 'B':
		return 0xB, nil
	case 'c', 'C':
		return 0xC, nil
	case 'd', 'D':
		return 0xD, nil
	case 'e', 'E':
		return 0xE, nil
	case 'f', 'F':
		return 0xF, nil
	}
	return 0x0, errors.New("unknown char")
}

// ================================================================= //
// PARSING
// ================================================================= //

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
