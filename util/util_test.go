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

func Test_IsDigit(t *testing.T) {
	Assert(t, IsDigit('0'))
	Assert(t, IsDigit('1'))
	Assert(t, IsDigit('2'))
	Assert(t, IsDigit('3'))
	Assert(t, IsDigit('4'))
	Assert(t, IsDigit('5'))
	Assert(t, IsDigit('6'))
	Assert(t, IsDigit('7'))
	Assert(t, IsDigit('8'))
	Assert(t, IsDigit('9'))
	Assert(t, !IsDigit('z'))
}
