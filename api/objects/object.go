//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package objects

// ObjectHeader is the deserialized (and more efficiently stored)
// version of a git object header
type ObjectHeader struct {
	otype ObjectType
	size  int64
}

func NewObjectHeader(t ObjectType, size int64) *ObjectHeader {
	return &ObjectHeader{t, size}
}

func (h *ObjectHeader) Type() ObjectType {
	return h.otype
}

func (h *ObjectHeader) Size() int64 {
	return h.size
}

// Object represents a generic git object: a blob, a tree,
// a tag, or a commit.
type Object interface {

	// Header returns the object header, which
	// contains the object's type and size.
	Header() *ObjectHeader

	// ObjectId returns the object id of the object.
	ObjectId() *ObjectId
}
