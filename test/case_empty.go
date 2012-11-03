//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
case_empty.go implements a repo test case.
*/
package test

// ================================================================= //
// TEST CASE: EMPTY REPO
// ================================================================= //

var Empty = NewRepoTestCase(
	"__empty",
	func(testCase *RepoTestCase) error {
		_, err := createRepo(testCase)
		return err
	},
)
