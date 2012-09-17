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
    parents   []*ObjectId
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

func parseCommit(b []byte) (c *Commit, err error) {
    buf := bytes.NewBuffer(b)
    c = new(Commit)

    marker, oid, e := parseOidLine(buf)
    if e != nil || marker != OBJECT_TREE_STR {
        return nil, parseErr("wrong marker")
    }
    c.tree = oid

    // c.parents = make([]*ObjectId, 0, 2)

    //    for {
    //        marker, oid, e := parseOidLine(buf, "parent") // TODO: const
    //        if 
    //    }
    return
}

func parseOidLine(buf *bytes.Buffer) (marker string, oid *ObjectId, err error) {
    var m, oidStr string
    n, e := fmt.Fscanf(buf, "%s %s\n", &m, &oidStr)
    if e != nil || n != 2 {
        return "", oid, parseErr("could not parse oid line")
    }
    oid, err = NewObjectIdFromString(oidStr)
    return m, oid, err
}

func parseHex(buf *bytes.Buffer, delim byte) (oid *ObjectId, err error) {
    oidStr, e := nextToken(buf, delim)
    if e == nil {
        if oid, err = NewObjectIdFromString(oidStr); err != nil {
            return
        }
    }
    return nil, e
}
