package ggit

import (
        "bytes"
        "crypto/sha1"
        "errors"
        "hash"
        "strconv"
        "strings"
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
}

type ObjectHeader struct {
        Type    ObjectType
        Size    int
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

func toObjectHeader(header string) (h *ObjectHeader, err error) {
        var (
                toks    []string
                otype   ObjectType
        )
        if toks := strings.Split(header, " "); len(toks) < 2 {
                return nil, errors.New("bad object header")
        }

        typeStr, sizeStr := toks[0], toks[1]
        if otype, err = toObjectType(typeStr); err != nil {
                return
        }

        osize, err := strconv.Atoi(sizeStr)
        if err != nil {
                return nil, errors.New("bad object size")
        }

        return &ObjectHeader{otype, osize}, nil
}

// interface for hashable objects
type Hashable interface {
        Bytes() []byte
}

// parses the header from the raw data
func (o *RawObject) Header() (h *ObjectHeader, err error) {
        buf := bytes.NewBuffer(o.bytes)
        var header string
        if header, err = buf.ReadString('\000'); err != nil {
                return nil, errors.New("can't find end of header")
        }

        h, err = toObjectHeader(header)
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

// returns the headerless payload of the object
func (o *RawObject) Payload() (bts []byte, err error) {
        buf := bytes.NewBuffer(o.bytes)
        if _, err = buf.ReadString('\000'); err != nil {
                return
        }
        return buf.Bytes(), nil
}

// the hash object used to build
// hashes of our objects
var sha hash.Hash = sha1.New()

// produce a hash for any object that
// can be construed as a bunch of bytes
func Hash(h Hashable) (o *ObjectId) {
        sha.Reset()
        sha.Write(h.Bytes())
        return NewObjectIdFromHash(sha)
}

type Blob struct {
        RawObject
}

type Commit struct {
        RawObject
}
