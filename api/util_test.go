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

func Test_makeHash(t *testing.T) {
	testRepo := test.DerefsPacked
	repo := Open(testRepo.Repo())
	info := testRepo.Info().(*test.InfoDerefsPacked)
	tagOid := OidNow(info.TagOid)
	o, err := repo.ObjectFromOid(tagOid)
	if err != nil {
		panic(err)
	}
	if h, err := MakeHash(o); err != nil {
		panic(err)
	} else {
		oid := OidFromHash(h)
		util.Assert(t, oid.String() == tagOid.String())
	}
}
