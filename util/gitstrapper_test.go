package util

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var (
	testDir1  string
	testDir2  string
	testFile1 string
)

func init() {
	testDir1 = TestDir("test001")
	testDir2 = TestDir("test002")
	testFile1 = path.Join(testDir2, "test.txt")
}

func Test_CreateRepo(t *testing.T) {
	_, err := CreateGitRepo(testDir1)
	AssertNoErr(t, err)
	Assert(t, IsValidRepo(testDir1))
}

func Test_GitExec(t *testing.T) {
	_, err := CreateGitRepo(testDir2)
	AssertNoErr(t, err)
	Assert(t, IsValidRepo(testDir2))

	// status on empty dir
	var out string
	out, err = GitExec(testDir2, "status")
	AssertNoErr(t, err)
	Assert(t, out == emptyRepoStatus)

	// add a test file
	err = ioutil.WriteFile(testFile1, []byte("hahaha"), 0644)
	AssertNoErr(t, err)

	// status with test file
	out, err = GitExec(testDir2, "status")
	AssertNoErr(t, err)
	Assert(t, out == newFileStatus)

	// hash an object in the repo's object db
	out, err = GitExec(testDir2, "hash-object", "-w", testFile1)
	AssertNoErr(t, err)
	Assert(t, out == oidOfTestFile)

	// hash an object in the repo's object db, using HashBlob
	var oid string
	oid, err = HashBlob(testDir2, testFile1, "hahaha")
	AssertNoErr(t, err)
	Assert(t, oid == oidOfTestFile)
}

// NOTE: this method should be run last
// for cleanup purposes. There may be
// a better way of doing this in Go.
func Test_CLEANUP(t *testing.T) {
	os.RemoveAll(testDir1)
	os.RemoveAll(testDir2)
}

const emptyRepoStatus = `# On branch master
#
# Initial commit
#
nothing to commit (create/copy files and use "git add" to track)`

const newFileStatus = `# On branch master
#
# Initial commit
#
# Untracked files:
#   (use "git add <file>..." to include in what will be committed)
#
#	test.txt
nothing added to commit but untracked files present (use "git add" to track)`

const oidOfTestFile = `1240583197c9a4507a2fb0d59eb1a82886844e57`
