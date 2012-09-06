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
