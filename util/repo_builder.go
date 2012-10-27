//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
repo_builder.go implements methods for building up test repos.
*/
package util

import (
	"errors"
	"fmt"
	"os"
)

const (
	paramLinearN = 10
)

// ================================================================= //
// ALL TEST CASES
// ================================================================= //

// mapping of name => RepoTestCase
var RepoTestCases = make(map[string]*RepoTestCase, len(RepoBuilders))

func CreateRepoTestCases() (err error) {
	// build the test case
	var testCase *RepoTestCase
	for _, testCaseFunc := range RepoBuilders {
		testCase, err = testCaseFunc()
		if err != nil {
			return err
		}
		fmt.Printf("Created case: %s\n\n", testCase.Name())
		RepoTestCases[testCase.Name()] = testCase
	}
	return
}

func RemoveRepoTestCases() {
	for _, testCase := range RepoTestCases {
		if testCase != nil {
			testCase.Remove()
		}
	}
}

type RepoTestCase struct {
	name string
	repo string // path
}

func (tc *RepoTestCase) Repo() string {
	return tc.repo
}

func (tc *RepoTestCase) Name() string {
	return tc.name
}

func (tc *RepoTestCase) Remove() {
	os.RemoveAll(tc.repo)
}

func NewRepoTestCase(name string) (*RepoTestCase, error) {
	repo := TempRepo(name)
	_, err := CreateGitRepo(repo)
	if err != nil {
		return nil, fmt.Errorf("Could not create case '%s': %s", name, err.Error())
	}
	return &RepoTestCase{
		name: name,
		repo: repo,
	}, nil
}

type RepoBuilder func() (*RepoTestCase, error)

var RepoBuilders = []RepoBuilder{
	RepoBuilderEmpty,
	RepoBuilderLinear,
}

var RepoBuilderEmpty = func() (*RepoTestCase, error) {
	return NewRepoTestCase("__empty")
}

var RepoBuilderLinear = func() (testCase *RepoTestCase, err error) {
	n := paramLinearN
	testCase, err = NewRepoTestCase("__linear")
	if err != nil {
		return nil, err
	}
	if n < 1 {
		return nil, errors.New("n must be > 0")
	}
	repo := testCase.repo
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("%d.txt", i)
		err = TestFile(repo, name, string(i))
		if err != nil {
			return nil, errors.New("could not create test file for repo: " + err.Error())
		}
		// create a few commits
		err = GitExecMany(repo,
			[]string{"add", "--all"},
			[]string{"commit", "-a", "-m", fmt.Sprintf("Commit: %d", i)},
		)
		if err != nil {
			return nil, errors.New("could not commit to repo")
		}

	}
	return
}
