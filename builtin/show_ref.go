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
	Description: "List references in a local repository",
	UsageLine:   "show-ref",
	ManPage:     "TODO",
}

func init() {
	// add to command list
	Add(showRefBuiltin)
}

func showRef(b *Builtin, args []string, repo api.Repository, w io.Writer) {
	r, e := repo.PackedRefs()
	if e != nil {
		fmt.Fprint(w, "could not read packed refs")
		return
	}
	for _, v := range r {
		fmt.Fprintf(w, "%s\n", v.String())
	}
}
