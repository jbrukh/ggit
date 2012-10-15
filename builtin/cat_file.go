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
	"fmt"
	"github.com/jbrukh/ggit/api"
)

type CatFileBuiltin struct {
	HelpInfo
	flags           flag.FlagSet
	flagShowType    bool
	flagShowSize    bool
	flagPrettyPrint bool
}

var CatFile = &CatFileBuiltin{
	HelpInfo: HelpInfo{
		Name:        "cat-file",
		Description: "Provide content or type and size information for repository objects",
		UsageLine:   "(-t|-s|-p) <object>",
		ManPage:     "TODO",
	},
}

func init() {
	CatFile.flags.BoolVar(&CatFile.flagShowType, "t", false, "show object type")
	CatFile.flags.BoolVar(&CatFile.flagPrettyPrint, "p", false, "pretty-print object's contents")
	CatFile.flags.BoolVar(&CatFile.flagShowSize, "s", false, "show object size")

	// add to command list
	Add(CatFile)
}

func (b *CatFileBuiltin) Execute(p *Params, args []string) {
	CatFile.flags.Parse(args)
	args = CatFile.flags.Args()

	if len(args) != 1 {
		b.HelpInfo.Usage(p.Werr)
		return
	}
	id := args[0]
	o, err := api.ObjectFromShortOid(p.Repo, id)
	if err != nil {
		fmt.Fprintln(p.Werr, err)
		return
	}

	switch {
	case b.flagPrettyPrint:
		f := api.Format{p.Wout}
		f.Object(o)
	case b.flagShowType:
		fmt.Fprintln(p.Wout, o.Header().Type())
	case b.flagShowSize:
		fmt.Fprintln(p.Wout, o.Header().Size())
	default:
		b.HelpInfo.Usage(p.Werr)
		return
	}
}
