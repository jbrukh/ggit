package ggit

import (
    "errors"
    "strconv"
)

// ================================================================= //
// OBJECTS AND RAWOBJECTS
// ================================================================= //

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
