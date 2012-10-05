package api

import (
	"fmt"
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
// OBJECT PARSER
// ================================================================= //

func (p *objectParser) parseTree() *Tree {
	t := &Tree{
		entries: make([]*TreeEntry, 0), // TODO
		repo:    nil,
	}
	p.ResetCount()
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
	t.size = p.hdr.Size
	if p.Count() != p.hdr.Size {
		panicErr("payload doesn't match prescibed size")
	}
	return t
}

// ================================================================= //
// OBJECT FORMATTER
// ================================================================= //

func (f *Format) Tree(t *Tree) (int, error) {
	N := 0
	for _, e := range t.entries {
		n, err := fmt.Fprintf(f.Writer, "%.6o %s %-43s %s\n", e.mode, e.otype, e.oid, e.name)
		N += n
		if err != nil {
			return N, err
		}
	}
	return N, nil
}
