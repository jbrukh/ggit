//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
refs_git_test.go implements tests for reading conrete and symbolic refs.
*/
package api

import (
	"github.com/jbrukh/ggit/test"
	"github.com/jbrukh/ggit/util"
	"testing"
)

func Test_refPaths(t *testing.T) {
	testRepo := test.Refs
	repo := Open(testRepo.Repo())
	output := testRepo.Output().(*test.OutputRefs)

	oid := OidNow(output.CommitOid)
	tagOid := OidNow(output.AnnTagOid)
	var (
		master   = "refs/heads/master"
		branch   = "refs/heads/" + output.BranchName
		annTag   = "refs/tags/" + output.AnnTagName
		lightTag = "refs/tags/" + output.LightTagName
	)

	// test reading these full path refs directly from
	// the repository files, loose or packed
	testRefPathPeeled(t, repo, master, oid)
	testRefPathPeeled(t, repo, branch, oid)
	testRefPathPeeled(t, repo, annTag, tagOid)
	testRefPathPeeled(t, repo, lightTag, oid)

	testRefPathSymbolic(t, repo, output.SymbolicRef1, output.SymbolicRef1Target)
	testRefPathSymbolic(t, repo, output.SymbolicRef2, output.SymbolicRef2Target)
}

func testRefPathPeeled(t *testing.T, repo Repository, spec string, oid *ObjectId) {
	ref, err := repo.Ref(spec)
	util.AssertNoErr(t, err)
	util.AssertEqualString(t, ref.Name(), spec)

	// make sure target is symbolic and matches
	symbolic, target := ref.Target()
	util.Assert(t, !symbolic)
	util.AssertEqualString(t, target.(*ObjectId).String(), oid.String())
	util.AssertEqualString(t, ref.ObjectId().String(), oid.String())
	util.AssertPanic(t, func() {
		s := target.(string)
		s += ""
	})
}

func testRefPathSymbolic(t *testing.T, repo Repository, spec string, tget string) {
	ref, err := repo.Ref(spec)
	util.AssertNoErr(t, err)
	util.AssertEqualString(t, ref.Name(), spec)
	// make sure target is symbolic and matches
	symbolic, target := ref.Target()
	util.Assert(t, symbolic)
	util.AssertEqualString(t, target.(string), tget)
	util.AssertPanic(t, func() {
		oid := target.(*ObjectId)
		oid.String()
	})
	util.AssertPanic(t, func() {
		oid := ref.ObjectId()
		oid.String()
	})

}
