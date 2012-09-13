package ggit

import (
    "fmt"
    "io"
)

func toCommit(obj *RawObject) (c *Commit, err error) {
    p, err := obj.Payload()
    if err != nil {
        return
    }
    // TODO implement the parsing
    fmt.Println(string(p))
    return new(Commit), nil // TODO
}

type Commit struct {
    author    *AuthorTimestamp // TODO: make this struct with time
    committer *AuthorTimestamp // TODO: make this struct with time
    message   string
    tree      *ObjectId
    parent    *ObjectId
    repo      *Repository
}

func (c *Commit) Type() ObjectType {
    return OBJECT_COMMIT
}

func (c *Commit) WriteTo(w io.Writer) (n int, err error) {
    // TODO
    return 0, nil
}
