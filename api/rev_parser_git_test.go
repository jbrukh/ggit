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
	testParentlessCommit(t, repo, OidNow(output.Commits[0].Oid))
	for _, c := range output.Commits[1:] {
		oid, expOid := OidNow(c.Oid), OidNow(c.ParentOid)
		testShortOid(t, repo, oid)
		testFirstParent(t, repo, oid, expOid)
		testFirstParentVariations(t, repo, oid, expOid)
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
		oid, expOid := OidNow(c.Oid), OidNow(output.Commits[i].Oid)
		testSecondAncestor(t, repo, oid, expOid)
		testSecondAncestorVariations(t, repo, oid, expOid)
	}
}

func Test_revParse__zeros(t *testing.T) {
	testCase := test.Linear
	repo := Open(testCase.Repo())
	output := testCase.Output().(*test.OutputCommits)

	util.Assert(t, output.N > 0)
	util.Assert(t, len(output.Commits) == output.N)

	// test the first, parentless commit
	for _, c := range output.Commits {
		oid := OidNow(c.Oid)
		testZeros(t, repo, oid)
	}
}

// testShortOid retrives the object by all possible combinations of
// shortening its id.
func testShortOid(t *testing.T, repo Repository, oid *ObjectId) {
	rev := oid.String()
	for i := 4; i <= 40; i++ {
		o, err := ObjectFromRevision(repo, rev[:i])
		util.AssertNoErr(t, err)
		util.AssertEqualString(t, o.ObjectId().String(), oid.String())
	}
}

func testFirstParent(t *testing.T, repo Repository, oid *ObjectId, parentOid *ObjectId) {
	rev := oid.String() // rev == oid here
	testObjectExpected(t, repo, rev, oid)
	for i := 4; i <= 40; i++ {
		testObjectExpected(t, repo, rev[:i]+"^", parentOid)
	}
}

func testFirstParentVariations(t *testing.T, repo Repository, oid *ObjectId, parentOid *ObjectId) {
	rev := oid.String()
	testObjectExpected(t, repo, rev+"^", parentOid)
	testObjectExpected(t, repo, rev+"^1", parentOid)
	testObjectExpected(t, repo, rev+"~", parentOid)
	testObjectExpected(t, repo, rev+"~1", parentOid)
}

func testSecondAncestor(t *testing.T, repo Repository, oid *ObjectId, parentOid *ObjectId) {
	rev := oid.String() // rev == oid here
	testObjectExpected(t, repo, rev, oid)
	for i := 4; i <= 40; i++ {
		testObjectExpected(t, repo, rev[:i]+"~2", parentOid)
	}
}

func testSecondAncestorVariations(t *testing.T, repo Repository, oid *ObjectId, parentOid *ObjectId) {
	rev := oid.String()
	testObjectExpected(t, repo, rev+"^^", parentOid)
	testObjectExpected(t, repo, rev+"^1^1", parentOid)
	testObjectExpected(t, repo, rev+"^^1", parentOid)
	testObjectExpected(t, repo, rev+"^1^", parentOid)
	testObjectExpected(t, repo, rev+"~~", parentOid)
	testObjectExpected(t, repo, rev+"~1~", parentOid)
	testObjectExpected(t, repo, rev+"~1~1", parentOid)
	testObjectExpected(t, repo, rev+"~~1", parentOid)
}

func testParentlessCommit(t *testing.T, repo Repository, oid *ObjectId) {
	rev := oid.String()
	testObjectExpected(t, repo, rev, oid)
	_, err := ObjectFromRevision(repo, rev+"~1")
	util.Assert(t, err != nil)
	_, err = ObjectFromRevision(repo, rev+"^")
	util.Assert(t, err != nil)
}

func testZeros(t *testing.T, repo Repository, oid *ObjectId) {
	rev := oid.String()
	testObjectExpected(t, repo, rev+"^0", oid)
	testObjectExpected(t, repo, rev+"~0", oid)
}

// ================================================================= //
// UTIL
// ================================================================= //

// testObjectExpected retrieves the commit with the given revision specification
// from the given repository and ensures that this operation went well and the
// returned object in fact has the expected oid.
func testObjectExpected(t *testing.T, repo Repository, rev string, expOid *ObjectId) {
	parent, err := ObjectFromRevision(repo, rev)
	util.AssertNoErr(t, err)
	util.Assert(t, parent.Header().Type() == ObjectCommit)
	util.AssertEqualString(t, parent.ObjectId().String(), expOid.String())
}
