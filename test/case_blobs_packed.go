//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
case_blobs_packed.go implements a test repository similar to case_blobs.go,
but with all the objects packed.
*/
package test

import (
	"github.com/jbrukh/ggit/util"
)

// ================================================================= //
// TEST CASE: BUNCHES OF BLOBS
// ================================================================= //

type OutputBlobsPacked struct {
	OutputBlobs
}

var BlobsPacked = NewRepoTestCase(
	"__blobs_packed",
	func(testCase *RepoTestCase) (err error) {
		err = Blobs.builder(testCase)
		if err != nil {
			return err
		}

		// pack all that shit and remove loose shit
		err = util.GitExecMany(testCase.Repo(),
			[]string{"repack", "-a"},
			[]string{"prune-packed"},
		)

		testCase.output = &OutputBlobsPacked{
			*testCase.output.(*OutputBlobs),
		}

		return
	},
)
