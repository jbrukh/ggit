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
	"path"
)

// ================================================================= //
// TEST CASE: A COMPLEX TREE
// ================================================================= //

type InfoTree struct {
	CommitOid   string
	TreeOid     string
	TreeSize    int
	TreeRepr    string
	File1Oid    string
	File2Oid    string
	File3Oid    string
	Subtree1Oid string
	Subtree2Oid string
	N           int
}

var Tree = NewRepoTestCase(
	"__tree",
	func(testCase *RepoTestCase) error {
		repo, err := createRepo(testCase)
		if err != nil {
			return err
		}
		var n int

		// add some files
		file1 := "1.txt"
		file2 := "2.txt"
		file3 := "3.txt"

		err = util.TestFile(repo, file1, "1")
		if err != nil {
			return err
		}
		n++

		err = util.TestFile(repo, file2, "2")
		if err != nil {
			return err
		}
		n++

		err = util.TestFile(repo, file3, "3")
		if err != nil {
			return err
		}
		n++

		// add some trees
		subtree1 := "4"
		subtree2 := "5"

		err = util.TestFile(repo, path.Join(subtree1, "haha.txt"), "hello")
		if err != nil {
			return err
		}
		n++

		err = util.TestFile(repo, path.Join(subtree2, "hehe.txt"), "bye")
		if err != nil {
			return err
		}
		n++

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
			CommitOid:   util.RevOid(repo, "HEAD"),
			TreeOid:     util.RevOid(repo, "HEAD^{tree}"),
			TreeSize:    util.ObjectSize(repo, "HEAD^{tree}"),
			TreeRepr:    util.ObjectRepr(repo, "HEAD^{tree}"),
			File1Oid:    util.TreeEntryOid(repo, file1),
			File2Oid:    util.TreeEntryOid(repo, file2),
			File3Oid:    util.TreeEntryOid(repo, file3),
			Subtree1Oid: util.TreeEntryOid(repo, subtree1),
			Subtree2Oid: util.TreeEntryOid(repo, subtree2),
			N:           n,
		}
		testCase.info = info
		return err
	},
)
