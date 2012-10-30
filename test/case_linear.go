//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
case_linear.go implements a repo test case.
*/
package test

import (
	"errors"
	"fmt"
	"github.com/jbrukh/ggit/util"
)

// ================================================================= //
// TEST CASE: A FEW LINEAR COMMITS
// ================================================================= //

type CommitInfo struct {
	Oid        string
	ParentOid  string // first parent oid
	Repr       string // representation of this commit as a string
	BranchName string
}

type OutputLinear struct {
	Commits []*CommitInfo
	N       int
}

var Linear = NewRepoTestCase(
	"__linear",
	func(testCase *RepoTestCase) (err error) {
		n := 10
		err = createRepo(testCase)
		if err != nil {
			return err
		}
		if n < 1 {
			return errors.New("n must be > 0")
		}
		repo := testCase.repo
		output := &OutputLinear{
			Commits: make([]*CommitInfo, n),
			N:       n,
		}
		for i := 0; i < n; i++ {
			name := fmt.Sprintf("%d.txt", i)
			err = util.TestFile(repo, name, string(i))
			if err != nil {
				return errors.New("could not create test file for repo: " + err.Error())
			}

			// create a commits
			err = util.GitExecMany(repo,
				[]string{"add", "--all"},
				[]string{"commit", "-a", "-m", fmt.Sprintf("\"Commit: %d\"", i)},
			)
			if err != nil {
				return fmt.Errorf("could not commit to repo: %s", err)
			}

			// create a branch for that commit
			branchName := fmt.Sprintf("branch_%d", i)
			_, err = util.GitExec(repo, "branch", branchName)
			if err != nil {
				return fmt.Errorf("could not create branch: %s", err)
			}

			// get the output data
			var oid, parentOid, repr string
			oid = util.RevOid(repo, "HEAD")
			repr, err = util.GitExec(repo, "cat-file", "-p", oid)
			if err != nil {
				return err
			}
			if i != 0 {
				parentOid = util.RevOid(repo, "HEAD^")
			}
			output.Commits[i] = &CommitInfo{
				oid,
				parentOid,
				repr,
				branchName,
			}
		}
		testCase.output = output
		return
	},
)
