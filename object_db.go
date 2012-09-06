package ggit

import (
    "crypto/sha1"
    "hash"
    "io"
)

type ObjectType int

var shaHash hash.Hash = sha1.New()

// the types of objects
const (
    OBJECT_BLOB ObjectType = iota
    OBJECT_TREE
    OBJECT_COMMIT
    OBJECT_TAG
)

func (t ObjectType) String() string {
    switch t {
    case OBJECT_BLOB:
        return "blob"
    case OBJECT_TREE:
        return "tree"
    case OBJECT_COMMIT:
        return "commit"
    case OBJECT_TAG:
        return "tag"
    default:
        panic("unknown type")
    }
    return ""
}

type Object interface {
    // return the type of the object
    // TODO: ascertain whether this is actually needed
    Type() ObjectType

    // writes a representatin of this object as
    // a sequence of bytes that are ready to be
    // stored in the object database to an io.Writer, and
    // the cryptographic id that is associated
    // with this representation is returned
    WriteTo(w io.Writer) (id *ObjectId, err error)
}

type ObjectDatabase interface {
	// read an object by id
    ReadObject(ObjectId) Object

	// write an object into the database
	WriteObject(Object) ObjectId
	
	// return the number of objects in the databse
	Size() int
}
