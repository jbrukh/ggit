package util

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func Test_AssertCreateRemove(t *testing.T) {
	repo := TempRepo("test001")
	AssertCreateGitRepo(t, repo)
	AssertRemoveGitRepo(t, repo)
}

func Test_GitExec(t *testing.T) {
	repo := TempRepo("test002")
	AssertCreateGitRepo(t, repo)
	defer AssertRemoveGitRepo(t, repo)

	// status on empty dir
	var out string
	var err error
	out, err = GitExec(repo, "status")
	AssertNoErr(t, err)
	Assert(t, out == emptyRepoStatus)

	// add a test file
	var testFile = path.Join(repo, "test.txt")
	err = ioutil.WriteFile(testFile, []byte("hahaha"), 0644)
	AssertNoErr(t, err)

	// status with test file
	out, err = GitExec(repo, "status")
	AssertNoErr(t, err)
	Assert(t, out == newFileStatus)

	// hash an object in the repo's object db
	out, err = GitExec(repo, "hash-object", "-w", testFile)
	AssertNoErr(t, err)
	Assert(t, strings.TrimSpace(out) == oidOfTestFile)

	// hash an object in the repo's object db, using HashBlob
	var oid string
	oid, err = HashBlob(repo, "hahaha")
	AssertNoErr(t, err)
	Assert(t, oid == oidOfTestFile)
}

func Test_TestFile(t *testing.T) {
	repo := TempRepo("test_file_add")
	AssertCreateGitRepo(t, repo)
	defer AssertRemoveGitRepo(t, repo)

	testFile := "hello"
	testContents := "hey!"
	err := TestFile(repo, testFile, testContents)
	AssertNoErrOrDie(t, err)

	pth := path.Join(repo, testFile)
	file, err := os.Open(pth)
	AssertNoErr(t, err)

	var contents bytes.Buffer
	io.Copy(&contents, file)
	AssertEqualString(t, contents.String(), testContents)
}

const emptyRepoStatus = `# On branch master
#
# Initial commit
#
nothing to commit (create/copy files and use "git add" to track)
`

const newFileStatus = `# On branch master
#
# Initial commit
#
# Untracked files:
#   (use "git add <file>..." to include in what will be committed)
#
#	test.txt
nothing added to commit but untracked files present (use "git add" to track)
`

const oidOfTestFile = "1240583197c9a4507a2fb0d59eb1a82886844e57"
