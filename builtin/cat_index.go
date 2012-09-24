package builtin

import (
	"fmt"
	"github.com/jbrukh/ggit/api"
	"io"
	"os"
)

// var (
// 	fType, fPrint, fSize bool
// )

var catIndexBuiltin = &Builtin{
	Execute:     catIndex,
	Name:        "cat-index",
	Description: "Provide a debug dump of the index file",
	UsageLine:   "cat-index",
	ManPage:     "TODO",
}

func init() {
	// add to command list
	Add(catIndexBuiltin)
}

func catIndex(b *Builtin, args []string, path string, w io.Writer) {
	repo, e := api.Open(path)
	if e != nil {
		fmt.Fprintf(os.Stderr, "could not open repo: %s", path)
		return
	}
	defer repo.Close()

	inx, e := repo.Index()
	if e != nil {
		return
	}
	fmt.Fprint(w, inx)
}
