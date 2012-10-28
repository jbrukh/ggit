//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
test_case.go implements methods for building up test repos.
*/
package test

import (
	"os"
)

const (
	paramLinearN = 10
)

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
