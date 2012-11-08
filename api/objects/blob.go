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
package objects

// ================================================================= //
// BLOB OBJECT
// ================================================================= //

// Blob represents the deserialized version
// of a Git blob object.
type Blob struct {
	data []byte
	oid  *ObjectId
	hdr  *ObjectHeader
}

func NewBlob(oid *ObjectId, hdr *ObjectHeader, data []byte) *Blob {
	return &Blob{data, oid, hdr}
}

func (b *Blob) Header() *ObjectHeader {
	return b.hdr
}

func (b *Blob) ObjectId() *ObjectId {
	return b.oid
}

func (b *Blob) Data() []byte {
	return b.data
}
