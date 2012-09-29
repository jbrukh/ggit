package builtin

import (
	"fmt"
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

func showRef(p *Params, b *Builtin, args []string) {
	r, e := p.Repo.PackedRefs()
	if e != nil {
		fmt.Fprint(p.Wout, "could not read packed refs")
		return
	}
	for _, v := range r {
		fmt.Fprintf(p.Wout, "%s\n", v.String())
	}
}
