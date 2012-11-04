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

func Test_Open(t *testing.T) {
	util.Assert(t, Open("test").path == "test/.git")
	util.Assert(t, Open("test/.git").path == "test/.git")
}

func Test_AssertDiskRepo(t *testing.T) {
	ggit := Open("./..")
	diskRepo, err := AssertDiskRepo(ggit)
	util.AssertNoErr(t, err)
	util.Assert(t, diskRepo != nil)
}

func Test_LooseObjectIds(t *testing.T) {
	testCase := test.Derefs
	repo := Open(testCase.Repo())
	info := testCase.Info().(*test.InfoDerefs)

	oids, err := repo.LooseObjectIds()
	util.AssertNoErr(t, err)

	util.AssertEqualInt(t, info.ObjectsN, len(oids))
	util.Assert(t, compareLists(oids, []*ObjectId{
		OidNow(info.BlobOid),
		OidNow(info.CommitOid),
		OidNow(info.TagOid),
		OidNow(info.TreeOid),
	}),
	)
}

func compareLists(one []*ObjectId, two []*ObjectId) bool {
	if len(one) != len(two) {
		return false
	}
	for _, v := range one {
		found := false
		for _, w := range two {
			if v.String() == w.String() {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func Test_PackedObjects(t *testing.T) {
	testCase := test.DerefsPacked
	repo := Open(testCase.Repo())
	info := testCase.Info().(*test.InfoDerefsPacked)

	objects, err := repo.PackedObjects()
	util.AssertNoErr(t, err)

	var oids []*ObjectId
	for _, o := range objects {
		oids = append(oids, o.Object().ObjectId())
	}

	util.AssertEqualInt(t, info.ObjectsN, len(oids))
	util.Assert(t, compareLists(oids, []*ObjectId{
		OidNow(info.BlobOid),
		OidNow(info.CommitOid),
		OidNow(info.TagOid),
		OidNow(info.TreeOid),
	}),
	)
}

func Test_PackedObjectIds(t *testing.T) {
	testCase := test.DerefsPacked
	repo := Open(testCase.Repo())
	info := testCase.Info().(*test.InfoDerefsPacked)

	oids, err := repo.PackedObjectIds()
	util.AssertNoErr(t, err)

	util.AssertEqualInt(t, info.ObjectsN, len(oids))
	util.Assert(t, compareLists(oids, []*ObjectId{
		OidNow(info.BlobOid),
		OidNow(info.CommitOid),
		OidNow(info.TagOid),
		OidNow(info.TreeOid),
	}),
	)
}

func Test_Refs(t *testing.T) {
	testCase := test.Derefs
	repo := Open(testCase.Repo())
	info := testCase.Info().(*test.InfoDerefs)
	refs, err := repo.Refs()
	util.AssertNoErr(t, err)

	util.AssertEqualInt(t, info.RefsN, len(refs))
}
