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
	"github.com/jbrukh/ggit/test"
	"github.com/jbrukh/ggit/util"
	"testing"
)

// this test implements basic blob reading from
// a git repository.
func Test_readBlobs(t *testing.T) {
	testRepo := test.Blobs

	repo := Open(testRepo.Repo())
	output := testRepo.Output().([]*test.OutputBlob)

	if len(output) < 1 {
		fmt.Println("warning: no blobs to test")
	}

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
