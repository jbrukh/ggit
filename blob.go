package ggit

import (
    "bufio"
    "io"
)

// Blob represents the deserialized version of a Git blob
// object.
type Blob struct {
    data []byte
    repo Repository
}

// parseBlob parses the payload of a binary blob object
// and converts it to Blob
func parseBlob(repo Repository, h *objectHeader, buf *bufio.Reader) (*Blob, error) {
    p := dataParser{buf}
    b := new(Blob)
    err := dataParse(func() {
        b.data = p.Bytes()
    })
    b.repo = repo
    return b, err
}

func (b *Blob) String() string {
    return string(b.data)
}

func (b *Blob) Type() ObjectType {
    return ObjectBlob
}

func (b *Blob) WriteTo(w io.Writer) (n int, err error) {
    return io.WriteString(w, b.String())
}
