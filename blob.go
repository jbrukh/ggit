package ggit

import (
    "bytes"
    "fmt"
    "io"
)

type Blob struct {
    // bytes stores the data contained in a blob
    bytes []byte
}

func NewBlob(input string) (b *Blob) {
    return &Blob{
        bytes: []byte(input),
    }
}

func (b *Blob) Type() ObjectType {
    return OBJECT_BLOB
}

func (b *Blob) WriteTo(w io.Writer) (id *ObjectId, err error) {
    var buf bytes.Buffer
    header := fmt.Sprintf("%s %d\000", OBJECT_BLOB, len(b.bytes))
    buf.Write([]byte(header))
    buf.Write(b.bytes)

    shaHash.Reset()
    mw := io.MultiWriter(w, shaHash)
    _, err = buf.WriteTo(mw)
    if err == nil {
        id = NewObjectIdFromHash(shaHash)
    }
    return
}


