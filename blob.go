package ggit

type Blob struct {
    // bytes stores the data contained in a blob
    bytes []byte
}

func (b *Blob) Type() ObjectType {
    return OBJECT_BLOB
}

func (b *Blob) Bytes() (id *ObjectId, bytes []byte) {
    
}


