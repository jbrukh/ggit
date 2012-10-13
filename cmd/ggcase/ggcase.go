package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/jbrukh/ggit/api"
	"github.com/jbrukh/ggit/test"
	"os"
	"sort"
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

	var err error
	if flagPrintObjects {
		err = printObjects(args[0])
	} else {
		err = generateCase(args[0])
	}

	if err != nil {
		die(err.Error())
	}
}

func generateCase(caseFile string) (err error) {
	dir, err := test.Repo(".", caseFile)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	fmt.Fprintf(os.Stdout, "created case '%s' in '%s'\n", caseFile, dir)
	return
}

// sort interface for sorting refs
type objectByType []api.Object

func (s objectByType) Len() int           { return len(s) }
func (s objectByType) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s objectByType) Less(i, j int) bool { return s[i].Header().Type() < s[j].Header().Type() }

func printObjects(dir string) (err error) {
	if !api.IsValidRepo(dir) {
		return fmt.Errorf("'%s' doesn't seem to be a valid repo", dir)
	}
	repo := api.Open(dir)
	oids, err := repo.ObjectIds()
	if err != nil {
		return fmt.Errorf("could not get object id's")
	}

	var objs []api.Object
	objs, err = getObjects(repo, oids)
	if err != nil {
		return err
	}

	// sort by type
	sort.Sort(objectByType(objs))

	counts := make(map[api.ObjectType]int)
	counts[api.ObjectBlob] = 0
	counts[api.ObjectTree] = 0
	counts[api.ObjectCommit] = 0
	counts[api.ObjectTag] = 0

	buf := bytes.NewBufferString("")
	fmt.Fprint(buf, "var (\n")
	for _, o := range objs {
		otype := o.Header().Type()
		cnt := counts[otype]
		cnt++
		counts[otype] = cnt
		fmt.Fprintf(buf, "\ttest_%s_%d %20s api.OidNow(\"%s\")\n", otype, cnt, "=", o.ObjectId())
	}
	fmt.Fprint(buf, ")")
	fmt.Println(buf)
	return
}
func getObjects(repo api.Repository, oids []*api.ObjectId) ([]api.Object, error) {
	var objs []api.Object
	for _, oid := range oids {
		o, err := repo.ObjectFromOid(oid)
		if err != nil {
			return nil, fmt.Errorf("cannot get object: %s", oid)
		}
		objs = append(objs, o)
	}
	return objs, nil
}

func die(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
