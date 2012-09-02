package ggit

import (
    "hash"
)

const (
    OID_SZ     = 20           // bytes
    OID_HEXSZ  = OID_SZ*2     // maximum length of hex string we can translate
)

type ObjectId struct {
    id []byte
}

func NewObjectIdFromString(hex string) *ObjectId {
    // TODO
    return nil
}

func NewObjectIdFromHash(h hash.Hash) *ObjectId {
    // TODO
    return nil
}

func (id *ObjectId) String() string {
    // TODO
    return ""
}
