package main

import (
	"flag"
	"fmt"
	"github.com/jbrukh/ggit/test"
	"os"
)

// ggcase lets you play around with test repository cases. It will
// create that case that you specify in the working directory.
func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stdout, "usage: ggcase <case_file>")
		return
	}
	caseFile := args[0]
	dir, err := test.Repo(".", caseFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Fprintf(os.Stdout, "Created case '%s' in '%s'\n", caseFile, dir)
}
