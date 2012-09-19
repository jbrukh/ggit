package ggit

import (
    "fmt"
    "strings"
    "testing"
)

func Test_toObjectType(t *testing.T) {
    test := func(s string, T ObjectType) {
        if tp, _ := toObjectType(s); tp != T {
            t.Error("mismatch")
        }
    }
    test("blob", ObjectBlob)
    test("tree", ObjectTree)
    test("tag", ObjectTag)
    test("commit", ObjectCommit)
}

func Test_Parse(t *testing.T) {
    const P1 = "blob 11\000  2 dfow aj"
    const P2 = "commit 10\000 1 2 3 4 0"
    const P3 = "tree 19\000hello world!!!!!!!!"
    const P4 = "blob 0\000"
    const P5 = "tag 10\000sex or die"
    const F1 = "commit"
    const F2 = "\000"
    const F3 = "   "
    const F4 = "hedgehog 11\000"
    const F5 = ""

    testOk := func(data string, otype ObjectType) {
        fmt.Println("testing: ", data)
        rawObj := RawObject{
            bytes: []byte(data),
        }
        _, p, err := rawObj.Parse()
        if err != nil {
            t.Error("gave error: ", err)
        }
        toks := strings.Split(data, string(NUL))
        pld := toks[1]
        if pld != string(p) {
            t.Error("parsed wrong payload: ", p)
        }
    }
    testFail := func(data string) {
        fmt.Println("testing: ", data)
        rawObj := RawObject{
            bytes: []byte(data),
        }
        _, _, err := rawObj.Parse()
        if err == nil {
            t.Error("should have failed")
        }
    }
    testOk(P1, ObjectBlob)
    testOk(P2, ObjectCommit)
    testOk(P3, ObjectTree)
    testOk(P4, ObjectBlob)
    testOk(P5, ObjectTag)
    testFail(F1)
    testFail(F2)
    testFail(F3)
    testFail(F4)
    testFail(F5)
}

func Test_Payload(t *testing.T) {
    testOk := func(payload string) {
        bts := fmt.Sprintf("blob %d\000%s", len(payload), payload)
        r := RawObject{
            bytes: []byte(bts),
        }
        h, err := r.Header()
        if err != nil {
            t.Error("could not parse: ", err)
        }
        if h.Type != ObjectBlob || h.Size != len(payload) {
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

func Test_PayloadFirst(t *testing.T) {
    const P1 = "blob 17\00012345678901234567"
    r := RawObject{
        bytes: []byte(P1),
    }
    p, err := r.Payload()
    if err != nil {
        t.Error("error: ", err)
    }
    if string(p) != "12345678901234567" {
        t.Error("wrong payload")
    }
}
