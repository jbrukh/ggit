package api

import (
	"fmt"
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

// FirstParent returns the first parent of the commit, or
// nil if no such parent exists.
func (c *Commit) FirstParent() *ObjectId {
	if len(c.parents) > 0 {
		return c.parents[0]
	}
	return nil
}

func (c *Commit) String() string {
	// TODO: move this to goddamn formatter
	const format = "tree %s\n%s\nauthor %s\ncommitter %s\n\n%s"
	parentsToString := func(p []*ObjectId) string {
		s := ""
		for i := 0; i < len(p); i++ {
			if i > 0 {
				s = s + "\n"
			}
			s = fmt.Sprintf("%sparent %s", s, p[i])
		}
		return s
	}
	return fmt.Sprintf(format, c.tree, parentsToString(c.parents), c.author, c.committer, c.message)
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

// ================================================================= //
// OBJECT FORMATTER
// ================================================================= //

func (f *Format) Commit(c *Commit) (int, error) {
	return fmt.Fprint(f.Writer, c.String())
}
