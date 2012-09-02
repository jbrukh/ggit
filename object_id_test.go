package ggit

import (
    "crypto/sha1"
    "io"
    "testing"
)

const (
    ZEROS  = "0000000000000000000000000000000000000000"
    ONES   = "1111111111111111111111111111111111111111"
    ALLSET = "ffffffffffffffffffffffffffffffffffffffff"
    CRAZY  = "abcdef1234567890000000000000000000000000"
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

    crazy := make([]byte, OID_SZ)
    crazy[0] = 0xAB
    crazy[1] = 0xCD
    crazy[2] = 0xEF
    crazy[3] = 0x12
    crazy[4] = 0x34
    crazy[5] = 0x56
    crazy[6] = 0x78
    crazy[7] = 0x90
    compareHexRepr(t, crazy, CRAZY)
}

func compareHexRepr(t *testing.T, bytes []byte, expected string) {
    id := NewObjectIdFromBytes(bytes)
    repr := id.String()
    if repr != expected {
        t.Error("representation is not correct, expected ", expected, " but got ", repr)
    }
}

func TestHex2Bytes(t *testing.T) {
    testHex2Bytes(t, ZEROS)
    testHex2Bytes(t, ONES)
    testHex2Bytes(t, ALLSET)
}

func testHex2Bytes(t *testing.T, hex string) {
    testBytes, err := computeBytes(hex)
    if err != nil {
        t.Error("failed with error: ", err)
    }
    converted := computeRepr(testBytes)
    if converted != hex[0:OID_HEXSZ] {
        t.Error("could not convert to bytes successfully: ", hex, " => ", converted) 
    }
}

func TestNewObjectIdFromString(t *testing.T) {
    id, err := NewObjectIdFromString(CRAZY)
    if err != nil || id == nil || id.bytes == nil {
        t.Error("did not initialize bytes properly")
    }
}

func TestNewObjectIdFromHash(t *testing.T) {
    h := sha1.New()
    id := NewObjectIdFromHash(h)
    if id.bytes == nil {
        t.Error("did not initialize bytes properly")
    }
    if id.String() != ZEROS {
        t.Error("representation should be 0's")
    }

    // now we will test a real hash
    io.WriteString(h, "I have always known that one day I would take this road, but yesterday I did not know it would be today.")
    hashBytes := make([]byte, 20)
    h.Sum(hashBytes)
    id = NewObjectIdFromHash(h)

    expected, actual := computeRepr(hashBytes), id.String()
    if actual != expected {
        t.Error("bad hash initialization: ", expected, " but got ", actual)
    }
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
