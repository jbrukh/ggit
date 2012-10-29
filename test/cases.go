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
	"time"
)

// ================================================================= //
// ALL TEST CASES
// ================================================================= //

// mapping of name => RepoTestCase
var repoTestCases = []*RepoTestCase{
	Empty,
	Linear,
	LinearPacked,
	Blobs,
	Derefs,
}

// init initializes all the repo test cases, if they haven't been
// initialized already. An unsuccessful initialization will cause
// the entire process to exit.
func init() {
	fmt.Println("Creating repo test cases...\n")
	for _, testCase := range repoTestCases {
		start := time.Now()
		err := testCase.Build()
		if err != nil {
			fmt.Printf("error (exiting!): %s\n", err)
			RemoveTestCases()
			os.Exit(1)
		}
		fmt.Printf("Created case: %s (%d ms)\n\n", testCase.Name(), int64(time.Since(start))/int64(time.Millisecond))
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

// ================================================================= //
// REPO TEST CASE
// ================================================================= //

type RepoTestCase struct {
	name    string
	repo    string // path
	builder RepoBuilder
	output  interface{}
}

func (tc *RepoTestCase) Repo() string {
	return tc.repo
}

func (tc *RepoTestCase) Name() string {
	return tc.name
}

func (tc *RepoTestCase) Output() interface{} {
	return tc.output
}

func (tc *RepoTestCase) Remove() {
	os.RemoveAll(tc.repo)
}

func (tc *RepoTestCase) Build() error {
	return tc.builder(tc)
}

func NewRepoTestCase(name string, builder RepoBuilder) *RepoTestCase {
	return &RepoTestCase{
		name:    name,
		builder: builder,
	}
}

// ================================================================= //
// REPO BUILDER
// ================================================================= //

type RepoBuilder func(testCase *RepoTestCase) error

// ================================================================= //
// UTIL
// ================================================================= //

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
