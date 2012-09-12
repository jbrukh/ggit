package ggit

import (
    "io"
)

type Blob struct {
    /*
    	Unlike trees and commits, the Blob object does not have a rawBlob
    	object intermediately because there is nothing to parse.
    */
    RawObject
    repo *Repository
}

func (b *Blob) String() string {
    p, _ := b.Payload()
    return string(p)
}

func (b *Blob) Type() ObjectType {
    return OBJECT_BLOB
}

func (b *Blob) WriteTo(w io.Writer) (n int, err error) {
    return io.WriteString(w, b.String())
}
