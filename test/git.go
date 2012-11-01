//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	//"path/filepath"
	"github.com/jbrukh/ggit/util"
	"strings"
	"testing"
)

// TempRepo returns a temporary location where we 
// can store the test repo.
func TempRepo(subdir string) string {
	return path.Join(os.TempDir(), subdir)
}

// CreateGitRepo creates an empty git repo in the
// specified location.
func CreateGitRepo(dir string) (string, error) {
	// ensure the directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return GitExec(dir, "init")
}

// GitExec executes a git command via the shell in
// the given workDir. The string returned is the
// output of the git command.
func GitExec(workDir string, args ...string) (string, error) {
	// execute the git command
	gitDir := path.Join(workDir, ".git")
	gitDirArg := fmt.Sprintf("--git-dir=%s", gitDir)
	workDirArg := fmt.Sprintf("--work-tree=%s", workDir)
	args = append([]string{"git", gitDirArg, workDirArg}, args...)

	// print the output

	fmt.Println(gitDir, ": ", strings.Join(args[3:], " "))

	cmd := exec.Command(args[0], args[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}

// GitExecMany executes multiple git commands in the 
// given repo (simular to GitExec) but swallows the
// output and returns on any error. This is meant for
// setting up test scenarios.
func GitExecMany(workDir string, cmds ...[]string) error {
	for i, cmd := range cmds {
		_, err := GitExec(workDir, cmd...)
		if err != nil {
			return fmt.Errorf("Failed on command %d ('%s') %s", i, cmd, err)
		}
	}
	return nil
}

// AssertCreateGitRepo is a convenience method for testing
// which creates a new repo and asserts that it was
// created successfully.
func AssertCreateGitRepo(t *testing.T, repo string) {
	_, err := CreateGitRepo(repo)
	util.AssertNoErr(t, err)
	util.Assert(t, util.IsValidRepo(repo))
}

// AssertRemoveGitRepo is a convenience method for testing
// which removes a new repo and asserts that is was
// removed successfully.
func AssertRemoveGitRepo(t *testing.T, repo string) {
	err := os.RemoveAll(repo)
	util.AssertNoErr(t, err)
}

// TestFile creates a file with name "name" inside of the
// repo "repo" with the specified contents.
func TestFile(repo string, name string, contents string) error {
	pth := path.Join(repo, name)
	err := ioutil.WriteFile(pth, []byte(contents), 0644)
	return err
}

// HashBlob creates a new blob object in the odb of the
// given repository.
func HashBlob(repo string, contents string) (oid string, err error) {
	if !util.IsValidRepo(repo) {
		return "", fmt.Errorf("does not appear to be a valid repo: %s", repo)
	}
	name := path.Join(os.TempDir(), util.UniqueHex16())
	err = ioutil.WriteFile(name, []byte(contents), 0644)
	if err != nil {
		return "", err
	}
	oid, err = GitExec(repo, "hash-object", "-w", name)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(oid), err
}

// TODO: turn the following two methods into GitNow()

// RevOid returns the git-rev-parse of the current
// revision and returns it as a 40-character string.
func RevOid(repo string, rev string) string {
	oid, err := GitExec(repo, "rev-parse", rev)
	if err != nil {
		msg := fmt.Sprintf("can't get oid for: %s (%s)", rev, err)
		panic(msg)
	}
	return strings.TrimSpace(oid)
}

// ObjectRepr returns the git-cat-file -p output for
// the given revision, or panics if there is an
// error
func ObjectRepr(repo string, rev string) string {
	oid, err := GitExec(repo, "cat-file", "-p", rev)
	if err != nil {
		panic("can't get repr for: " + rev)
	}
	return oid
}
