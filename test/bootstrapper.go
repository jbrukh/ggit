package test

import (
	"github.com/jbrukh/ggit/api"
	"os"
	"os/exec"
	"path"
)

const varDir = "var"

func CreateRepo(script string, name string) (api.Repository, error) {
	dir := path.Join(varDir, name)
	if e := os.MkdirAll(dir, 0755); e != nil {
		return nil, e
	}
	cmd := exec.Command(string, dir)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return api.Open(dir)
}

func CleanRepo(name string) error {
	dir := path.Join(varDir, name)
	return os.Remove(dir)
}
