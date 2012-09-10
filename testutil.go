package ggit

import (
    "testing"
)

func assert(t *testing.T, b bool, items ...interface{}) {
    if !b {
        t.Error("error: ", items)
    }
}

func newBlobFromString(blobStr string) *Blob {
    return &Blob{
        RawObject: RawObject{
            bytes: []byte(blobStr),
        },
        repo: nil,
    }
}
