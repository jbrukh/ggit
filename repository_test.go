package ggit

import (
	"testing"
	"bytes"
)

func TestObjectPath(t *testing.T) {
	oid, _ := NewObjectIdFromString("aaaaabbbbbcccccddddd11111222223333344444")
	path := objectPath(oid)
	if path != "objects/aa/aaabbbbbcccccddddd11111222223333344444" {
		t.Error("Wrong path:", path)
	}
}

func TestBuffer(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.Write([]byte("aaa\000bbb"))
	
	first, _ := buf.ReadString('\000')
	second := string(buf.Bytes())
	println("first: ", len(first), " ", first)
	println("second: ", len(second), " ", second)
}