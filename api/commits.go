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
)

// ================================================================= //
// COMMIT OPERATIONS
// ================================================================= //

func CommitNthParent(repo Repository, c *objects.Commit, n int) (rc *objects.Commit, err error) {
	if n == 0 {
		return c, nil
	}
	l := len(c.Parents())
	if 0 < n && n <= l {
		oid := c.Parents()[n-1]
		return CommitFromOid(repo, oid)
	}
	return nil, fmt.Errorf("cannot find parent n=%d", n)
}

// CommitNthAncestor will look up a chain of n objects by 
// following the first parent. If n == 0, then the parameterized
// commit is returned. If along the way, a commit is found
// to not have a first parent, an error is returned.
func CommitNthAncestor(repo Repository, c *objects.Commit, n int) (rc *objects.Commit, err error) {
	rc = c
	for i := 0; i < n; i++ {
		if len(rc.Parents()) > 0 {
			rc, err = CommitFromOid(repo, rc.Parents()[0])
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
func CommitFromObject(repo Repository, o objects.Object) (*objects.Commit, error) {
	switch t := o.(type) {
	case *objects.Commit:
		return t, nil
	case *objects.Tag:
		obj, err := ObjectFromOid(repo, t.Object())
		if err != nil {
			return nil, err
		}
		return obj.(*objects.Commit), nil
	}
	return nil, errors.New("not a commit or tag")
}

// CommitFromOid takes an oid and turns it into a commit object. If the
// oid points at a commit, the Commit object is returned. If the oid 
// points at an annotated tag, then the target commit is returned. If
// the oid points to another type of object, an error is returned.
func CommitFromOid(repo Repository, oid *objects.ObjectId) (*objects.Commit, error) {
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
func CommitFromRef(repo Repository, spec string) (*objects.Commit, error) {
	o, err := ObjectFromRef(repo, spec)
	if err != nil {
		return nil, err
	}
	return CommitFromObject(repo, o)
}
