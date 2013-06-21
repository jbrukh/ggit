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

type InfoTreeDiff struct {
	//unchanged files
	File1Oid, File2Oid, File3Oid string
	//changed files
	ModifiedFileName, RemovedFileName, AddedFileName                          string
	ModifiedFileBeforeOid, ModifiedFileAfterOid, RemovedFileOid, AddedFileOid string
	CommitOid1, TreeOid1                                                      string
	CommitOid2, TreeOid2                                                      string
}

// ================================================================= //
// TEST CASE: A COMPLEX TREE
// ================================================================= //

var TreeDiff = NewRepoTestCase(
	"__tree_diff",
	func(testCase *RepoTestCase) error {
		repo, err := createRepo(testCase)
		if err != nil {
			return err
		}

		// add some files
		file1 := "1.txt"
		file2 := "2.txt"
		file3 := "3.txt"

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
		subtree1 := "4"
		subtree2 := "5"

		deepFile1 := path.Join(subtree1, "haha.txt")
		err = util.TestFile(repo, deepFile1, "hello, good sir!")
		if err != nil {
			return err
		}

		deepFile2 := path.Join(subtree2, "a_new_file.txt")
		err = util.TestFile(repo, deepFile2, "nothing to see here")
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
		commit1 := util.RevOid(repo, "HEAD")
		tree1 := util.RevOid(repo, "HEAD^{tree}")
		fileOid1 := util.TreeEntryOid(repo, file1)
		fileOid2 := util.TreeEntryOid(repo, file2)
		fileOid3 := util.TreeEntryOid(repo, file3)
		deepFileOid1a := util.TreeEntryOid(repo, deepFile1)
		deepFileOid2 := util.TreeEntryOid(repo, deepFile2)

		err = util.TestFile(repo, deepFile1, "this file been done changed")
		if err != nil {
			return err
		}

		if err = util.DeleteFile(repo, deepFile2); err != nil {
			return err
		}

		deepFile3 := path.Join(subtree2, "a_newer_file.txt")
		if err = util.TestFile(repo, deepFile3, "time and chance happeneth"); err != nil {
			return err
		}

		if err = util.GitExecMany(repo,
			[]string{"add", "--all"},
			[]string{"commit", "-a", "-m", "\"A modified tree.\""},
		); err != nil {
			return err
		}

		deepFileOid1b := util.TreeEntryOid(repo, deepFile1)

		commit2 := util.RevOid(repo, "HEAD")
		tree2 := util.RevOid(repo, "HEAD^{tree}")
		deepFileOid3 := util.TreeEntryOid(repo, deepFile3)

		testCase.info = &InfoTreeDiff{
			File1Oid:              fileOid1,
			File2Oid:              fileOid2,
			File3Oid:              fileOid3,
			ModifiedFileBeforeOid: deepFileOid1a,
			ModifiedFileAfterOid:  deepFileOid1b,
			RemovedFileOid:        deepFileOid2,
			AddedFileOid:          deepFileOid3,
			ModifiedFileName:      deepFile1,
			RemovedFileName:       deepFile2,
			AddedFileName:         deepFile3,
			CommitOid1:            commit1,
			CommitOid2:            commit2,
			TreeOid1:              tree1,
			TreeOid2:              tree2,
		}
		return err
	},
)
