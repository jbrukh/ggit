package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
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
	args = append([]string{gitDirArg, workDirArg}, args...)

	fmt.Println(args)

	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
