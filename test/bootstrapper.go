package test

import (
	"fmt"
	"github.com/jbrukh/ggit/api"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// Repo creates a directory in "root" and a subdirectory
// therein where it executes the test case which builds up
// a test repository. The test case script is handed the
// destination directory as its sole parameter. That directory
// is returned.
// The name of the script can be a nickname, such as "empty_repo",
// which will automatically resolve to "test/cases/empty_repo.sh".
// Otherwise, the parameter will be regarded as a path to the
// file.
func Repo(root string, script string) (string, error) {
	dir := path.Join(root, intuitName(script))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	if api.IsValidRepo(dir) {
		// already exists
		return "", fmt.Errorf("the repo '%s' already exists", dir)
	}

	cmd := exec.Command(resolvePath(script), dir)
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

func resolvePath(script string) string {
	// if it looks like a nickname, resolve automatically
	if !strings.HasSuffix(script, ".sh") && strings.Index(script, "/") < 0 {
		// TODO: relative path makes this a bug
		return fmt.Sprintf("test/cases/%s.sh", script)
	}
	return script
}
