package ggit

import (
    "bufio"
    "io"
)

// Blob represents the deserialized version of a Git blob
// object.
type Blob struct {
    data []byte
    size int
    repo Repository
}

// parseBlob parses the payload of a binary blob object
// and converts it to Blob
func parseBlob(repo Repository, h *objectHeader, buf *bufio.Reader) (*Blob, error) {
    p := dataParser{buf}
    b := new(Blob)
    err := dataParse(func() {
        data := p.Bytes()
        if len(data) != h.Size {
            panicErrf("wrong payload size, expecting: %d", h.Size)
        }
        b.data = data
    })
    b.repo = repo
    b.size = h.Size
    return b, err
}

func (b *Blob) String() string {
    return string(b.data)
}

func (b *Blob) Type() ObjectType {
    return ObjectBlob
}

func (b *Blob) Size() int {
    return b.size
}

func (b *Blob) WriteTo(w io.Writer) (n int, err error) {
    return io.WriteString(w, b.String())
}
