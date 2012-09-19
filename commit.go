package ggit

import (
    "bufio"
    "bytes"
    "fmt"
    "io"
)

const (
    markerTree      = "tree"
    markerParent    = "parent"
    markerAuthor    = "author"
    markerCommitter = "committer"
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

func (c *Commit) String() string {
    const FMT = "Commit{author=%s, committer=%s, tree=%s, parents=%v, message='%s'}"
    return fmt.Sprintf(FMT, c.author, c.committer, c.tree, c.parents, c.message)
}

func (c *Commit) WriteTo(w io.Writer) (n int, err error) {
    return io.WriteString(w, c.String())
}

func (c *Commit) addParent(oid *ObjectId) {
    if c.parents == nil {
        c.parents = make([]*ObjectId, 0, 2)
    }
    c.parents = append(c.parents, oid)
}

func toCommit(repo Repository, obj *RawObject) (c *Commit, err error) {
    var p []byte
    if p, err = obj.Payload(); err != nil {
        return
    }
    if c, err = parseCommit(p); err != nil {
        fmt.Println("could not parse: ", err)
        return
    }
    return c, nil // TODO
}

func parseCommit(b []byte) (c *Commit, err error) {
    buf := bufio.NewReader(bytes.NewBuffer(b))
    c = new(Commit)

    p := &dataParser{buf}
    err = dataParse(func() {

        // read the tree line
        p.ConsumeString(markerTree)
        p.ConsumeByte(SP)
        c.tree = p.ParseObjectId()
        p.ConsumeByte(LF)

        // read an arbitrary number of parent lines
        n := len(markerParent)
        for p.PeekString(n) == markerParent {
            p.ConsumeString(markerParent)
            p.ConsumeByte(SP)
            c.addParent(p.ParseObjectId())
            p.ConsumeByte(LF)
        }

        // read the author
        p.ConsumeString(markerAuthor)
        p.ConsumeByte(SP)
        line := p.ReadString(LF)                  // gets rid of the LF!
        c.author = &AuthorTimestamp{line, "", ""} // TODO

        // read the committer
        p.ConsumeString(markerCommitter)
        p.ConsumeByte(SP)
        line = p.ReadString(LF)                      // gets rid of the LF!
        c.committer = &AuthorTimestamp{line, "", ""} // TODO

        // read the commit message
        p.ConsumeByte(LF)
        c.message = p.String()
    })
    return
}

func isParentMarker(buf *bufio.Reader) (bool, error) {
    peek, err := buf.Peek(len(markerParent))
    if err != nil {
        return false, err
    }
    return string(peek) == markerParent, nil
}

func parseOidLine(buf *bufio.Reader) (marker string, oid *ObjectId, err error) {
    var m, oidStr string
    _, e := fmt.Fscanf(buf, "%s %s\n", &m, &oidStr)
    if e != nil {
        return "", nil, parseErrn("could not parse oid line: ", e.Error())
    }
    oid, err = NewObjectIdFromString(oidStr)
    return m, oid, err
}

func parseAuthorTimestamp(buf *bufio.Reader) (string, *AuthorTimestamp, error) {
    var marker, name, email, date string
    if _, e := fmt.Fscanf(buf, "%s %s <%s> %s\n", &marker, &name, &email, &date); e != nil {
        return "", nil, e
    }
    return marker, &AuthorTimestamp{name, email, date}, nil
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
