package builtin

import (
	"fmt"
)

func init() {
	// add to command list
	Add(RevParse)
}

type RevParseBuiltin struct {
	HelpInfo
}

var RevParse = &RevParseBuiltin{
	HelpInfo: HelpInfo{
		Name:        "rev-parse",
		Description: "Translate a revision specification into a SHA1 object id",
		UsageLine:   "",
		ManPage:     "TODO",
	},
}

func (b *RevParseBuiltin) Execute(p *Params, args []string) {
	if len(args) != 1 {
		b.Usage(p.Werr)
		return
	}
	spec := args[0]
	fmt.Fprintln(p.Wout, "doing: rev-parse TODO: ", spec)
}
