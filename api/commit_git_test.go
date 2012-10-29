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
	output := testCase.Output().(*test.OutputCommits)

	util.Assert(t, output.N > 1)
	util.Assert(t, len(output.Commits) == output.N)

	f := NewStrFormat()
	for _, c := range output.Commits {
		o, err := repo.ObjectFromOid(OidNow(c.Oid))
		util.AssertNoErr(t, err)
		util.Assert(t, o.ObjectId().String() == c.Oid)
		util.Assert(t, o.Header().Type() == ObjectCommit)
		f.Reset()
		f.Object(o)
		util.AssertEqualString(t, c.Repr, f.String())
	}
}
