// Code in this package originally based on https://github.com/jordanorelli/multicommand.
package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/jbrukh/ggit/api"
	"github.com/jbrukh/ggit/builtin"
	"io"
	"os"
	"path"
)

var (
	flagVersion     bool
	flagWhichGitDir bool

	Wout = os.Stdout
	Werr = os.Stderr
)

func init() {
	flag.BoolVar(&flagVersion, "version", false, "")
	flag.BoolVar(&flagWhichGitDir, "which-git-dir", false, "show the path of the enclosing repo")
}

// ================================================================= //
// GGIT COMMAND
// ================================================================= //

func main() {
	flag.Usage = usage
	flag.Parse()

	// --which-git-dir
	if flagWhichGitDir {
		path, err := findRepo()
		if err != nil {
			fmt.Fprintf(Werr, "Could not discern parent repo: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Fprintln(Wout, path)
		os.Exit(0)
	}

	// --version
	if flagVersion {
		fmt.Fprintln(Wout, "ggit version", Version)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	// what builtin are we trying to call?
	name, args := args[0], args[1:]

	// get the builtin
	cmd, ok := builtin.Get(name)
	if ok {
		repo, e := openRepo()
		if e != nil {
			fmt.Fprintln(Werr, msgNotARepo)
			os.Exit(1)
		}
		cmd.Execute(&builtin.Params{
			repo,
			os.Stdout,
			os.Stderr,
		}, args)
	} else {
		fmt.Fprintf(os.Stderr, fmtUnknownCommand, name)
		usage()
	}
}

func findRepo() (string, error) {
	dir, e := os.Getwd()
	if e != nil {
		return "", e
	}
	var file string
	for {
		// check this directory to see if it contains
		// the git directory, usually .git
		gitDir := path.Join(dir, api.DefaultGitDir)
		if _, e = os.Stat(gitDir); os.IsNotExist(e) {
			// try the directory up
			dir, file = path.Split(dir)
			if file == "" { // nothing more to go up
				return "", errors.New("no repo found")
			}
		} else if e != nil {
			// there is some other error
			break
		} else {
			// no, error the git dir exists
			return gitDir, nil
		}
	}
	return "", e
}

// openRepo opens the repository, if any, which is
// the enclosing repository of this directory.
func openRepo() (repo api.Repository, err error) {
	var path string
	path, err = findRepo()
	if err != nil {
		return nil, err
	}

	if repo, err = api.Open(path); err != nil {
		return nil, err
	}
	return
}

// ================================================================= //
// GGIT USAGE
// ================================================================= //

func usage() {
	printUsage(Wout)
	os.Exit(2)
}

func printUsage(w io.Writer) {
	tmpl(w, tmplUsage, builtin.All())
}
