package builtin

import (
	"fmt"
	"github.com/jbrukh/ggit/api"
	"io"
)

// var (
// 	fType, fPrint, fSize bool
// )

var showRefBuiltin = &Builtin{
	Execute:     showRef,
	Name:        "show-ref",
	Description: "Provide a debug dump of the index file",
	UsageLine:   "cat-index",
	ManPage:     "TODO",
}

func init() {
	// add to command list
	Add(catIndexBuiltin)
}

func catIndex(b *Builtin, args []string, repo api.Repository, w io.Writer) {
	inx, e := repo.Index()
	if e != nil {
		return
	}
	fmt.Fprint(w, inx)
}
