package util

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
)

func CreateGitRepo(dir string) (string, error) {
	// ensure the directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return GitExec(dir, "init")
}

func RemoveGitRepo(dir string) error {
	return os.RemoveAll(dir)
}

// IsValidRepo validates a repository path to make sure it has
// the right format and that it exists.	
func IsValidRepo(pth string) bool {
	p := inferGitDir(pth)
	if _, e := os.Stat(p); e != nil {
		return false
	}
	// TODO: may want to do other checks here...
	return true
}

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
