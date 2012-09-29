package builtin

import (
	"fmt"
)

func init() {
	// add to command list
	Add(ShowRef)
}

type ShowRefBuiltin struct {
	HelpInfo
}

var ShowRef = &ShowRefBuiltin{
	HelpInfo: HelpInfo{
		Name:        "show-ref",
		Description: "List references in a local repository",
		UsageLine:   "show-ref",
		ManPage:     "TODO",
	},
}

func (b *ShowRefBuiltin) Info() *HelpInfo {
	return &b.HelpInfo
}

func (b *ShowRefBuiltin) Execute(p *Params, args []string) {
	//fmt.Println("getting ", args)
	if len(args) > 0 {
		// we want to see a particular ref
		which := args[0]
		ref, e := p.Repo.ReadRef(which)
		if e != nil {
			// nothing to show
			return
		}
		fmt.Fprint(p.Wout, ref.String())
	} else {
		r, e := p.Repo.PackedRefs()
		if e != nil {
			fmt.Fprint(p.Wout, "could not read packed refs")
			return
		}
		for _, v := range r {
			fmt.Fprintf(p.Wout, "%s\n", v.String())
		}
	}
}
