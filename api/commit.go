package api

import (
	"fmt"
	"io"
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
// GGIT COMMIT OBJECT
// ================================================================= //

type Commit struct {
	author    *WhoWhen
	committer *WhoWhen
	message   string
	tree      *ObjectId
	parents   []*ObjectId
	repo      Repository
	size      int
}

func (c *Commit) Type() ObjectType {
	return ObjectCommit
}

func (c *Commit) Size() int {
	return c.size
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

// ================================================================= //
// OBJECT PARSER COMMIT PARSING METHODS
// ================================================================= //

func (p *objectParser) parseCommit() *Commit {

	c := new(Commit)
	p.ResetCount()

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

	c.author = p.parseWhoWhen(markerAuthor)
	p.ConsumeByte(LF)

	c.committer = p.parseWhoWhen(markerCommitter)
	p.ConsumeByte(LF)

	// commit message
	p.ConsumeByte(LF)
	c.message = p.String()

	c.size = p.hdr.Size

	if p.Count() != p.hdr.Size {
		panicErr("payload doesn't match prescibed size")
	}

	return c
}
