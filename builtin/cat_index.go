package builtin

import (
	"fmt"
	"github.com/jbrukh/ggit/api"
	"io"
)

// var (
// 	fType, fPrint, fSize bool
// )

var catIndexBuiltin = &Builtin{
	Execute:     CatIndex,
	Name:        "cat-index",
	Description: "Provide a debug dump of the index file",
	UsageLine:   "cat-index",
	ManPage:     "TODO",
}

func init() {
	// add to command list
	Add(catIndexBuiltin)
}

func CatIndex(b *Builtin, args []string, repo api.Repository, w io.Writer) {
	inx, e := repo.Index()
	if e != nil {
		return
	}
	fmt.Fprint(w, inx)
}
