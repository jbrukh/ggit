package ggit

type Blob struct {
    RawObject
    parent    *Repository
}

func (b *Blob) String() string {
	p, _ := b.Payload()
	return string(p)
}