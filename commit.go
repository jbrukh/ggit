package ggit

import (
    "bytes"
    "fmt"
    "io"
)

type Commit struct {
    author    *AuthorTimestamp
    committer *AuthorTimestamp
    message   string
    tree      *ObjectId
    parent    []*ObjectId
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

type BinaryParseErr string

func (p BinaryParseErr) Error() string {
    return string(p)
}

func parseCommit(b []byte) (c *Commit, err error) {
    buf := bytes.NewBuffer(b)

    tree, e := buf.ReadBytes(SP)
    if e != nil {
        return nil, e
    }
    treeStr := trimLastStr(tree)

    if treeStr != OBJECT_TREE_STR {
        return nil, BinaryParseErr("dd")
    }
    return
}
