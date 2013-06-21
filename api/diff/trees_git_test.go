//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
trees_git_test.go implements git-comparison tests for ggit tree parsing.
*/
package diff

import (
	"fmt"
	"github.com/jbrukh/ggit/api"
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/test"
	"github.com/jbrukh/ggit/util"
	"testing"
)

// Test_readCommits will compare the commit output of
// git and ggit for a string of commits.
func Test_readTree(t *testing.T) {
	testCase := test.TreeDiff
	repo := api.Open(testCase.Repo())
	info, _ := testCase.Info().(*test.InfoTreeDiff)

	oid := objects.OidNow(info.TreeOid1)
	o, err := repo.ObjectFromOid(oid)
	util.AssertNoErr(t, err)
	// get the tree
	tree1, _ := o.(*objects.Tree)

	oid = objects.OidNow(info.TreeOid2)
	o, err = repo.ObjectFromOid(oid)
	util.AssertNoErr(t, err)
	// get the tree
	tree2, _ := o.(*objects.Tree)

	diff, err := NewTreeDiffer(repo).Diff(tree1, tree2)
	modified := diff.Modified()

	if len(modified) != 1 {
		t.Errorf("expected 1 modification in the tree diff. got %d", len(modified))
	}

	edit := modified[0]
	before, after := edit.Before, edit.After
	beforeOid, afterOid := before.ObjectId(), after.ObjectId()
	beforeName := before.Name()

	if beforeName != info.ModifiedFileName {
		t.Errorf("expected modified file name to be [%s], but got [%s]", info.ModifiedFileName, beforeName)
	}

	if beforeOid.String() != info.ModifiedFileBeforeOid {
		t.Errorf("expected modified file's initial id to be [%s], but got [%s]", info.ModifiedFileBeforeOid, beforeOid)
	}

	if afterOid.String() != info.ModifiedFileAfterOid {
		t.Errorf("expected modified file's initial id to be [%s], but got [%s]", info.ModifiedFileAfterOid, afterOid)
	}

	deleted := diff.Deleted()

	util.Assert(t, len(deleted) == 1, fmt.Sprintf("expected 1 file to be deleted, got [%d]", len(deleted)))
	entry := deleted[0]
	deletedName, deletedOid := entry.Name(), entry.ObjectId().String()
	if deletedName != info.RemovedFileName {
		t.Errorf("expected deleted file's name to be [%s], but got [%s]", info.RemovedFileName, deletedName)
	}
	if deletedOid != info.RemovedFileOid {
		t.Errorf("expected deleted file's id to be [%s], but got [%s]", info.RemovedFileOid, deletedOid)
	}
}
