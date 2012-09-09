package ggit

import (
	"io"
)

type Blob struct {
    RawObject
    parent    *Repository
}

func (b *Blob) String() string {
	p, _ := b.Payload()
	return string(p)
}

func (b *Blob) Type() ObjectType {
	return OBJECT_BLOB
}

func (b *Blob) WriteTo(w io.Writer) (n int, err error) {
	return io.WriteString(w, string(b.String()))
}