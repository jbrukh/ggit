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
		UsageLine:   "cat-index",
		ManPage:     "TODO",
	},
}

func (b *CatIndexBuiltin) Info() *HelpInfo {
	return &b.HelpInfo
}

func (b *CatIndexBuiltin) Execute(p *Params, args []string) {
	inx, e := p.Repo.Index()
	if e != nil {
		return
	}
	fmt.Fprint(p.Wout, inx)
}
