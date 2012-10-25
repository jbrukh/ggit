//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
blob.go implements ggit Blob objects and their parsing and formatting.
*/
package api

import (
	"fmt"
)

// ================================================================= //
// BLOB OBJECT
// ================================================================= //

// Blob represents the deserialized version
// of a Git blob object.
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
// and converts it to Blob. If there are parsing errors,
// it panics with parseErr, so this method should be
// called as a parameter a safeParse().
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

// Blob formats the contents of the blog as a string
// for output to the screen.
func (f *Format) Blob(b *Blob) (int, error) {
	return fmt.Fprintf(f.Writer, "%s", string(b.data))
}
