//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
trees.go implements ggit Tree objects, TreeEntries, their parsing and
formatting.
*/
package api

import (
	"fmt"
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/util"
)

// ================================================================= //
// TREE
// ================================================================= //

type Tree struct {
	entries []*TreeEntry
	hdr     *objects.ObjectHeader
	oid     *objects.ObjectId
}

// TODO: is this necessary?
func (t *Tree) Entries() []*TreeEntry {
	return t.entries
}

func (t *Tree) Header() *objects.ObjectHeader {
	return t.hdr
}

func (t *Tree) ObjectId() *objects.ObjectId {
	return t.oid
}

// ================================================================= //
// TREE ENTRY
// ================================================================= //

type TreeEntry struct {
	mode  objects.FileMode
	otype objects.ObjectType
	name  string
	oid   *objects.ObjectId
}

func (e *TreeEntry) Mode() objects.FileMode {
	return e.mode
}

func (e *TreeEntry) ObjectId() *objects.ObjectId {
	return e.oid
}

func (e *TreeEntry) ObjectType() objects.ObjectType {
	return e.otype
}

func (e *TreeEntry) Name() string {
	return e.name
}

func (e *TreeEntry) String() (s string) {
	s = fmt.Sprintf("%.6o %s %s\t%s", e.mode, e.otype, e.oid, e.name)
	return
}

// ================================================================= //
// PARSING
// ================================================================= //

// parseTree performs the parsing of binary data into a Tree
// object, or panics with panicErr if there is a problem parsing.
// For this reason, it should be called as a parameter to
// safeParse().
func (p *objectParser) parseTree() *Tree {
	t := &Tree{
		entries: make([]*TreeEntry, 0), // TODO
		oid:     p.oid,
	}
	p.ResetCount()
	for !p.EOF() {
		mode := p.ParseFileMode(SP)
		name := p.ReadString(NUL)
		oid := p.ParseOidBytes()
		entry := &TreeEntry{
			mode:  mode,
			otype: deduceObjectType(mode),
			name:  name,
			oid:   oid,
		}
		t.entries = append(t.entries, entry)
	}
	t.hdr = p.hdr
	if p.Count() != p.hdr.Size() {
		util.PanicErrf("payload of size %d isn't of expected size %d", p.Count(), p.hdr.Size())
	}
	return t
}

// The file mode of a tree entry implies an object type.
func deduceObjectType(mode objects.FileMode) objects.ObjectType {
	switch mode {
	case objects.ModeNew, objects.ModeBlob, objects.ModeBlobExec, objects.ModeLink:
		return objects.ObjectBlob
	case objects.ModeTree:
		return objects.ObjectTree
	}
	// TODO
	panic("unknown mode")
}

// ================================================================= //
// FORMATTING
// ================================================================= //

// Tree formats this tree object into an API-friendly string that is
// the same as the output of git-cat-file tree <tree>.
func (f *Format) Tree(t *Tree) (int, error) {
	N := 0
	for _, e := range t.entries {
		n, err := fmt.Fprintf(f.Writer, "%o %s%s%s", e.mode, e.name, string(NUL), string(e.oid.Bytes()))
		N += n
		if err != nil {
			return N, err
		}
	}
	return N, nil
}

// TreePretty formats this tree object into a human-friendly table
// that is the same as the output of git-cat-file -p <tree>.
func (f *Format) TreePretty(t *Tree) (int, error) {
	N := 0
	for _, e := range t.entries {
		n, err := fmt.Fprintf(f.Writer, "%.6o %s %s\t%s\n", e.mode, e.otype, e.oid, e.name)
		N += n
		if err != nil {
			return N, err
		}
	}
	return N, nil
}
