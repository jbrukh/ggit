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

	var (
		oid    = OidNow(output.CommitOid)
		tagOid = OidNow(output.AnnTagOid)

		master   = "refs/heads/master"
		branch   = expandHeadRef(output.BranchName)
		annTag   = expandTagRef(output.AnnTagName)
		lightTag = expandTagRef(output.LightTagName)
	)

	// test reading these full path refs directly from
	// the repository files, loose or packed
	testRefPathPeeled(t, repo, master, oid)
	testRefPathPeeled(t, repo, branch, oid)
	testRefPathPeeled(t, repo, annTag, tagOid)
	testRefPathPeeled(t, repo, lightTag, oid)

	// test reading symbolic refs and asserting that the
	// targets are in fact symbols and are correct
	testRefPathSymbolic(t, repo, output.SymbolicRef1, output.SymbolicRef1Target)
	testRefPathSymbolic(t, repo, output.SymbolicRef2, output.SymbolicRef2Target)
	testRefPathSymbolic(t, repo, "HEAD", master)

	// test that packed refs have correct commit
	// dereferencing information stored in the packed-refs file
	testPackedTagDerefInfo(t, repo, annTag, oid)

	// test ref peeling
	testPeelRef(t, repo, master, oid)
	testPeelRef(t, repo, branch, oid)
	testPeelRef(t, repo, output.SymbolicRef1, oid)
	testPeelRef(t, repo, output.SymbolicRef2, oid)

	// make sure we read loose refs correctly
	testRefRetrieval(t, repo, func() ([]Ref, error) {
		return repo.LooseRefs()
	}, []string{master, branch})

	// make sure we read packed refs correctly
	testRefRetrieval(t, repo, func() ([]Ref, error) {
		return repo.PackedRefs()
	}, []string{annTag, lightTag})

	// make sure we get all refs correctly
	testRefRetrieval(t, repo, func() ([]Ref, error) {
		return repo.Refs()
	}, []string{master, branch, annTag, lightTag})

	// test non existent refs
	testNonexistent(t, repo)

	// test resolution of short refs
	testShortRefResolvesPeeled(t, repo, "master", oid)
	testShortRefResolvesPeeled(t, repo, output.BranchName, oid)
	testShortRefResolvesPeeled(t, repo, output.AnnTagName, tagOid)
	testShortRefResolvesPeeled(t, repo, output.LightTagName, oid)
	testShortRefResolvesSymbolic(t, repo, "HEAD", master, oid)
	testShortRefResolvesSymbolic(t, repo, output.SymbolicRef1, output.SymbolicRef1Target, oid)
	testShortRefResolvesSymbolic(t, repo, output.SymbolicRef2, output.SymbolicRef2Target, oid)

}

func testShortRefResolvesSymbolic(t *testing.T, repo Repository, spec string, tget string, oid *ObjectId) {
	symbolicRef, err := RefFromSpec(repo, spec)
	util.AssertNoErr(t, err)
	assertSymbolicRef(t, symbolicRef, tget)

	peeledRef, err := PeeledRefFromSpec(repo, spec)
	util.AssertNoErr(t, err)
	assertPeeledRef(t, peeledRef, oid)
}

func testShortRefResolvesPeeled(t *testing.T, repo Repository, spec string, oid *ObjectId) {
	peeledRef, err := RefFromSpec(repo, spec)
	util.AssertNoErr(t, err)
	assertPeeledRef(t, peeledRef, oid)

	peeledRef, err = PeeledRefFromSpec(repo, spec)
	util.AssertNoErr(t, err)
	assertPeeledRef(t, peeledRef, oid)
}

func testNonexistent(t *testing.T, repo Repository) {
	cases := []string{
		"refs/heads/elephant",
		"refs/tags/tiger",
		"woodchuck",
		"goldmine",
	}
	for _, spec := range cases {
		assertBadRef(t, repo, spec)
	}
}

func assertBadRef(t *testing.T, repo Repository, spec string) {
	_, err := repo.Ref(spec)
	util.Assert(t, err != nil)
	util.Assert(t, IsNoSuchRef(err))
	_, err = RefFromSpec(repo, spec)
	util.Assert(t, err != nil)
	util.Assert(t, IsNoSuchRef(err))
	_, err = PeeledRefFromSpec(repo, spec)
	util.Assert(t, err != nil)
	util.Assert(t, IsNoSuchRef(err))
}

func testRefRetrieval(t *testing.T, repo Repository, f func() ([]Ref, error), expected []string) {
	refs, err := f()
	util.AssertNoErr(t, err)
	util.AssertEqualInt(t, len(expected), len(refs))
	for i, r := range expected {
		util.AssertEqualString(t, r, refs[i].Name())
	}
}

func testRefPathPeeled(t *testing.T, repo Repository, spec string, oid *ObjectId) {
	ref, err := repo.Ref(spec)
	util.AssertNoErr(t, err)
	util.AssertEqualString(t, ref.Name(), spec)
	assertPeeledRef(t, ref, oid)
}

func testRefPathSymbolic(t *testing.T, repo Repository, spec string, tget string) {
	symbolicRef, err := repo.Ref(spec)
	util.AssertNoErr(t, err)
	util.AssertEqualString(t, symbolicRef.Name(), spec)
	// make sure target is symbolic and matches
	assertSymbolicRef(t, symbolicRef, tget)
}

func testPackedTagDerefInfo(t *testing.T, repo Repository, spec string, oid *ObjectId) {
	ref, err := repo.Ref(spec)
	util.AssertNoErr(t, err)
	util.AssertEqualString(t, ref.Name(), spec)
	// make sure target is symbolic and matches
	symbolic, target := ref.Target()
	util.Assert(t, !symbolic)
	if target == nil || ref.Commit() == nil {
		t.Fatalf("nil target or commit")
	}
	util.AssertEqualString(t, ref.Commit().String(), oid.String())
}

func testPeelRef(t *testing.T, repo Repository, spec string, oid *ObjectId) {
	// first peel it manually
	ref, err := repo.Ref(spec)
	util.AssertNoErr(t, err)
	util.AssertEqualString(t, ref.Name(), spec)
	peeledRef, err := PeelRef(repo, ref)
	util.AssertNoErr(t, err)
	assertPeeledRef(t, peeledRef, oid)

	// now, peel it automatically
	ref, err = PeeledRefFromSpec(repo, spec)
	util.AssertNoErr(t, err)
	peeledRef, err = PeelRef(repo, ref)
	util.AssertNoErr(t, err)
	assertPeeledRef(t, peeledRef, oid)
}

func assertPeeledRef(t *testing.T, peeledRef Ref, oid *ObjectId) {
	symbolic, target := peeledRef.Target()
	util.Assert(t, !symbolic)
	if target == nil {
		t.Fatalf("nil target")
	}
	util.AssertEqualString(t, target.(*ObjectId).String(), oid.String())
	util.AssertEqualString(t, peeledRef.ObjectId().String(), oid.String())
	util.AssertPanic(t, func() {
		s := target.(string)
		s += "" // for compilation
	})
}

func assertSymbolicRef(t *testing.T, symbolicRef Ref, tget string) {
	symbolic, target := symbolicRef.Target()
	util.Assert(t, symbolic)
	if target == nil {
		t.Fatalf("nil target")
	}
	util.AssertEqualString(t, target.(string), tget)
	util.AssertPanic(t, func() {
		oid := target.(*ObjectId)
		oid.String() // for compilation
	})
	util.AssertPanic(t, func() {
		oid := symbolicRef.ObjectId()
		oid.String() // for compilation
	})
}
