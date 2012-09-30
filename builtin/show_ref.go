package builtin

import (
	"flag"
	"fmt"
)

type ShowRefBuiltin struct {
	HelpInfo
	flags     flag.FlagSet
	flagQuiet bool
}

var ShowRef = &ShowRefBuiltin{
	HelpInfo: HelpInfo{
		Name:        "show-ref",
		Description: "List references in a local repository",
		UsageLine:   "show-ref [<refPath>]",
		ManPage:     "TODO",
	},
}

//var flags flag.FlagSet

func init() {
	ShowRef.flags.BoolVar(&ShowRef.flagQuiet, "q", false, "Do not print any results to stdout.")
	// add to command list
	Add(ShowRef)
}

func (b *ShowRefBuiltin) Info() *HelpInfo {
	return &b.HelpInfo
}

func (b *ShowRefBuiltin) Execute(p *Params, args []string) {
	ShowRef.flags.Parse(args)
	args = ShowRef.flags.Args()

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
		fmt.Fprintln(p.Wout, "File refs:\n")
		if refs, e := p.Repo.ReadRefs(); e != nil {
			fmt.Fprintln(p.Werr, e.Error())
			return
		} else {
			for _, v := range refs {
				fmt.Fprintln(p.Wout, v.String())
			}
		}

		fmt.Fprintln(p.Wout, "\nPacked refs:\n")
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
