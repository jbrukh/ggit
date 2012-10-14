package api

import (
	"fmt"
)

// ================================================================= //
// BLOB OBJECT
// ================================================================= //

// Blob represents the deserialized version of a Git blob
// object.
type Blob struct {
	data []byte
	oid  *ObjectId
	hdr  ObjectHeader
}

func (b *Blob) Header() ObjectHeader {
	return b.hdr
}

func (b *Blob) ObjectId() *ObjectId {
	return b.oid
}

func (b *Blob) String() string {
	return string(b.data)
}

// ================================================================= //
// OBJECT PARSER
// ================================================================= //

// parseBlob parses the payload of a binary blob object
// and converts it to Blob
func (p *objectParser) parseBlob() *Blob {
	b := new(Blob)
	b.oid = p.oid

	p.ResetCount()

	b.data = p.Bytes()
	b.hdr = p.hdr

	if p.Count() != p.hdr.Size() {
		panicErr("payload doesn't match prescibed size")
	}

	return b
}

// ================================================================= //
// OBJECT FORMATTER
// ================================================================= //

func (f *Format) Blob(b *Blob) (int, error) {
	return fmt.Fprintf(f.Writer, string(b.data))
}
