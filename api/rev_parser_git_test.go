//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import (
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/test"
	"github.com/jbrukh/ggit/util"
	"testing"
)

func Test_revParse__firstParent(t *testing.T) {
	testCase := test.Linear
	repo := Open(testCase.Repo())
	info := testCase.Info().(*test.InfoLinear)

	util.Assert(t, info.N > 1)
	util.Assert(t, len(info.Commits) == info.N)

	// test the first, parentless commit
	testParentlessCommit(t, repo, objects.OidNow(info.Commits[0].CommitOid))
	for _, c := range info.Commits[1:] {
		oid, expOid := objects.OidNow(c.CommitOid), objects.OidNow(c.ParentOid)
		testShortOid(t, repo, oid)
		testFirstParent(t, repo, oid, expOid)
		testFirstParentVariations(t, repo, oid, expOid)
	}
}

func Test_revParse__secondAncestor(t *testing.T) {
	testCase := test.Linear
	repo := Open(testCase.Repo())
	info := testCase.Info().(*test.InfoLinear)

	util.Assert(t, info.N > 2)
	util.Assert(t, len(info.Commits) == info.N)

	// test the first, parentless commit
	for i, c := range info.Commits[2:] {
		oid, expOid := objects.OidNow(c.CommitOid), objects.OidNow(info.Commits[i].CommitOid)
		testSecondAncestor(t, repo, oid, expOid)
		testSecondAncestorVariations(t, repo, oid, expOid)
	}
}

func Test_revParse__zeros(t *testing.T) {
	testCase := test.Linear
	repo := Open(testCase.Repo())
	info := testCase.Info().(*test.InfoLinear)

	util.Assert(t, info.N > 0)
	util.Assert(t, len(info.Commits) == info.N)

	// test the first, parentless commit
	for _, c := range info.Commits {
		oid := objects.OidNow(c.CommitOid)
		testZeros(t, repo, oid)
	}
}

func Test_revParse__derefs(t *testing.T) {
	testCase := test.Derefs
	repo := Open(testCase.Repo())
	info := testCase.Info().(*test.InfoDerefs)

	commitOid := objects.OidNow(info.CommitOid)
	tagOid := objects.OidNow(info.TagOid)
	treeOid := objects.OidNow(info.TreeOid)
	testObjectExpected(t, repo, "HEAD", commitOid, objects.ObjectCommit)
	testObjectExpected(t, repo, "HEAD^{commit}", commitOid, objects.ObjectCommit)
	testObjectExpected(t, repo, "HEAD^{tree}", treeOid, objects.ObjectTree)
	testObjectExpected(t, repo, "HEAD^{commit}^{tree}", treeOid, objects.ObjectTree)
	testObjectExpected(t, repo, info.TagName, tagOid, objects.ObjectTag)
	testObjectExpected(t, repo, info.TagName+"^{commit}", commitOid, objects.ObjectCommit)
	testObjectExpected(t, repo, info.TagName+"^{commit}^{tree}", treeOid, objects.ObjectTree)

}

// testShortOid retrives the object by all possible combinations of
// shortening its id.
func testShortOid(t *testing.T, repo Repository, oid *objects.ObjectId) {
	rev := oid.String()
	for i := 4; i <= 40; i++ {
		o, err := ObjectFromRevision(repo, rev[:i])
		util.AssertNoErr(t, err)
		util.AssertEqualString(t, o.ObjectId().String(), oid.String())
	}
}

func testFirstParent(t *testing.T, repo Repository, oid *objects.ObjectId, parentOid *objects.ObjectId) {
	rev := oid.String() // rev == oid here
	testObjectExpected(t, repo, rev, oid, objects.ObjectCommit)
	for i := 4; i <= 40; i++ {
		testObjectExpected(t, repo, rev[:i]+"^", parentOid, objects.ObjectCommit)
	}
}

func testFirstParentVariations(t *testing.T, repo Repository, oid *objects.ObjectId, parentOid *objects.ObjectId) {
	rev := oid.String()
	testObjectExpected(t, repo, rev+"^", parentOid, objects.ObjectCommit)
	testObjectExpected(t, repo, rev+"^1", parentOid, objects.ObjectCommit)
	testObjectExpected(t, repo, rev+"~", parentOid, objects.ObjectCommit)
	testObjectExpected(t, repo, rev+"~1", parentOid, objects.ObjectCommit)
}

func testSecondAncestor(t *testing.T, repo Repository, oid *objects.ObjectId, parentOid *objects.ObjectId) {
	rev := oid.String() // rev == oid here
	testObjectExpected(t, repo, rev, oid, objects.ObjectCommit)
	for i := 4; i <= 40; i++ {
		testObjectExpected(t, repo, rev[:i]+"~2", parentOid, objects.ObjectCommit)
	}
}

func testSecondAncestorVariations(t *testing.T, repo Repository, oid *objects.ObjectId, parentOid *objects.ObjectId) {
	rev := oid.String()
	testObjectExpected(t, repo, rev+"^^", parentOid, objects.ObjectCommit)
	testObjectExpected(t, repo, rev+"^1^1", parentOid, objects.ObjectCommit)
	testObjectExpected(t, repo, rev+"^^1", parentOid, objects.ObjectCommit)
	testObjectExpected(t, repo, rev+"^1^", parentOid, objects.ObjectCommit)
	testObjectExpected(t, repo, rev+"~~", parentOid, objects.ObjectCommit)
	testObjectExpected(t, repo, rev+"~1~", parentOid, objects.ObjectCommit)
	testObjectExpected(t, repo, rev+"~1~1", parentOid, objects.ObjectCommit)
	testObjectExpected(t, repo, rev+"~~1", parentOid, objects.ObjectCommit)
}

func testParentlessCommit(t *testing.T, repo Repository, oid *objects.ObjectId) {
	rev := oid.String()
	testObjectExpected(t, repo, rev, oid, objects.ObjectCommit)
	_, err := ObjectFromRevision(repo, rev+"~1")
	util.Assert(t, err != nil)
	_, err = ObjectFromRevision(repo, rev+"^")
	util.Assert(t, err != nil)
}

func testZeros(t *testing.T, repo Repository, oid *objects.ObjectId) {
	rev := oid.String()
	testObjectExpected(t, repo, rev+"^0", oid, objects.ObjectCommit)
	testObjectExpected(t, repo, rev+"~0", oid, objects.ObjectCommit)
}

// ================================================================= //
// UTIL
// ================================================================= //

// testObjectExpected retrieves the commit with the given revision specification
// from the given repository and ensures that this operation went well and the
// returned object in fact has the expected oid.
func testObjectExpected(t *testing.T, repo Repository, rev string, expOid *objects.ObjectId, expType objects.ObjectType) {
	parent, err := ObjectFromRevision(repo, rev)
	util.AssertNoErr(t, err)
	util.Assert(t, parent.Header().Type() == expType)
	util.AssertEqualString(t, parent.ObjectId().String(), expOid.String())
}
