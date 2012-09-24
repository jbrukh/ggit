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

var fVersion bool

func init() {
	flag.BoolVar(&fVersion, "version", false, "")
}

// ================================================================= //
// GGIT COMMAND
// ================================================================= //

func main() {
	flag.Usage = usage
	flag.Parse()

	if fVersion {
		fmt.Println("ggit version", Version)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	// what builtin are we trying to call?
	name := args[0]

	// TODO make this into a command
	if name == "help" {
		help(args[1:])
		return
	}

	// get the builtin
	cmd, ok := builtin.Get(name)
	if ok {
		cmd.FlagSet.Usage = func() {
			cmd.Usage(os.Stderr)
		}
		cmd.FlagSet.Parse(args[1:])
		args = cmd.FlagSet.Args()

		path, e := findRepo()
		if e != nil {
			fmt.Println(msgNotARepo)
		}
		cmd.Execute(cmd, args, path, os.Stdout)
	} else {
		fmt.Fprintf(os.Stderr, fmtUnknownCommand, name)
		usage()
	}
}

func findRepo() (string, error) {
	return api.DEFAULT_GIT_DIR, nil
}

// ================================================================= //
// GGIT USAGE
// ================================================================= //

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

func printUsage(w io.Writer) {
	tmpl(w, tmplUsage, builtin.All())
}

// ================================================================= //
// GGIT HELP
// ================================================================= //

// help implements the 'help' command
func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		return
	}

	name := args[0]
	println("TODO: look up help for: ", name)

	// for _, cmd := range api.builtin {
	// 	if cmd.Name() == name {
	// 		tmpl(os.Stdout, helpTemplate, cmd)
	// 		// not exit 2: succeeded at 'go help cmd'.
	// 		return
	// 	}
	// }
}
