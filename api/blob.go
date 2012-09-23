package api

import (
	"io"
)

// ================================================================= //
// BLOB OBJECT
// ================================================================= //

// Blob represents the deserialized version of a Git blob
// object.
type Blob struct {
	data []byte
	size int
	repo Repository
}

func (b *Blob) String() string {
	return string(b.data)
}

func (b *Blob) Type() ObjectType {
	return ObjectBlob
}

func (b *Blob) Size() int {
	return b.size
}

func (b *Blob) WriteTo(w io.Writer) (n int, err error) {
	return io.WriteString(w, b.String())
}

// ================================================================= //
// OBJECT PARSER BLOB PARSING METHODS
// ================================================================= //

// parseBlob parses the payload of a binary blob object
// and converts it to Blob
func (p *objectParser) parseBlob() *Blob {
	b := new(Blob)
	p.ResetCount()

	b.data = p.Bytes()
	b.size = p.hdr.Size

	if p.Count() != p.hdr.Size {
		panicErr("payload doesn't match prescibed size")
	}

	return b
}
