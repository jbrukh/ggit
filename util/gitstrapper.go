package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

// TestDir returns a temporary location where we 
// can store the test repo.
func TestDir(subdir string) string {
	return path.Join( /*os.TempDir()*/ "var", subdir)
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
	fmt.Println(strings.Join(args, " "))

	cmd := exec.Command(args[0], args[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// GitExecMany executes multiple git commands in the 
// given repo (simular to GitExec) but swallows the
// output and returns on any error. This is meant for
// setting up test scenarios.
func GitExecMany(workDir string, cmds ...[]string) error {
	for _, cmd := range cmds {
		_, err := GitExec(workDir, cmd...)
		if err != nil {
			return err
		}
	}
	return nil
}
