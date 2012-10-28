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

func Test_revParse__firstParent(t *testing.T) {
	testCase := test.Linear
	repo := Open(testCase.Repo())
	output := testCase.Output().(*test.OutputCommits)

	util.Assert(t, output.N > 1)
	util.Assert(t, len(output.Commits) == output.N)

	// test the first, parentless commit
	testParentlessCommit(t, repo, output.Commits[0].Oid)
	for _, c := range output.Commits[1:] {
		testShortOid(t, repo, c.Oid)
		testFirstParent(t, repo, c.Oid, c.ParentOid)
	}
}

func Test_revParse__secondAncestor(t *testing.T) {
	testCase := test.Linear
	repo := Open(testCase.Repo())
	output := testCase.Output().(*test.OutputCommits)

	util.Assert(t, output.N > 2)
	util.Assert(t, len(output.Commits) == output.N)

	// test the first, parentless commit
	for i, c := range output.Commits[2:] {
		testSecondAncestor(t, repo, c.Oid, output.Commits[i].Oid)
	}
}

func testShortOid(t *testing.T, repo Repository, oid string) {
	for i := 4; i <= 40; i++ {
		o, err := ObjectFromRevision(repo, oid[:i])
		util.AssertNoErr(t, err)
		util.AssertEqualString(t, o.ObjectId().String(), oid)
	}
}

func testFirstParent(t *testing.T, repo Repository, oid string, parentOid string) {
	testCommit(t, repo, oid)

	for i := 4; i <= 40; i++ {
		parent, err := ObjectFromRevision(repo, oid[:i]+"^")
		util.AssertNoErr(t, err)
		util.Assert(t, parent.Header().Type() == ObjectCommit)
		util.AssertEqualString(t, parent.ObjectId().String(), parentOid)
	}
}

func testSecondAncestor(t *testing.T, repo Repository, oid string, parentOid string) {
	testCommit(t, repo, oid)

	for i := 4; i <= 40; i++ {
		parent, err := ObjectFromRevision(repo, oid[:i]+"~2")
		util.AssertNoErr(t, err)
		util.Assert(t, parent.Header().Type() == ObjectCommit)
		util.AssertEqualString(t, parent.ObjectId().String(), parentOid)
	}
}

func testParentlessCommit(t *testing.T, repo Repository, oid string) {
	testCommit(t, repo, oid)
	_, err := ObjectFromRevision(repo, oid+"~1")
	util.Assert(t, err != nil)
	_, err = ObjectFromRevision(repo, oid+"^")
	util.Assert(t, err != nil)
}

func testCommit(t *testing.T, repo Repository, oid string) {
	o, err := ObjectFromRevision(repo, oid)
	util.AssertNoErr(t, err)
	util.Assert(t, o.Header().Type() == ObjectCommit)
	util.Assert(t, o.ObjectId().String() == oid)
}
