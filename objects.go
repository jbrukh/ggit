package ggit

import (
    "errors"
    "io"
    "strconv"
)

type ObjectType int

// the types of objects
const (
    OBJECT_BLOB ObjectType = iota
    OBJECT_TREE
    OBJECT_COMMIT
    OBJECT_TAG
)

const (
    OBJECT_BLOB_STR   = "blob"
    OBJECT_TREE_STR   = "tree"
    OBJECT_COMMIT_STR = "commit"
    OBJECT_TAG_STR    = "tag"
)

type Object interface {
    // Type returns the type of the object
    Type() ObjectType

    // write the string representation to the writer
    WriteTo(w io.Writer) (n int, err error)
}

// raw (but uncompressed) data for a
// git object that contains the header;
type RawObject struct {
    bytes []byte
    pInx  uint // start of payload bytes
    h     *ObjectHeader
}

// ObjectHeader is the deserialized (and more efficiently stored)
// version of a git object header
type ObjectHeader struct {
    Type ObjectType
    Size int
}

func toObjectType(typeStr string) (otype ObjectType, err error) {
    switch typeStr {
    case OBJECT_BLOB_STR:
        return OBJECT_BLOB, nil
    case OBJECT_TREE_STR:
        return OBJECT_TREE, nil
    case OBJECT_TAG_STR:
        return OBJECT_TAG, nil
    case OBJECT_COMMIT_STR:
        return OBJECT_COMMIT, nil
    }
    return 0, errors.New("unknown object type")
}

func (otype ObjectType) String() string {
    switch otype {
    case OBJECT_BLOB:
        return OBJECT_BLOB_STR
    case OBJECT_TREE:
        return OBJECT_TREE_STR
    case OBJECT_COMMIT:
        return OBJECT_COMMIT_STR
    case OBJECT_TAG:
        return OBJECT_TAG_STR
    }
    panic("unknown type")
}

// parses the header from the raw data
func (obj *RawObject) Header() (h *ObjectHeader, err error) {
    if obj.h == nil {
        obj.h, err = obj.parse()
        return obj.h, err
    }
    return obj.h, nil
}

func (obj *RawObject) parse() (h *ObjectHeader, err error) {
    if len(obj.bytes) < 1 {
        return nil, errors.New("no data bytes")
    }
    var typeStr, sizeStr string
    typeStr, sizeStr, obj.pInx = parseHeader(obj.bytes)
    if obj.pInx <= 0 {
        return nil, errors.New("bad header")
    }
    otype, err := toObjectType(typeStr)
    if err != nil {
        return
    }
    osize, err := strconv.Atoi(sizeStr)
    if err != nil {
        return nil, errors.New("bad object size")
    }
    return &ObjectHeader{otype, osize}, nil
}

func parseHeader(b []byte) (typeStr, sizeStr string, pInx uint) {
    const MAX_SZ = 32
    var i, j uint
    l := uint(min(MAX_SZ, len(b)))
    for i = 0; i < l; i++ {
        if b[i] == ' ' {
            typeStr = string(b[:i])
            for j = i; j < l; j++ {
                if b[j] == '\000' {
                    pInx = j
                    sizeStr = string(b[i+1 : j])
                    return
                }
            }
        }
    }
    return
}

// returns the headerless payload of the object
func (o *RawObject) Payload() (bts []byte, err error) {
    if o.pInx <= 0 {
        // must parse the header
        if _, err = o.Header(); err != nil {
            return
        }
    }
    return o.bytes[o.pInx+1:], nil
}

func (o *RawObject) Parse() (h *ObjectHeader, payload []byte, err error) {
    if h, err = o.Header(); err != nil {
        return
    }

    if payload, err = o.Payload(); err != nil {
        return
    }

    // check size!
    if h.Size != len(payload) {
        err = errors.New("object corrupted (checksize is wrong)")
    }
    return
}

func (o *RawObject) Write(b []byte) (n int, err error) {
    if o.bytes == nil {
        o.bytes = make([]byte, len(b))
        return copy(o.bytes, b), nil
    }
    return 0, errors.New("object already has data")
}

// returns the raw byte representation of
// the object
func (o *RawObject) Bytes() []byte {
    return o.bytes
}
