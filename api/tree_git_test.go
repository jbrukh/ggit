//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
tree_git_test.go implements git-comparison tests for ggit tree parsing.
*/
package api

import (
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

	var (
		oid = OidNow(info.TreeOid)
	)
	o, err := ObjectFromOid(repo, oid)
	util.AssertNoErr(t, err)

	// check the id
	util.Assert(t, o.ObjectId().String() == info.TreeOid)

	// check the header
	util.Assert(t, o.Header().Type() == ObjectTree)
	util.AssertEqualInt(t, int(o.Header().Size()), info.TreeSize)

}
