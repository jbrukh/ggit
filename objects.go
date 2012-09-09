package ggit

import (
    "errors"
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

// raw (but uncompressed) data for a
// git object that contains the header;
type RawObject struct {
    bytes []byte
    pInx  uint // start of payload bytes 
}

type ObjectHeader struct {
    Type ObjectType
    Size int
}

func toObjectType(typeStr string) (otype ObjectType, err error) {
    switch typeStr {
    case "blob":
        return OBJECT_BLOB, nil
    case "tree":
        return OBJECT_TREE, nil
    case "tag":
        return OBJECT_TAG, nil
    case "commit":
        return OBJECT_COMMIT, nil
    }
    return 0, errors.New("unknown object type")
}

func (otype ObjectType) String() string {
	switch otype {
		case OBJECT_BLOB:
			return "blob"
		case OBJECT_TREE:
			return "tree"
		case OBJECT_COMMIT:
			return "commit"
		case OBJECT_TAG:
			return "tag"
	}
	panic("unknown type")
}

// parses the header from the raw data
func (o *RawObject) Header() (h *ObjectHeader, err error) {
    if len(o.bytes) < 1 {
        return nil, errors.New("no data bytes")
    }
    var typeStr, sizeStr string
    typeStr, sizeStr, o.pInx = parseHeader(o.bytes)
    if o.pInx <= 0 {
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