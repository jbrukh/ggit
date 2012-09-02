package ggit

import (
    "testing"
)

const (
    ZEROS  = "0000000000000000000000000000000000000000"
    ONES   = "1111111111111111111111111111111111111111"
    ALLSET = "ffffffffffffffffffffffffffffffffffffffff"
)

func TestObjectIdString(t *testing.T) {
    zeros := make([]byte, OID_SZ)
    compareHexRepr(t, zeros, ZEROS)

    ones := make([]byte, OID_SZ)
    for inx, _ := range ones {
       ones[inx] |= 0x11 
    }
    compareHexRepr(t, ones, ONES)

    allset := make([]byte, OID_SZ)
    for inx, _ := range allset {
        allset[inx] |= 0xff
    }
    compareHexRepr(t, allset, ALLSET)
}

func compareHexRepr(t *testing.T, bytes []byte, expected string) {
    id := NewObjectIdFromBytes(bytes)
    repr := id.String()
    if repr != expected {
        t.Error("representation is not correct, expected ", expected, " but got ", repr)
    }
}

func TestNewObjectIdFromString(t *testing.T) {
    // TODO
}

func TestNewObjectIdFromHash(t *testing.T) {
    // TODO
}

func TestNewObjectIdFromBytes(t *testing.T) {
    bytes := make([]byte, OID_SZ)
    id := NewObjectIdFromBytes(bytes)
    if id.bytes == nil {
        t.Error("did not initialize bytes properly")
    }
    if id.repr != "" {
        t.Error("prematurely initialized string repr")
    }
    id.String()
    if id.repr == "" {
        t.Error("lazy init of string repr didn't work")
    }
}

func TestHex2Byte(t *testing.T) {
    tester := func(hex byte, b byte) {
        v, e := hex2byte(hex)
        if e != nil || v != b {
            t.Error("didn't get the right result for ", hex)
        }
    }
    tester('a', 0xA)
    tester('B', 0xB)
    tester('7', 0x7)
}
