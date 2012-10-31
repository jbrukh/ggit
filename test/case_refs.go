//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
case_refs.go implements a repo test case, which contains a bunch of refs
which point to the one and only commit.
*/
package test

import (
	"fmt"
	"github.com/jbrukh/ggit/util"
)

// ================================================================= //
// TEST CASE: A BUNCH OF CONCRETE AND SYMBOLIC REFS
// ================================================================= //

type OutputRefs struct {
	CommitOid          string // the one and only commit
	BlobOid            string // the one and only blob
	AnnTagName         string // name of annotated tag
	AnnTagOid          string // the oid of the annotated tag
	LightTagName       string // name of lightweight tag
	BranchName         string // name of the branch that is the same as master
	SymbolicRef1       string // points to the branch
	SymbolicRef1Target string // the name branch it points to
	SymbolicRef2       string // points to the first symbolic ref
	SymbolicRef2Target string // the first symbolic ref
}

var Refs = NewRepoTestCase(
	"__refs",
	func(testCase *RepoTestCase) (err error) {
		err = createRepo(testCase)
		if err != nil {
			return err
		}

		repo := testCase.repo

		name := "myfile1.txt"
		contents := "various refs lol"
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

		// create a single commit
		err = util.GitExecMany(repo,
			[]string{"add", "--all"},
			[]string{"commit", "-a", "-m", "\"First and only commit\""},
		)
		if err != nil {
			return fmt.Errorf("could not commit to repo: %s", err)
		}

		// create an annotated tag
		annTagName := "annotated_tag"
		_, err = util.GitExec(repo, "tag", "-a", annTagName, "-m", "My tag!")
		if err != nil {
			return fmt.Errorf("could not create tag: %s", err)
		}

		// create a lightweight tag
		lightTagName := "lightweight_tag"
		_, err = util.GitExec(repo, "tag", lightTagName)
		if err != nil {
			return fmt.Errorf("could not create tag: %s", err)
		}

		// create a branch
		branchName := "regular_branch"
		_, err = util.GitExec(repo, "branch", branchName)
		if err != nil {
			return fmt.Errorf("could not create branch: %s", err)
		}

		// create a symbolic ref (1 deep)
		symbolicRef1 := "symbolic1"
		_, err = util.GitExec(repo, "symbolic-ref", symbolicRef1, branchName)
		if err != nil {
			return fmt.Errorf("could not create symbolic ref: %s", err)
		}

		// create a symbolic ref (2 deep)
		symbolicRef2 := "symbolic2"
		_, err = util.GitExec(repo, "symbolic-ref", symbolicRef2, symbolicRef1)
		if err != nil {
			return fmt.Errorf("could not create symbolic ref: %s", err)
		}

		// get the output data
		output := &OutputRefs{
			CommitOid:          util.RevOid(repo, "HEAD"),
			BlobOid:            blobOid,
			AnnTagName:         annTagName,
			AnnTagOid:          util.RevOid(repo, annTagName),
			LightTagName:       lightTagName,
			BranchName:         branchName,
			SymbolicRef1:       symbolicRef1,
			SymbolicRef1Target: branchName,
			SymbolicRef2:       symbolicRef2,
			SymbolicRef2Target: symbolicRef1,
		}

		testCase.output = output
		return
	},
)
