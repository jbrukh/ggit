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

type CommitDetail struct {
	CommitOid  string // oid of the given commit
	ParentOid  string // first parent oid
	CommitRepr string // representation of this commit as a string
	CommitSize int    // size of the commit object
	TreeOid    string
	BranchName string
	TagName    string
	TagSize    int
	TagOid     string
	TagRepr    string
}

type InfoLinear struct {
	Commits []*CommitDetail
	N       int
}

var Linear = NewRepoTestCase(
	"__linear",
	func(testCase *RepoTestCase) error {
		n := 10
		repo, err := createRepo(testCase)
		if err != nil {
			return err
		}
		if n < 1 {
			return errors.New("n must be > 0")
		}
		info := &InfoLinear{
			Commits: make([]*CommitDetail, n),
			N:       n,
		}
		for i := 0; i < n; i++ {
			name := fmt.Sprintf("%d.txt", i)
			err = util.TestFile(repo, name, string(i))
			if err != nil {
				return errors.New("could not create test file for repo: " + err.Error())
			}

			var (
				branchName = fmt.Sprintf("branch_%d", i)
				tagName    = fmt.Sprintf("tag_%d", i)
				commitMsg  = fmt.Sprintf("\"Commit: %d\"", i)
			)
			// create a commits
			err = util.GitExecMany(repo,
				[]string{"add", "--all"},
				[]string{"commit", "-a", "-m", commitMsg},
				[]string{"branch", branchName},
				[]string{"tag", "-a", tagName, "-m", commitMsg},
			)
			if err != nil {
				return fmt.Errorf("could not commit to repo: %s", err)
			}

			// get the output data
			var (
				parentOid = ""
				oid       = util.RevOid(repo, "HEAD")
				tagOid    = util.RevOid(repo, tagName)
			)
			if i != 0 {
				parentOid = util.RevOid(repo, "HEAD^")
			}
			info.Commits[i] = &CommitDetail{
				CommitOid:  oid,
				ParentOid:  parentOid,
				CommitRepr: util.ObjectRepr(repo, oid),
				CommitSize: util.ObjectSize(repo, oid),
				TreeOid:    util.RevOid(repo, oid+"^{tree}"),
				BranchName: branchName,
				TagName:    tagName,
				TagOid:     tagOid,
				TagRepr:    util.ObjectRepr(repo, tagName),
				TagSize:    util.ObjectSize(repo, tagOid),
			}
		}
		testCase.info = info
		return nil
	},
)
