//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
tags_git_test.go implements git comparison tests for tag reading.
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
func Test_readTags(t *testing.T) {
	testCase := test.Linear
	repo := Open(testCase.Repo())
	info := testCase.Info().(*test.InfoLinear)

	util.Assert(t, info.N > 1)
	util.Assert(t, len(info.Commits) == info.N)

	f := NewStrFormat()
	for _, detail := range info.Commits {
		tagOid := objects.OidNow(detail.TagOid)
		o, err := repo.ObjectFromOid(tagOid)
		util.AssertNoErr(t, err)

		// check the id
		util.Assert(t, o.ObjectId().String() == detail.TagOid)

		// check the header
		util.Assert(t, o.Header().Type() == ObjectTag)
		util.AssertEqualInt(t, int(o.Header().Size()), detail.TagSize)

		// now convert to a tag and check the fields
		var tag *Tag
		util.AssertPanicFree(t, func() {
			tag = o.(*Tag)
		})

		// check the name
		util.AssertEqualString(t, tag.Name(), detail.TagName)

		// check the target object
		util.Assert(t, tag.Object() != nil)
		util.AssertEqualString(t, tag.Object().String(), detail.CommitOid)
		util.Assert(t, tag.ObjectType() == ObjectCommit)

		// check the whole representation, which will catch
		// most of the other stuff
		f.Reset()
		f.Object(o)
		util.AssertEqualString(t, detail.TagRepr, f.String())
	}
}
