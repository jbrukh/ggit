//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import (
	"github.com/jbrukh/ggit/test"
	"github.com/jbrukh/ggit/util"
	"testing"
)

// Test_readCommits will compare the commit output of
// git and ggit for a string of commits.
func Test_readCommits(t *testing.T) {
	testCase := test.Linear
	repo := Open(testCase.Repo())
	info := testCase.Info().(*test.InfoLinear)

	util.Assert(t, info.N > 1)
	util.Assert(t, len(info.Commits) == info.N)

	f := NewStrFormat()
	for _, c := range info.Commits {
		o, err := repo.ObjectFromOid(OidNow(c.CommitOid))
		util.AssertNoErr(t, err)

		// check the id
		util.Assert(t, o.ObjectId().String() == c.CommitOid)

		// check the header
		util.Assert(t, o.Header().Type() == ObjectCommit)
		util.AssertEqualInt(t, int(o.Header().Size()), c.Size)

		// now convert to a commit and check the fields
		var cmt *Commit
		util.AssertPanicFree(t, func() {
			cmt = o.(*Commit)
		})

		// check the tree
		util.Assert(t, cmt.Tree() != nil)
		util.AssertEqualString(t, cmt.Tree().String(), c.TreeOid)

		// check the whole representation, which will catch
		// most of the other stuff
		f.Reset()
		f.Object(o)
		util.AssertEqualString(t, c.Repr, f.String())
	}
}
