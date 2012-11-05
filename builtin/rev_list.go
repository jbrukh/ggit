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
	// flagShowType    bool
	// flagShowSize    bool
	// flagPrettyPrint bool
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
	// CatFile.BoolVar(&CatFile.flagShowType, "t", false, "show object type")
	// CatFile.BoolVar(&CatFile.flagPrettyPrint, "p", false, "pretty-print object's contents")
	// CatFile.BoolVar(&CatFile.flagShowSize, "s", false, "show object size")

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
	api.RevWalkFromRevision(p.Repo, rev, api.RevPrinter)
}
