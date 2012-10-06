package api

import (
	"errors"
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
// TODO: remove this
func (c *Commit) FirstParent() *ObjectId {
	if len(c.parents) > 0 {
		return c.parents[0]
	}
	return nil
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

	// parse author
	c.author = p.parseWhoWhen(markerAuthor)
	p.ConsumeByte(LF)

	// parse committer
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

// TODO: the return values of this method are broken
func (f *Format) Commit(c *Commit) (int, error) {
	// tree
	fmt.Fprintf(f.Writer, "tree %\n", c.tree)

	// parents
	for _, p := range c.parents {
		fmt.Fprintf(f.Writer, "parent %s\n", p)
	}

	// author
	fmt.Fprintf(f.Writer, "author %\n", c.author)
	fmt.Fprintf(f.Writer, "committer %\n", c.committer)

	// commit message
	fmt.Fprintf(f.Writer, "\n%s", c.message)
	return 0, nil // TODO TODO
}

// ================================================================= //
// COMMIT OPERATIONS
// ================================================================= //

// ParentCommit selects the n-th parent commit of the given oid, which
// should point to a commit object. If n == 0, then this is considered
// to be the oid itself.
func ParentCommit(repo Repository, oid *ObjectId, n int) (*ObjectId, error) {
	if n == 0 {
		return oid, nil
	}
	o, err := ObjectFromOid(repo, oid)
	if err != nil {
		return nil, err
	}
	if o.Type() != ObjectCommit {
		return nil, errors.New("wrong object type")
	}
	p := o.(*Commit).parents
	if n > len(p) {
		return nil, fmt.Errorf("Parent %d doesn't exit", n)
	}
	return p[n], nil
}
