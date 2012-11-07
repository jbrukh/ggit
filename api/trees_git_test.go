//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
trees_git_test.go implements git-comparison tests for ggit tree parsing.
*/
package api

import (
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/test"
	"github.com/jbrukh/ggit/util"
	"testing"
)

// Test_readCommits will compare the commit output of
// git and ggit for a string of commits.
func Test_readTree(t *testing.T) {
	testCase := test.Tree
	repo := Open(testCase.Repo())
	info := testCase.Info().(*test.InfoTree)
	f := NewStrFormat()

	var (
		oid = objects.OidNow(info.TreeOid)
	)
	o, err := ObjectFromOid(repo, oid)
	util.AssertNoErr(t, err)

	// check the id
	util.Assert(t, o.ObjectId().String() == info.TreeOid)

	// check the header
	util.Assert(t, o.Header().Type() == objects.ObjectTree)
	util.AssertEqualInt(t, int(o.Header().Size()), info.TreeSize)

	// get the tree
	// now convert to a tag and check the fields
	var tree *Tree
	util.AssertPanicFree(t, func() {
		tree = o.(*Tree)
	})

	// check entries
	entries := tree.Entries()
	util.AssertEqualInt(t, info.N, len(entries))
	util.Assert(t, info.N > 2)

	// check a file
	file := entries[0]
	util.AssertEqualString(t, info.File1Oid, file.ObjectId().String())
	util.Assert(t, file.Mode() == ModeBlob)
	// TODO: add checks for type, etc.

	// check the output
	// check the whole representation, which will catch
	// most of the other stuff
	f.Reset()
	f.Object(o)
	util.AssertEqualString(t, info.TreeRepr, f.String())
}
