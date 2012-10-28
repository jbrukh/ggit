//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
cases.go
*/
package test

import (
	"fmt"
	"github.com/jbrukh/ggit/util"
	"os"
)

// ================================================================= //
// ALL TEST CASES
// ================================================================= //

// mapping of name => RepoTestCase
var repoTestCases = []*RepoTestCase{
	Empty,
	Linear,
	Blobs,
}

// init initializes all the repo test cases, if they haven't been
// initialized already. An unsuccessful initialization will cause
// the entire process to exit.
func init() {
	fmt.Println("Creating repo test cases...\n")
	for _, testCase := range repoTestCases {
		err := testCase.Build()
		if err != nil {
			fmt.Printf("error (exiting!): %s\n", err)
			RemoveTestCases()
			os.Exit(1)
		}
		fmt.Printf("Created case: %s\n\n", testCase.Name())
	}
	fmt.Println("Done.\n")
}

func RemoveTestCases() {
	fmt.Println("Cleaning.")
	for _, testCase := range repoTestCases {
		fmt.Println(testCase.Name(), "\t", testCase.Repo())
		testCase.Remove()
	}
}

func createRepo(testCase *RepoTestCase) (err error) {
	repo := util.TempRepo(testCase.name)

	// clean that shit
	os.RemoveAll(repo)
	_, err = util.CreateGitRepo(repo)
	if err != nil {
		return fmt.Errorf("Could not create case '%s': %s", testCase.name, err.Error())
	}
	testCase.repo = repo
	return
}
