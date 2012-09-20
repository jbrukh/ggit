package ggit

import (
    "bufio"
    "fmt"
    "io"
)

type FileMode uint16

const (
    // TODO: there are more modes
    MODE_DLTD FileMode = 0000000
    MODE_FILE FileMode = 0100644
    MODE_EXEC FileMode = 0100755
    MODE_TREE FileMode = 0040000
    MODE_LINK FileMode = 0120000
)

type Tree struct {
    entries []*TreeEntry
    repo    Repository
}

// TODO: is this necessary?
func (t *Tree) Entries() []*TreeEntry {
    return t.entries
}

func (t *Tree) Type() ObjectType {
    return ObjectTree
}

func (t *Tree) WriteTo(w io.Writer) (n int, err error) {
    for _, e := range t.entries {
        s := fmt.Sprintf("%.6o %s %-43s %s\n", e.mode, e.otype, e.oid, e.name)
        N, err := io.WriteString(w, s)
        if err != nil {
            break // TODO: is the error above shadowing err??
        }
        n += N
    }
    return
}

type TreeEntry struct {
    mode  FileMode
    otype ObjectType
    name  string
    oid   *ObjectId
}

func (e *TreeEntry) String() (s string) {
    s = fmt.Sprintf("%.6o %s %-43s %s", e.mode, e.otype, e.oid, e.name)
    return
}

func parseTree(repo Repository, buf *bufio.Reader) (*Tree, error) {
    p := dataParser{buf}
    t := &Tree{
        entries: make([]*TreeEntry, 0),
        repo:    repo,
    }
    dataParse(func() {
        mode := p.ReadString(SP)
        println(mode)
    })
    return t, nil // TODO
}

/*
func toTree(repo Repository, obj *RawObject) (t *Tree, err error) {
    p, err := obj.Payload()
    if err != nil {
        return
    }
    entries := make([]*TreeEntry, 0, 10)
    for len(p) > 0 {
        e, size, err := parseTreeEntry(p)
        if err != nil {
            return nil, err
        }
        entries = append(entries, e)
        p = p[size:]
    }
    t = &Tree{
        entries,
        repo,
    }
    return
}

func parseTreeEntry(p []byte) (e *TreeEntry, size int, err error) {
    const MAX_SZ = 64
    l := min(MAX_SZ, len(p))
    size = 0
    for i := 0; i < l; i++ {
        if p[i] == ' ' {
            modeStr := string(p[:i])
            // fmt.Printf("modeStr:\t%s\n", modeStr)
            for j := i; j < l; j++ {
                if p[j] == NUL {
                    fileName := string(p[i+1 : j])
                    // fmt.Printf("fileName:\t%s\n", fileName)
                    j++ // skip the null
                    size = j + OID_SZ
                    if size > l {
                        err = parseErr("not enough bytes for hash")
                        return
                    }
                    hsh := p[j:size]
                    // fmt.Printf("hash:\t%s\n", NewObjectIdFromBytes(hsh))
                    e, err = toTreeEntry(modeStr, fileName, hsh)
                    return
                }
            }
        }
    }
    err = errors.New("malformed object")
    return
}

// getTreeEntry converts raw data for the entry into a TreeEntry object.
func toTreeEntry(modeStr, fileName string, hsh []byte) (e *TreeEntry, err error) {
    mode, err := strconv.ParseInt(modeStr, 8, 32)
    if err != nil {
        return
    }
    m := FileMode(mode)
    var oid *ObjectId
    if oid, err = NewObjectIdFromBytes(hsh); err != nil {
        return nil, err
    }
    e = &TreeEntry{
        mode:  m,
        otype: deduceObjectType(m),
        name:  fileName,
        oid:   oid,
    }
    return
}
*/
