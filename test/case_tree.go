//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
case_tree.go implements a repo test case, which contains a single commit,
with a complex tree that covers all entry cases.
*/
package test

import (
	//"fmt"
	"github.com/jbrukh/ggit/util"
)

// ================================================================= //
// TEST CASE: A COMPLEX TREE
// ================================================================= //

type InfoTree struct {
	CommitOid string
	TreeOid   string
	File1Oid  string
	File2Oid  string
	File3Oid  string
	Tree1Oid  string
	Tree2Oid  string
}

var Tree = NewRepoTestCase(
	"__tree",
	func(testCase *RepoTestCase) error {
		repo, err := createRepo(testCase)
		if err != nil {
			return err
		}

		// add some files
		file1 := "one.txt"
		file2 := "two.txt"
		file3 := "three.txt"

		err = util.TestFile(repo, file1, "1")
		if err != nil {
			return err
		}

		err = util.TestFile(repo, file2, "2")
		if err != nil {
			return err
		}

		err = util.TestFile(repo, file3, "3")
		if err != nil {
			return err
		}

		// add some trees
		tree1 := "mytree/hello.txt"
		tree2 := "anothertree/bye.txt"

		err = util.TestFile(repo, tree1, "hello")
		if err != nil {
			return err
		}

		err = util.TestFile(repo, tree2, "bye")
		if err != nil {
			return err
		}

		// create a single commit
		err = util.GitExecMany(repo,
			[]string{"add", "--all"},
			[]string{"commit", "-a", "-m", "\"A complicated tree.\""},
		)
		if err != nil {
			return err
		}

		// get the output data
		info := &InfoTree{
			CommitOid: util.RevOid(repo, "HEAD"),
			TreeOid:   util.RevOid(repo, "HEAD^{tree}"),
			File1Oid:  util.BlobOid(repo, file1),
			File2Oid:  util.BlobOid(repo, file2),
			File3Oid:  util.BlobOid(repo, file3),
		}

		testCase.info = info
		return err
	},
)
