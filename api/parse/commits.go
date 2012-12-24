//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

package parse

import (
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/api/token"
	"github.com/jbrukh/ggit/util"
)

// ================================================================= //
// CONSTANTS
// ================================================================= //

const (
	markerTree      = "tree"
	markerParent    = "parent"
	markerAuthor    = "author"
	markerCommitter = "committer"
)

// ================================================================= //
// OBJECT PARSER COMMIT PARSING METHODS
// ================================================================= //

func (p *objectParser) parseCommit() *objects.Commit {
	parents := make([]*objects.ObjectId, 0)

	p.ResetCount()

	// read the tree line
	p.ConsumeString(markerTree)
	p.ConsumeByte(token.SP)
	treeOid := p.ParseOid()
	p.ConsumeByte(token.LF)

	// read an arbitrary number of parent lines
	n := len(markerParent)
	for p.PeekString(n) == markerParent {
		p.ConsumeString(markerParent)
		p.ConsumeByte(token.SP)
		parents = append(parents, p.ParseOid())
		p.ConsumeByte(token.LF)
	}

	// parse author
	author := p.parseWhoWhen(markerAuthor)
	p.ConsumeByte(token.LF)

	// parse committer
	committer := p.parseWhoWhen(markerCommitter)
	p.ConsumeByte(token.LF)

	// commit message
	p.ConsumeByte(token.LF)
	message := p.String()

	if p.Count() != p.hdr.Size() {
		util.PanicErr("payload doesn't match prescibed size")
	}

	return objects.NewCommit(p.oid, treeOid, p.hdr.Size(), parents, author, committer, message)
}
