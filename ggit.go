// Code in this package originally based on https://github.com/jordanorelli/multicommand.
package main

import (
	"flag"
	"fmt"
	"github.com/jbrukh/ggit/builtin"
	"io"
	"os"
)

// ================================================================= //
// GGIT COMMAND
// ================================================================= //

func main() {
	flag.Usage = usage
	flag.Parse()
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
			cmd.Usage()
		}

		cmd.FlagSet.Parse(args[1:])
		args = cmd.FlagSet.Args()

		cmd.Execute(cmd, args)
		return

	}

	fmt.Fprintf(os.Stderr, unknownCommandFormat, name)
	usage()
}

// ================================================================= //
// GGIT USAGE
// ================================================================= //

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

func printUsage(w io.Writer) {
	tmpl(w, api.UsageTemplate, api.Builtins())
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
