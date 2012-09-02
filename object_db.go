package ggit

type ObjectType int

// the types of objects
const (
    OBJECT_BLOB ObjectType = iota
    OBJECT_TREE
    OBJECT_COMMIT
    OBJECT_TAG
)

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
