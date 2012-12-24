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
package parse

import (
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/api/token"
	"github.com/jbrukh/ggit/util"
)

// ================================================================= //
// PARSING
// ================================================================= //

// parseTree performs the parsing of binary data into a Tree
// object, or panics with panicErr if there is a problem parsing.
// For this reason, it should be called as a parameter to
// safeParse().
func (p *objectParser) parseTree() *objects.Tree {
	entries := make([]*objects.TreeEntry, 0)
	p.ResetCount()
	for !p.EOF() {
		mode := p.ParseFileMode(token.SP)
		name := p.ReadString(token.NUL)
		oid := p.ParseOidBytes()
		t := deduceObjectType(mode)
		entry := objects.NewTreeEntry(mode, t, name, oid)
		entries = append(entries, entry)
	}

	if p.Count() != p.hdr.Size() {
		util.PanicErrf("payload of size %d isn't of expected size %d", p.Count(), p.hdr.Size())
	}
	return objects.NewTree(p.oid, p.hdr, entries)
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
