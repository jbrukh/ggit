package ggit

import (
    "bufio"
    "io"
)

// ================================================================= //
// CONSTANTS RELATED TO TYPES
// ================================================================= //

// the types of Git objects
type ObjectType string

// return a human-readable representation of an ObjectType
func (otype ObjectType) String() string {
    return string(otype)
}

const (
    ObjectBlob   ObjectType = "blob"
    ObjectTree   ObjectType = "tree"
    ObjectCommit ObjectType = "commit"
    ObjectTag    ObjectType = "tag"
)

var objectTypes []string = []string{
    ObjectBlob.String(),
    ObjectTree.String(),
    ObjectCommit.String(),
    ObjectTag.String(),
}

// ================================================================= //
// OBJECTS AND RAWOBJECTS
// ================================================================= //

type Object interface {
    Type() ObjectType

    // write the string representation of 
    // this object to the writer
    WriteTo(w io.Writer) (n int, err error)
}

// ObjectHeader is the deserialized (and more efficiently stored)
// version of a git object header
type objectHeader struct {
    Type ObjectType
    Size int
}

func parseObjectHeader(buf *bufio.Reader) (*objectHeader, error) {
    h := new(ObjectHeader)
    p := dataParser{buf}
    err := dataParse(func() {
        h.Type = ObjectType(p.ConsumeStrings(objectTypes))
        p.ConsumeByte(SP)
        h.Size = p.ParseInt(NUL)
    })
    return h, err
}

/*
func parseObjectHeader(b []byte) (typeStr, sizeStr string, pInx uint) {
    const MAX_SZ = 32
    var i, j uint
    l := uint(min(MAX_SZ, len(b)))
    for i = 0; i < l; i++ {
        if b[i] == SP {
            typeStr = string(b[:i])
            for j = i; j < l; j++ {
                if b[j] == NUL {
                    pInx = j
                    sizeStr = string(b[i+1 : j])
                    return
                }
            }
        }
    }
    return
}
*/
