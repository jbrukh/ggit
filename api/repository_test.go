package api

import (
	"github.com/jbrukh/ggit/util"
	"os"
	"path"
	"testing"
)

func Test_Open(t *testing.T) {
	util.Assert(t, Open("test").path == "test/.git")
	util.Assert(t, Open("test/.git").path == "test/.git")
}

func Test_IsValidRepo(t *testing.T) {
	var (
		repo   = "test-000"
		gitDir = path.Join(repo, ".git")
	)
	err := os.MkdirAll(gitDir, 0755)
	util.AssertNoErr(t, err)

	util.Assert(t, util.IsValidRepo(repo))
	util.Assert(t, util.IsValidRepo(gitDir))

	err = os.RemoveAll(repo)
	util.AssertNoErr(t, err)
}
