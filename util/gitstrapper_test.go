package util

import (
	"bytes"
	"github.com/jbrukh/ggit/api"
	"io"
	"os"
	"path"
	"testing"
)

var (
	testDir  string
	testFile string
)

func init() {
	testDir = path.Join(os.TempDir(), "ggit_test", "test000")
	testFile = path.Join(testDir, "test.txt")
}

func Test_CreateRepo(t *testing.T) {
	_, err := CreateGitRepo(testDir)
	AssertNoErr(t, err)
	Assert(t, api.IsValidRepo(testDir))
	AssertNoErr(t, RemoveGitRepo(testDir))
}
