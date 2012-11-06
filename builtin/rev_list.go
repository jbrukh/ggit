//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package builtin

import (
	"flag"
	//"fmt"
	"github.com/jbrukh/ggit/api"
)

type RevListBuiltin struct {
	HelpInfo
	flag.FlagSet
	flagNoMerges bool
}

var RevList = &RevListBuiltin{
	HelpInfo: HelpInfo{
		Name:        "rev-list",
		Description: "Lists commit objects in reverse chronological order",
		UsageLine:   "<revision>",
		ManPage:     "TODO",
	},
}

func init() {
	RevList.BoolVar(&RevList.flagNoMerges, "no-merges", false, "Do not print commits with more than one parent.")

	// add to command list
	Add(RevList)
}

func (b *RevListBuiltin) Execute(p *Params, args []string) {
	RevList.Parse(args)
	args = RevList.Args()

	if len(args) < 1 {
		b.WriteUsage(p.Werr)
		return
	}

	rev := args[0]
	opts := new(api.RevWalkOptions)
	if b.flagNoMerges {
		opts.NoMerges = true
	}
	api.RevWalkDateOrder(p.Repo, rev, opts, api.RevPrinter)
}
