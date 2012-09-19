package ggit

import (
    "fmt"
    "io"
)

type Tag struct {
    repo *Repository
    //commit-SHA1
    object ObjectId
    //tag name
    tag string
    //author with timestamp
    tagger *PersonTimestamp
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

func toTag(repo Repository, obj *RawObject) (t *Tag, err error) {
    p, e := obj.Payload()
    if e != nil {
        return nil, e
    }
    // TODO implement the parsing
    fmt.Println(string(p))
    return new(Tag), nil // TODO
}
