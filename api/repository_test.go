package api

import (
	"os"
	"path"
	"testing"
)

func Test_Open(t *testing.T) {
	assert(t, Open("test").path == "test/.git")
	assert(t, Open("test/.git").path == "test/.git")
}

func Test_IsValidRepo(t *testing.T) {
	var (
		repo   = "test-000"
		gitDir = path.Join(repo, ".git")
	)
	err := os.MkdirAll(gitDir, 0755)
	assertNoErr(t, err)

	assert(t, IsValidRepo(repo))
	assert(t, IsValidRepo(gitDir))

	err = os.RemoveAll(repo)
	assertNoErr(t, err)
}
