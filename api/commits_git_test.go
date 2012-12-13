//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import (
	"github.com/jbrukh/ggit/api/format"
	"github.com/jbrukh/ggit/api/objects"
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

	f := format.NewStrFormat()
	for _, detail := range info.Commits {
		o, err := repo.ObjectFromOid(objects.OidNow(detail.CommitOid))
		util.AssertNoErr(t, err)

		// check the id
		util.Assert(t, o.ObjectId().String() == detail.CommitOid)

		// check the header
		util.Assert(t, o.Header().Type() == objects.ObjectCommit)
		util.AssertEqualInt(t, int(o.Header().Size()), detail.CommitSize)

		// now convert to a commit and check the fields
		var cmt *objects.Commit
		util.AssertPanicFree(t, func() {
			cmt = o.(*objects.Commit)
		})

		// check the tree
		util.Assert(t, cmt.Tree() != nil)
		util.AssertEqualString(t, cmt.Tree().String(), detail.TreeOid)

		// check the whole representation, which will catch
		// most of the other stuff
		f.Reset()
		f.Object(o)
		util.AssertEqualString(t, detail.CommitRepr, f.String())
	}
}
