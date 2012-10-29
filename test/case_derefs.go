//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
linear.go implements a repo test case.
*/
package test

import (
	"errors"
	"fmt"
	"github.com/jbrukh/ggit/util"
)

// ================================================================= //
// TEST CASE: A COMMIT AND TAG AND TREE AND BLOB FOR DEREFERENCING
// ================================================================= //

type OutputDerefs struct {
	TagName   string
	CommitOid string
	TreeOid   string
	TagOid    string
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
		err = util.TestFile(repo, name, "one")
		if err != nil {
			return errors.New("could not create test file for repo: " + err.Error())
		}

		// create a single commit
		err = util.GitExecMany(repo,
			[]string{"add", "--all"},
			[]string{"commit", "-a", "-m", "\"First and only commit\""},
		)
		if err != nil {
			return fmt.Errorf("could not commit to repo: %s", err)
		}

		// create a single tag
		tagName := "0.0.0"
		_, err = util.GitExec(repo, "tag", "-a", tagName, "-m", "My tag!")
		if err != nil {
			return fmt.Errorf("could not create tag: %s", err)
		}

		// get the output data
		output := new(OutputDerefs)
		output.TagName = tagName

		output.CommitOid = util.RevOid(repo, "HEAD")
		output.TreeOid = util.RevOid(repo, "HEAD^{tree}")
		output.TagOid = util.RevOid(repo, tagName)

		testCase.output = output
		return
	},
)
