//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
commits.go implements ggit Commit objects, their parsing and formatting,
and useful operations that allow users to resolve and navigate commits.
*/
package api

import (
	"errors"
	"fmt"
	"github.com/jbrukh/ggit/api/objects"
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
// GGIT COMMIT OBJECT
// ================================================================= //

type Commit struct {
	hdr       *objects.ObjectHeader
	oid       *objects.ObjectId
	treeOid   *objects.ObjectId
	parents   []*objects.ObjectId
	author    *objects.WhoWhen
	committer *objects.WhoWhen
	message   string
}

func (c *Commit) Header() *objects.ObjectHeader {
	return c.hdr
}

func (c *Commit) ObjectId() *objects.ObjectId {
	return c.oid
}

func (c *Commit) Tree() *objects.ObjectId {
	return c.treeOid
}

func (c *Commit) Parents() []*objects.ObjectId {
	return c.parents
}

func (c *Commit) Author() *objects.WhoWhen {
	return c.author
}

func (c *Commit) Committer() *objects.WhoWhen {
	return c.committer
}

func (c *Commit) Message() string {
	return c.message
}

func (c *Commit) addParent(oid *objects.ObjectId) {
	c.parents = append(c.parents, oid)
}

// ================================================================= //
// OBJECT PARSER COMMIT PARSING METHODS
// ================================================================= //

func (p *objectParser) parseCommit() *Commit {
	c := &Commit{
		parents: make([]*objects.ObjectId, 0),
		oid:     p.oid,
	}
	p.ResetCount()

	// read the tree line
	p.ConsumeString(markerTree)
	p.ConsumeByte(SP)
	c.treeOid = p.ParseOid()
	p.ConsumeByte(LF)

	// read an arbitrary number of parent lines
	n := len(markerParent)
	for p.PeekString(n) == markerParent {
		p.ConsumeString(markerParent)
		p.ConsumeByte(SP)
		c.addParent(p.ParseOid())
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

	c.hdr = p.hdr
	if p.Count() != p.hdr.Size() {
		util.PanicErr("payload doesn't match prescibed size")
	}
	return c
}

// ================================================================= //
// OBJECT FORMATTER
// ================================================================= //

// TODO: the return values of this method are broken
func (f *Format) Commit(c *Commit) (int, error) {
	// tree
	fmt.Fprintf(f.Writer, "tree %s\n", c.treeOid)

	// parents
	for _, p := range c.parents {
		fmt.Fprintf(f.Writer, "parent %s\n", p)
	}

	// author
	sf := NewStrFormat()
	sf.WhoWhen(c.author)
	fmt.Fprintf(f.Writer, "author %s\n", sf.String())
	sf.Reset()
	sf.WhoWhen(c.committer)
	fmt.Fprintf(f.Writer, "committer %s\n", sf.String())

	// commit message
	fmt.Fprintf(f.Writer, "\n%s", c.message)
	return 0, nil // TODO TODO
}

// ================================================================= //
// COMMIT OPERATIONS
// ================================================================= //

func CommitNthParent(repo Repository, c *Commit, n int) (rc *Commit, err error) {
	if n == 0 {
		return c, nil
	}
	l := len(c.parents)
	if 0 < n && n <= l {
		oid := c.parents[n-1]
		return CommitFromOid(repo, oid)
	}
	return nil, fmt.Errorf("cannot find parent n=%d", n)
}

// CommitNthAncestor will look up a chain of n objects by 
// following the first parent. If n == 0, then the parameterized
// commit is returned. If along the way, a commit is found
// to not have a first parent, an error is returned.
func CommitNthAncestor(repo Repository, c *Commit, n int) (rc *Commit, err error) {
	rc = c
	for i := 0; i < n; i++ {
		if len(rc.parents) > 0 {
			rc, err = CommitFromOid(repo, rc.parents[0])
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("no first parent")
		}
	}
	return rc, nil
}

// CommitFromObject returns the commit being referred to; that is, if
// the object is a commit object, it is converted and returned. If the
// object is a tag, then the target of the tag is returned. Other object
// types cause an error to be returned.
func CommitFromObject(repo Repository, o objects.Object) (*Commit, error) {
	switch t := o.(type) {
	case *Commit:
		return t, nil
	case *Tag:
		obj, err := ObjectFromOid(repo, t.Object())
		if err != nil {
			return nil, err
		}
		return obj.(*Commit), nil
	}
	return nil, errors.New("not a commit or tag")
}

// CommitFromOid takes an oid and turns it into a commit object. If the
// oid points at a commit, the Commit object is returned. If the oid 
// points at an annotated tag, then the target commit is returned. If
// the oid points to another type of object, an error is returned.
func CommitFromOid(repo Repository, oid *objects.ObjectId) (*Commit, error) {
	o, err := ObjectFromOid(repo, oid)
	if err != nil {
		return nil, err
	}
	return CommitFromObject(repo, o)
}

// CommitFromRef turns a full reference to be converted to the commit
// object. If the reference is a reference to a commit, the commit
// object is returned. If the reference is a tag, then the target commit
// of the tag is returned.
func CommitFromRef(repo Repository, spec string) (*Commit, error) {
	o, err := ObjectFromRef(repo, spec)
	if err != nil {
		return nil, err
	}
	return CommitFromObject(repo, o)
}
