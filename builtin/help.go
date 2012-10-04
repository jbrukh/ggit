package builtin

import (
	"fmt"
)

func init() {
	// add to command list
	Add(Help)
}

type HelpBuiltin struct {
	HelpInfo
}

var Help = &HelpBuiltin{
	HelpInfo: HelpInfo{
		Name:        "help",
		Description: "Display help information about ggit",
		UsageLine:   "[command]",
		ManPage:     "TODO",
	},
}

func (b *HelpBuiltin) Execute(p *Params, args []string) {
	if len(args) < 1 {
		b.Usage(p.Werr)
		return
	}
	name := args[0]
	cmd, ok := Get(name)
	if ok {
		cmd.Info().Usage(p.Wout)
	} else {
		fmt.Fprintf(p.Werr, "No manual entry for poop %s\n", name)
	}
}
