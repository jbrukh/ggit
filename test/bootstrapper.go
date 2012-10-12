package test

import (
	"github.com/jbrukh/ggit/api"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const varDir = "var"

// CreateTestRepo creates a temporary directory and a subdirectory where
// a test repo will be created. It passes this path to a script which
// it executes. It then returns a ggit Repository based on that directory.
func Repo(script string) (api.Repository, error) {
	dir := path.Join(varDir, intuitName(script))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	cmd := exec.Command(script, dir)
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return api.Open(dir)
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
