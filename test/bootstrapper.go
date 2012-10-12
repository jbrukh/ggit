package test

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// CreateTestRepo creates a temporary directory and a subdirectory where
// a test repo will be created. It passes this path to a script which
// it executes. It then returns the resulting directory.
func Repo(root string, script string) (string, error) {
	dir := path.Join(root, intuitName(script))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	cmd := exec.Command(script, dir)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return dir, nil
}

// intuitName returns the name of the shell script that is being used
// without any tail extention. For example, "test/cases/x.sh" would result
// in the output "x".
func intuitName(script string) string {
	_, file := filepath.Split(script)
	if file == "" {
		panic("must provide a path to a shell script")
	}
	toks := strings.Split(file, ".")
	if len(toks) == 1 {
		return toks[0]
	}
	return toks[len(toks)-2]
}
