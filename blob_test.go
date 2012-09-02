package ggit

import (
    "crypto/sha1"
    "fmt"
    "os"
    "testing"
)

func TestNewBlob(t *testing.T) {
    data := "this is a blob of sorts"
    expected := "blob 23\000this is a blob of sorts"

    b := NewBlob(data)
    if b.bytes == nil {
        t.Error("did not initialize blob")
    }
    id, err := b.WriteTo(os.Stdout)
    fmt.Println()
    if err != nil {
        t.Error("failed writing out")
    }
    h := sha1.New()
    h.Write([]byte(expected))
    expectedId := NewObjectIdFromHash(h)

    if id.String() != expectedId.String() {
        t.Error("expected hash doesn't match: ", id.String(), " should be ", expectedId)
    }
}
