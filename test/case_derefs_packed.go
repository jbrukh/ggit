//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
blobs.go implements a test repository.
*/
package test

import (
	"github.com/jbrukh/ggit/util"
)

// ================================================================= //
// TEST CASE: BUNCHES OF BLOBS
// ================================================================= //

type OutputDerefsPacked struct {
	OutputDerefs
}

var DerefsPacked = NewRepoTestCase(
	"__linear_packed",
	func(testCase *RepoTestCase) (err error) {
		err = Derefs.builder(testCase)
		if err != nil {
			return err
		}

		// pack all that shit and remove loose shit
		err = util.GitExecMany(testCase.Repo(),
			[]string{"repack", "-a"},
			[]string{"prune-packed"},
		)

		testCase.output = &OutputDerefsPacked{
			*testCase.output.(*OutputDerefs),
		}

		return
	},
)
