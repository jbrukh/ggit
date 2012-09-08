package ggit

import (
    "testing"
	"fmt"
)

func Test_toObjectType(t *testing.T) {
    test := func(s string, T ObjectType) {
        if tp, _ := toObjectType(s); tp != T {
            t.Error("mismatch")
        }
    }
    test("blob", OBJECT_BLOB)
    test("tree", OBJECT_TREE)
    test("tag",  OBJECT_TAG)
    test("commit", OBJECT_COMMIT)
}

func Test_toObjectHeader(t *testing.T) {
    const P1 = "blob 11"
    const P2 = "commit 10323"
    const P3 = "tree 19\000"
    const F1 = "commit"
    const F2 = "\000"
    const F3 = "   "
    const F4 = "hedgehog 11\000"
	const F5 = ""

    testOk := func(header string, otype ObjectType, osize int) {
        h, err := toObjectHeader(header)
        if err != nil {
            t.Error("gave error: ", err)
        }
        if h.Type != otype || h.Size != osize {
            t.Error("mismatch, expecting ", otype, " ", osize)
        }
    }
    testFail := func(header string) {
        h, err := toObjectHeader(header)
        if err == nil || h != nil {
            t.Error("should have failed on: ", header)
        }
    }
    testOk(P1, OBJECT_BLOB, 11)
    testOk(P2, OBJECT_COMMIT, 10323)
    testOk(P3, OBJECT_TREE, 19)
    testFail(F1)
    testFail(F2)
    testFail(F3)
    testFail(F4)
    testFail(F5)
}

func Test_Payload(t *testing.T) {
	const P1 = "blob 17\00012345678901234567"
	testOk := func(payload string) {
		bts := fmt.Sprintf("blob %d\000%s", len(payload), payload)
		r := RawObject {
			bytes: []byte(bts),
		}
		h, err := r.Header()
		if err != nil {
			t.Error("could not parse: ", err)
		}
		if h.Type != OBJECT_BLOB || h.Size != len(payload) {
			t.Error("wrong type or size")
		}
		p, err := r.Payload()
		if string(p) != payload {
			t.Errorf("wrong payload; got: '%s'", string(p))
		}
	}
	testOk("haha")
	testOk("")
	testOk(" ")
	testOk("          ")
	testOk("13123241324342342341242341234123421")
}
