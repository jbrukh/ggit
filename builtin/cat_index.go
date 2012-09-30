package builtin

import (
	"fmt"
)

func init() {
	// add to command list
	Add(CatIndex)
}

type CatIndexBuiltin struct {
	HelpInfo
}

var CatIndex = &CatIndexBuiltin{
	HelpInfo: HelpInfo{
		Name:        "cat-index",
		Description: "Provide a debug dump of the index file",
		UsageLine:   "",
		ManPage:     "TODO",
	},
}

func (b *CatIndexBuiltin) Execute(p *Params, args []string) {
	inx, e := p.Repo.Index()
	if e != nil {
		return
	}
	fmt.Fprint(p.Wout, inx)
}
