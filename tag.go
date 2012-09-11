package ggit

import (
    "io"
)

type Tag struct {
    repo *Repository
    //commit-SHA1
    object ObjectId
    //tag name
    tag string
    //author with timestamp
    tagger *AuthorTimestamp
    //message
    message string
}

func (t *Tag) String() string {
    //TODO
    return ""
}

func (t *Tag) Type() ObjectType {
    return OBJECT_TAG
}

func (t *Tag) WriteTo(w io.Writer) (n int, err error) {
    return io.WriteString(w, t.String())
}
