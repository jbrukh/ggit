package ggit

import (
    "crypto/sha1"
    "hash"
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

    // return a representatin of this object as
    // a sequence of bytes that are ready to be
    // stored in the object database, as well as
    // the cryptographic id that is associated
    // with this representation
    Bytes() (id *ObjectId, bytes []byte)
}

type ObjectDatabase interface {
    
}

func digest(bytes []byte) (id *ObjectId) {
    shaHash.Reset()
    shaHash.Write(bytes)
    id = NewObjectIdFromHash(shaHash)
    return
}
