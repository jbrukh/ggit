//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
git.go implements functions for interactive with git repositories using
command-line git.
*/
package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// check to make sure git is installed on this system,
// for testing purposes
func init() {
	cmd := exec.Command("git", "--version")
	cmd.Stdout = nil
	if err := cmd.Run(); err != nil {
		panic("no git installed: " + err.Error())
	}
}

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

	fmt.Printf("%s: %s\n", gitDir, strings.Join(args[3:], " "))

	cmd := exec.Command(args[0], args[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return out.String(), err
	}
	return out.String(), nil
}

// GitExecMany executes multiple git commands in the 
// given repo (simular to GitExec) but swallows the
// output and returns on any error. This is meant for
// setting up test scenarios.
func GitExecMany(workDir string, cmds ...[]string) error {
	for i, cmd := range cmds {
		out, err := GitExec(workDir, cmd...)
		if err != nil {
			return fmt.Errorf("Failed on command %d ('%s') %s (git: '%s')", i, cmd, err, out)
		}
	}
	return nil
}

// TestFile creates a file with name "name" inside of the
// repo "repo" with the specified contents.
func TestFile(repo string, name string, contents string) error {
	pth := path.Join(repo, name)

	// create the directory, if applicable
	dir, _ := filepath.Split(pth)
	if dir != "" {
		os.MkdirAll(dir, 0755)
	}

	err := ioutil.WriteFile(pth, []byte(contents), 0644)
	if err != nil {
		return fmt.Errorf("could not create test file '%s' for repo: %s", name, err)
	}
	return nil
}

// HashBlob creates a new blob object in the odb of the
// given repository.
func HashBlob(repo string, contents string) (oid string, err error) {
	if !IsValidRepo(repo) {
		return "", fmt.Errorf("does not appear to be a valid repo: %s", repo)
	}
	name := path.Join(os.TempDir(), UniqueHex16())
	err = ioutil.WriteFile(name, []byte(contents), 0644)
	if err != nil {
		return "", err
	}
	defer func() {
		err := os.Remove(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error removing file: %s\n", err)
		}
	}()
	oid, err = GitExec(repo, "hash-object", "-w", name)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(oid), err
}

// ================================================================= //
// UTILS
// ================================================================= //

// GitNow executes a git command or panics
func GitNow(repo string, params ...string) string {
	out, err := GitExec(repo, params...)
	if err != nil {
		msg := fmt.Sprintf("can't execute cmd ('%s') in %s (%s)", strings.Join(params, " "), repo, err)
		panic(msg)
	}
	return out
}

func TreeEntryOid(repo string, file string) string {
	line := GitNow(repo, "ls-files", "-s", file)
	return line[7:47] // get out the sha
}

// RevOid returns the git-rev-parse of the current
// revision and returns it as a 40-character string.
func RevOid(repo string, rev string) string {
	oid := GitNow(repo, "rev-parse", rev)
	return strings.TrimSpace(oid)
}

func ObjectType(repo string, rev string) string {
	otype := GitNow(repo, "cat-file", "-t", rev)
	return strings.TrimSpace(otype)
}

// ObjectRepr returns the git-cat-file -p output for
// the given revision, or panics if there is an
// error
func ObjectRepr(repo string, rev string) (repr string) {
	otype := ObjectType(repo, rev)
	return GitNow(repo, "cat-file", otype, rev)
}

func ObjectSize(repo string, rev string) int {
	size := strings.TrimSpace(GitNow(repo, "cat-file", "-s", rev))
	i, err := strconv.Atoi(size) // should be an int, or git is busted
	if err != nil {
		panic(err.Error())
	}
	return i
}
