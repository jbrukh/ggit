package ggit

import (
    "fmt"
    "io"
)

type Commit struct {
    author    *AuthorTimestamp 
    committer *AuthorTimestamp 
    message   string
    tree      *ObjectId
    parent    *ObjectId
    repo      Repository
}

func (c *Commit) Type() ObjectType {
    return OBJECT_COMMIT
}

func (c *Commit) WriteTo(w io.Writer) (n int, err error) {
    // TODO
    return 0, nil
}

func toCommit(repo Repository, obj *RawObject) (c *Commit, err error) {
    p, err := obj.Payload()
    if err != nil {
        return
    }
    // TODO implement the parsing
    fmt.Println(string(p))
    return new(Commit), nil // TODO
}
