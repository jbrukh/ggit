package ggit

import (
    "errors"
    "io"
    "strconv"
)

// ================================================================= //
// CONSTANTS RELATED TO TYPES
// ================================================================= //

// the types of Git objects
type ObjectType int8

// return a human-readable representation of an ObjectType
// TODO: turn this into a to-function
func (otype ObjectType) String() string {
    switch otype {
    case ObjectBlob:
        return objectBlobStr
    case ObjectTree:
        return objectTreeStr
    case ObjectCommit:
        return objectCommitStr
    case ObjectTag:
        return objectTagStr
    }
    panic("unknown type")
}

const (
    ObjectBlob ObjectType = iota
    ObjectTree
    ObjectCommit
    ObjectTag
)

// string representations of Git objects
const (
    objectBlobStr   = "blob"
    objectTreeStr   = "tree"
    objectCommitStr = "commit"
    objectTagStr    = "tag"
)

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
type ObjectHeader struct {
    Type ObjectType
    Size int
}

// raw (but uncompressed) data for a
// git object that contains the header;
type RawObject struct {
    bytes []byte
    pInx  uint // start of payload bytes
    h     *ObjectHeader
}

func (obj *RawObject) parse() (h *ObjectHeader, err error) {
    if len(obj.bytes) < 1 {
        return nil, parseErr("no data bytes")
    }
    var typeStr, sizeStr string
    typeStr, sizeStr, obj.pInx = parseObjectHeader(obj.bytes)
    if obj.pInx <= 0 {
        return nil, parseErr("bad header")
    }
    otype, e := toObjectType(typeStr)
    if e != nil {
        return nil, e
    }
    osize, e := strconv.Atoi(sizeStr)
    if e != nil {
        return nil, parseErr("bad object size")
    }
    return &ObjectHeader{otype, osize}, nil
}

// parses the header from the raw data
func (obj *RawObject) Header() (h *ObjectHeader, err error) {
    // the header is parsed lazily, and then cached
    if obj.h == nil {
        obj.h, err = obj.parse()
        return obj.h, err
    }
    return obj.h, nil
}

// returns the headerless payload of the object
func (obj *RawObject) Payload() (bts []byte, err error) {
    if obj.pInx <= 0 {
        // must parse the header
        if _, err = obj.Header(); err != nil {
            return
        }
    }
    bts = obj.bytes[obj.pInx+1:]
    if obj.h.Size != len(bts) {
        err = errors.New("object corrupted (checksize is wrong)")
    }
    return
}

func (obj *RawObject) Parse() (h *ObjectHeader, payload []byte, err error) {
    // Payload() will automatically cache the header
    if payload, err = obj.Payload(); err != nil {
        return
    }
    return obj.h, payload, nil
}

func (obj *RawObject) Write(b []byte) (n int, err error) {
    if obj.bytes == nil {
        obj.bytes = make([]byte, len(b))
        return copy(obj.bytes, b), nil
    }
    return 0, errors.New("object already has data")
}

// returns the raw byte representation of
// the object
func (o *RawObject) Bytes() []byte {
    return o.bytes
}

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

func toObjectType(typeStr string) (otype ObjectType, err error) {
    switch typeStr {
    case objectBlobStr:
        return ObjectBlob, nil
    case objectTreeStr:
        return ObjectTree, nil
    case objectTagStr:
        return ObjectTag, nil
    case objectCommitStr:
        return ObjectCommit, nil
    }
    return 0, errors.New("unknown object type")
}
