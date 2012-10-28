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

func Test_readBlobs(t *testing.T) {
	testRepo := test.Blobs

	// create a ggit repo
	repo := Open(testRepo.Repo())
	output := testRepo.Output().([]*test.OutputBlob)

	// hash the test objects
	for _, out := range output {

		// read the blob
		oid := OidNow(out.Oid)
		o, err := repo.ObjectFromOid(oid)
		util.AssertNoErr(t, err)
		util.Assert(t, o.Header().Type() == ObjectBlob)
		b := o.(*Blob)
		util.AssertEqualString(t, b.String(), out.Contents)
		util.AssertEqualInt(t, int(b.Header().Size()), len(out.Contents))
	}

}
