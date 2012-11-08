//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
tree.go implements ggit Tree objects, TreeEntries, their parsing and
formatting.
*/
package objects

import "fmt"

// ================================================================= //
// TREE
// ================================================================= //

type Tree struct {
	entries []*TreeEntry
	hdr     *ObjectHeader
	oid     *ObjectId
}

func NewTree(oid *ObjectId, header *ObjectHeader, entries []*TreeEntry) *Tree {
	return &Tree{
		entries: entries,
		hdr:     header,
		oid:     oid,
	}
}

// TODO: is this necessary?
func (t *Tree) Entries() []*TreeEntry {
	return t.entries
}

func (t *Tree) Header() *ObjectHeader {
	return t.hdr
}

func (t *Tree) ObjectId() *ObjectId {
	return t.oid
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

func NewTreeEntry(mode FileMode, objectType ObjectType, name string, objectId *ObjectId) *TreeEntry {
	return &TreeEntry{
		mode,
		objectType,
		name,
		objectId,
	}
}

func (e *TreeEntry) Mode() FileMode {
	return e.mode
}

func (e *TreeEntry) ObjectId() *ObjectId {
	return e.oid
}

func (e *TreeEntry) ObjectType() ObjectType {
	return e.otype
}

func (e *TreeEntry) Name() string {
	return e.name
}

func (e *TreeEntry) String() (s string) {
	s = fmt.Sprintf("%.6o %s %s\t%s", e.mode, e.otype, e.oid, e.name)
	return
}
