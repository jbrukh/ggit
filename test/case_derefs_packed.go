//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
case_derefs_packed.go implements a test repository similar to case_derefs.go,
but with all the objects packed.
*/
package test

import (
	"github.com/jbrukh/ggit/util"
)

// ================================================================= //
// TEST CASE: DEREFS
// ================================================================= //

type InfoDerefsPacked struct {
	InfoDerefs
}

var DerefsPacked = NewRepoTestCase(
	"__derefs_packed",
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

		testCase.info = &InfoDerefsPacked{
			*testCase.info.(*InfoDerefs),
		}

		return err
	},
)
