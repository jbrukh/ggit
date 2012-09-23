package ggit

import (
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

func (p *objectParser) parseTree() *Tree {
	t := &Tree{
		entries: make([]*TreeEntry, 0), // TODO
		repo:    nil,
	}
	p.ResetRead()
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

	if p.read != p.hdr.Size {
		// panic
	}

	return t
}
