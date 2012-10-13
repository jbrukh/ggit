package main

import (
	"flag"
	"fmt"
	"github.com/jbrukh/ggit/api"
	"github.com/jbrukh/ggit/test"
	"os"
)

var flagPrintObjects bool

func init() {
	flag.BoolVar(&flagPrintObjects, "objects", false, "print Go code for repo objects")
}

const usage = `Generate cases:

    ggcase (<case_script>|<nickname>)

Generate consts:

    ggcase --objects <repo>
`

// ggcase lets you play around with test repository cases. It will
// create that case that you specify in the working directory.
func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	if flagPrintObjects {
		printObjects(args[0])
	} else {
		generateCase(args[0])
	}
}

func generateCase(caseFile string) {
	dir, err := test.Repo(".", caseFile)
	if err != nil {
		die("%s\n", err.Error())
	}
	fmt.Fprintf(os.Stdout, "created case '%s' in '%s'\n", caseFile, dir)
}

func printObjects(dir string) {
	if !api.IsValidRepo(dir) {
		die("'%s' doesn't seem to be a valid repo.\n", dir)
	}
	repo := api.Open(dir)
	oids, err := repo.ObjectIds()
	if err != nil {
		die("could not get object id's")
	}

	for _, oid := range oids {
		fmt.Println(oid.String())
	}
}

func die(format string, items ...interface{}) {
	fmt.Fprintf(os.Stderr, format, items)
	os.Exit(1)
}
