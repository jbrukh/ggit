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

// interface for hashable objects
type Hashable interface {
    Bytes() []byte
}

// parses the header from the raw data
// TODO: reconsider the return values here
func (o *RawObject) Header() (header string, size int, err error) {
    buf := bytes.NewBuffer(o.bytes)
    header, err = buf.ReadString('\000')
    if err != nil {
        err = errors.New("can't find end of header")
        return
    }
    toks := strings.Split(header, " ")
    // TODO: fix this garbage
    if len(toks) > 1 {
        if size, err = strconv.Atoi(toks[1]); err != nil {
            return "", 0, errors.New("malformed size in header")
        }
        if size != len(buf.Bytes()) {
            err = errors.New("header size doesn't match payload size")
        }
        otype := toks[1]
        if otype != "blob" && otype != "tree" && otype != "commit" {
            err = errors.New("unknown otype in header")
        }
    } else {
        err = errors.New("missing size in header")
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

type Tree struct {
    RawObject
}

type Commit struct {
    RawObject
}


