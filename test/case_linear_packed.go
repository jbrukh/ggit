//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
case_linear_packed.go implements a test repository similar to case_linear.go,
but with all the objects packed.
*/
package test

// ================================================================= //
// TEST CASE: BUNCHES OF BLOBS
// ================================================================= //

type OutputLinearPacked struct {
	OutputLinear
}

var LinearPacked = NewRepoTestCase(
	"__linear_packed",
	func(testCase *RepoTestCase) (err error) {
		err = Linear.builder(testCase)
		if err != nil {
			return err
		}

		// pack all that shit and remove loose shit
		err = GitExecMany(testCase.Repo(),
			[]string{"repack", "-a"},
			[]string{"prune-packed"},
		)

		testCase.output = &OutputLinearPacked{
			*testCase.output.(*OutputLinear),
		}

		return
	},
)
