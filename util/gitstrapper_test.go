package util

import (
	"github.com/jbrukh/ggit/api"
	"testing"
)

func Test_CreateRepo(t *testing.T) {
	dir := "var/test_repo"
	_, err := CreateGitRepo(dir)
	AssertNoErr(t, err)
	Assert(t, api.IsValidRepo(dir))
	AssertNoErr(t, RemoveGitRepo(dir))
}
