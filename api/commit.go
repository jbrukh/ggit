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
	oid       *ObjectId
}

func (c *Commit) Type() ObjectType {
	return ObjectCommit
}

func (c *Commit) Size() int {
	return c.size
}

func (c *Commit) ObjectId() *ObjectId {
	return c.oid
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

func (p *objectParser) parseCommit(oid *ObjectId) *Commit {
	c := new(Commit)
	c.oid = oid
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

func CommitNthParent(repo Repository, c *Commit, n int) (rc *Commit, err error) {
	if n == 0 {
		return c, nil
	}
	l := len(c.parents)
	if 0 < n && n < l {
		oid := c.parents[n]
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
			return CommitFromOid(repo, rc.parents[0])
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
func CommitFromObject(repo Repository, o Object) (*Commit, error) {
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
func CommitFromOid(repo Repository, oid *ObjectId) (*Commit, error) {
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
