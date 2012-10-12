package test

import (
	"github.com/jbrukh/ggit/api"
	"os"
	"os/exec"
	"path"
)

const varDir = "var"

// CreateTestRepo creates a temporary directory and a subdirectory where
// a test repo will be created. It passes this path to a script which
// it executes. It then returns a ggit Repository based on that directory.
func CreateTestRepo(script string, name string) (api.Repository, error) {
	dir := path.Join(varDir, name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	cmd := exec.Command(string, dir)
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return api.Open(dir)
}

func CleanTestRepo(name string) error {
	dir := path.Join(varDir, name)
	return os.Remove(dir)
}
