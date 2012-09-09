package ggit

import (
    "errors"
    "fmt"
    "io"
    "strconv"
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
    repo    *Repository
}

// TODO: is this necessary?
func (t *Tree) Entries() []*TreeEntry {
    return t.entries
}

func (t *Tree) Type() ObjectType {
    return OBJECT_TREE
}

func (t *Tree) WriteTo(w io.Writer) (n int, err error) {
    for _, e := range t.entries {
        s := fmt.Sprintf("%.6o %s %-43s %s\n", e.mode, e.otype, e.oid, e.name)
        N, err := io.WriteString(w, s)
        if err != nil {
            break
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

type rawTree struct {
    RawObject
}

func newRawTree(rawObj *RawObject) (rt *rawTree) {
    // TODO: decide if object check should be done here
    return &rawTree{
        RawObject: *rawObj,
    }
}

func (rt *rawTree) ParseTree() (t *Tree, err error) {
    p, err := rt.Payload()
    if err != nil {
        return
    }
    entries := make([]*TreeEntry, 0, 10)
    for len(p) > 0 {
        e, size, err := parseEntry(p)
        if err != nil {
            return nil, err
        }
        entries = append(entries, e)
        p = p[size:]
    }
    t = &Tree{
        entries,
        nil,
    }
    return
}

func parseEntry(p []byte) (e *TreeEntry, size int, err error) {
    const MAX_SZ = 64
    l := min(MAX_SZ, len(p))
    size = 0
    for i := 0; i < l; i++ {
        if p[i] == ' ' {
            modeStr := string(p[:i])
            // fmt.Printf("modeStr:\t%s\n", modeStr)
            for j := i; j < l; j++ {
                if p[j] == '\000' {
                    fileName := string(p[i+1 : j])
                    // fmt.Printf("fileName:\t%s\n", fileName)
                    j++ // skip the null
                    size = j + OID_SZ
                    if size > l {
                        err = errors.New("not enough bytes for hash")
                        return
                    }
                    hsh := p[j:size]
                    // fmt.Printf("hash:\t%s\n", NewObjectIdFromBytes(hsh))
                    e, err = getTreeEntry(modeStr, fileName, hsh)
                    return
                }
            }
        }
    }
    err = errors.New("malformed object")
    return
}

func getTreeEntry(modeStr, fileName string, hsh []byte) (e *TreeEntry, err error) {
    mode, err := strconv.ParseInt(modeStr, 8, 32)
    if err != nil {
        return
    }
    m := FileMode(mode)
    e = &TreeEntry{
        mode:  m,
        otype: deduceObjectType(m), // TODO: fix this
        name:  fileName,
        oid:   NewObjectIdFromBytes(hsh),
    }
    return
}
