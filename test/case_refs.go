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

type InfoRefs struct {
	CommitOid          string // the one and only commit
	AnnTagName         string // name of annotated tag
	AnnTagOid          string // the oid of the annotated tag
	LightTagName       string // name of lightweight tag
	BranchName         string // name of the branch that is the same as master
	SymbolicRef1       string // points to the branch
	SymbolicRef1Target string // the name branch it points to
	SymbolicRef2       string // points to the first symbolic ref
	SymbolicRef2Target string // the first symbolic ref
	HeadTarget         string // where HEAD points (symbolically)
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

		var (
			annTagName   = "annotated_tag"
			lightTagName = "lightweight_tag"
			branchName   = "regular_branch"
			symbolicRef1 = "symbolic1"
			symbolicRef2 = "symbolic2"
		)

		// create a single commit
		err = util.GitExecMany(repo,
			// add a commit, the one and only
			[]string{"add", "--all"},
			[]string{"commit", "-a", "-m", "\"First and only commit\""},

			// add an annotated and lightweight tag
			[]string{"tag", "-a", annTagName, "-m", "\"My tag!\""},
			[]string{"tag", lightTagName},

			// add a branch pointing to master
			[]string{"branch", branchName},

			// add a symbolic ref, 1 level deep
			[]string{"symbolic-ref", symbolicRef1, "refs/heads/" + branchName},

			// add a symbolic ref, 2 levels deep
			[]string{"symbolic-ref", symbolicRef2, symbolicRef1},

			// pack the refs
			[]string{"pack-refs"},
		)
		if err != nil {
			return err
		}

		// get the output data
		info := &InfoRefs{
			CommitOid:          util.RevOid(repo, "HEAD"),
			AnnTagName:         annTagName,
			AnnTagOid:          util.RevOid(repo, annTagName),
			LightTagName:       lightTagName,
			BranchName:         branchName,
			SymbolicRef1:       symbolicRef1,
			SymbolicRef1Target: "refs/heads/" + branchName,
			SymbolicRef2:       symbolicRef2,
			SymbolicRef2Target: symbolicRef1,
			HeadTarget:         "refs/heads/master",
		}

		testCase.info = info
		return
	},
)
