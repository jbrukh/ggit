package util

import (
	"os"
	"path"
	"testing"
)

func Test_IsValidRepo(t *testing.T) {
	var (
		repo   = "test-000"
		gitDir = path.Join(repo, ".git")
	)
	err := os.MkdirAll(gitDir, 0755)
	AssertNoErr(t, err)

	Assert(t, IsValidRepo(repo))
	Assert(t, IsValidRepo(gitDir))

	err = os.RemoveAll(repo)
	AssertNoErr(t, err)
}
