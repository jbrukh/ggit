// Code in this package originally based on https://github.com/jordanorelli/multicommand.
package main

import (
	"flag"
	"fmt"
	"github.com/jbrukh/ggit/api"
	"github.com/jbrukh/ggit/builtin"
	"io"
	"os"
)

var (
	flagVersion bool

	Wout = os.Stdout
	Werr = os.Stderr
)

func init() {
	flag.BoolVar(&flagVersion, "version", false, "")
}

// ================================================================= //
// GGIT COMMAND
// ================================================================= //

func main() {
	flag.Usage = usage
	flag.Parse()

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
			fmt.Println(e.Error())
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
	return api.DefaultGitDir, nil
}

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
