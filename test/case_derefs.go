//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
case_derefs.go implements a repo test case, which contains a single commit,
a tag on that commit, and data about the underlying tag, tree, and commit oids.
*/
package test

import (
	"fmt"
	"github.com/jbrukh/ggit/util"
)

// ================================================================= //
// TEST CASE: A COMMIT AND TAG AND TREE AND BLOB FOR DEREFERENCING
// ================================================================= //

type InfoDerefs struct {
	TagName    string
	CommitOid  string
	TreeOid    string
	TagOid     string
	BranchName string
	BlobOid    string
}

var Derefs = NewRepoTestCase(
	"__derefs",
	func(testCase *RepoTestCase) (err error) {
		err = createRepo(testCase)
		if err != nil {
			return err
		}

		repo := testCase.repo

		name := "myfile1.txt"
		contents := "one"
		err = util.TestFile(repo, name, contents)
		if err != nil {
			return fmt.Errorf("could not create test file for repo: %s", err)
		}

		// hacky: figure out the blob oid of the file above
		var blobOid string
		blobOid, err = util.HashBlob(repo, contents)
		if err != nil {
			return fmt.Errorf("could not figure out blob oid: %s", err)
		}

		tagName := "0.0.0"
		branchName := "brooklyn"

		// create a single commit
		err = util.GitExecMany(repo,
			[]string{"add", "--all"},
			[]string{"commit", "-a", "-m", "\"First and only commit\""},
			[]string{"tag", "-a", tagName, "-m", "My tag!"},
			[]string{"branch", branchName},
		)
		if err != nil {
			return fmt.Errorf("could not commit to repo: %s", err)
		}

		// get the output data
		info := &InfoDerefs{
			TagName:    tagName,
			BranchName: branchName,
			CommitOid:  util.RevOid(repo, "HEAD"),
			TreeOid:    util.RevOid(repo, "HEAD^{tree}"),
			TagOid:     util.RevOid(repo, tagName),
			BlobOid:    blobOid,
		}

		testCase.info = info
		return
	},
)
