//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import (
	"fmt"
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/test"
	"github.com/jbrukh/ggit/util"
	"testing"
)

// this test implements basic blob reading from
// a git repository.
func Test_readBlobs(t *testing.T) {
	testRepo := test.Blobs

	repo := Open(testRepo.Repo())
	info := testRepo.Info().(*test.InfoBlobs)

	if len(info.Blobs) < 1 {
		fmt.Println("warning: no blobs to test")
	}

	for _, detail := range info.Blobs {
		// read the blob
		oid := objects.OidNow(detail.Oid)
		o, err := repo.ObjectFromOid(oid)
		util.AssertNoErr(t, err)
		util.Assert(t, o.Header().Type() == objects.ObjectBlob)
		b := o.(*Blob)
		util.AssertEqualString(t, string(b.Data()), detail.Contents)
		util.AssertEqualInt(t, int(b.Header().Size()), len(detail.Contents))
	}

}
