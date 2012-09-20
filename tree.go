package ggit

import (
    "bufio"
    "fmt"
    "io"
)

// ================================================================= //
// TREE
// ================================================================= //

type Tree struct {
    entries []*TreeEntry
    repo    Repository
    size    int
}

// TODO: is this necessary?
func (t *Tree) Entries() []*TreeEntry {
    return t.entries
}

func (t *Tree) Type() ObjectType {
    return ObjectTree
}

func (t *Tree) Size() int {
    return t.size
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

// ================================================================= //
// TREE ENTRY
// ================================================================= //

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

// ================================================================= //
// TREE PARSING
// ================================================================= //

func parseTree(repo Repository, h *objectHeader, buf *bufio.Reader) (*Tree, error) {
    p := dataParser{buf}
    t := &Tree{
        entries: make([]*TreeEntry, 0),
        repo:    repo,
    }
    t.entries = make([]*TreeEntry, 0, 10)
    err := dataParse(func() {
        for !p.EOF() {
            mode := p.ParseFileMode(SP)
            name := p.ReadString(NUL)
            oid := p.ParseObjectIdBytes()
            entry := &TreeEntry{
                mode:  mode,
                otype: deduceObjectType(mode),
                name:  name,
                oid:   oid,
            }
            t.entries = append(t.entries, entry)
        }
    })
    t.repo = repo
    t.size = h.Size
    return t, err
}
