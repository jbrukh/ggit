package main

import (
	"flag"
	"fmt"
	"github.com/jbrukh/ggit/test"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stdout, "usage: ggtest <case_file>")
		return
	}
	caseFile := args[0]

	dir, err := test.Repo(".", caseFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Fprintf(os.Stdout, "Created case '%s' in '%s", caseFile, dir)
}
