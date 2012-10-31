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

func Test_refs(t *testing.T) {
	testRepo := test.Refs

	repo := Open(testRepo.Repo())
	output := testRepo.Output().(*test.OutputRefs)

	testSymbolicRef(t, repo, "HEAD", output.HeadTarget)
	testSymbolicRef(t, repo, output.SymbolicRef1, output.SymbolicRef1Target)
	testSymbolicRef(t, repo, output.SymbolicRef2, output.SymbolicRef2Target)

	oid := OidNow(output.CommitOid)
	testConcreteRef(t, repo, "refs/heads/master", oid)

}

func testSymbolicRef(t *testing.T, repo Repository, spec string, expTarget string) {
	ref, err := repo.Ref(spec)
	util.AssertNoErr(t, err)

	// assert the name sticks
	util.AssertEqualString(t, ref.Name(), spec)

	// make sure target is symbolic and matches
	symbolic, target := ref.Target()
	util.Assert(t, symbolic)
	util.AssertEqualString(t, target.(string), expTarget)
}

func testConcreteRef(t *testing.T, repo Repository, spec string, oid *ObjectId) {
	ref, err := repo.Ref(spec)
	util.AssertNoErr(t, err)

	// assert the name sticks
	util.AssertEqualString(t, ref.Name(), spec)

	// make sure target is symbolic and matches
	symbolic, target := ref.Target()
	util.Assert(t, !symbolic)
	util.AssertEqualString(t, target.(*ObjectId).String(), oid.String())
	util.AssertEqualString(t, ref.ObjectId().String(), oid.String())
}
