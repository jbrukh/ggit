package util

import (
	"os"
	"path"
	"testing"
)

var (
	testDir  string
	testFile string
)

func init() {
	testDir = TestDir("test000")
	testFile = path.Join(testDir, "test.txt")
}

func Test_CreateRepo(t *testing.T) {
	_, err := CreateGitRepo(testDir)
	AssertNoErr(t, err)
	Assert(t, IsValidRepo(testDir))
}

// NOTE: this method should be run last
// for cleanup purposes. There may be
// a better way of doing this in Go.
func Test_CLEANUP(t *testing.T) {
	os.RemoveAll(testDir)
}
