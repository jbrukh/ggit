package ggit

import (
    "fmt"
    "io"
)

type rawCommit struct {
    RawObject
}

func newRawCommit(obj *RawObject) *rawCommit {
    return &rawCommit{
        RawObject: *obj,
    }
}

func (rc *rawCommit) ParseCommit() (c *Commit, err error) {
    p, err := rc.Payload()
    if err != nil {
        return
    }
    // TODO implement the parsing
    fmt.Println(string(p))
    return new(Commit), nil // TODO
}

type Commit struct {
    author    string // TODO: make this struct with time
    committer string // TODO: make this struct with time
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
