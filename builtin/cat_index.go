package builtin

import (
	"fmt"
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

func catIndex(p *Params, b *Builtin, args []string) {
	inx, e := p.Repo.Index()
	if e != nil {
		return
	}
	fmt.Fprint(p.Wout, inx)
}
