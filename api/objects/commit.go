//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
commit.go implements ggit Commit objects, their parsing and formatting,
and useful operations that allow users to resolve and navigate commits.
*/
package objects

// ================================================================= //
// GGIT COMMIT OBJECT
// ================================================================= //

type Commit struct {
	hdr       *ObjectHeader
	oid       *ObjectId
	treeOid   *ObjectId
	parents   []*ObjectId
	author    *WhoWhen
	committer *WhoWhen
	message   string
}

func NewCommit(oid, tree *ObjectId, size int64, parents []*ObjectId, author, committer *WhoWhen, msg string) *Commit {
	return &Commit{
		&ObjectHeader{
			ObjectCommit,
			size,
		},
		oid,
		tree,
		parents,
		author,
		committer,
		msg,
	}
}

func (c *Commit) Header() *ObjectHeader {
	return c.hdr
}

func (c *Commit) ObjectId() *ObjectId {
	return c.oid
}

func (c *Commit) Tree() *ObjectId {
	return c.treeOid
}

func (c *Commit) Parents() []*ObjectId {
	return c.parents
}

func (c *Commit) Author() *WhoWhen {
	return c.author
}

func (c *Commit) Committer() *WhoWhen {
	return c.committer
}

func (c *Commit) Message() string {
	return c.message
}
