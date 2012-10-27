//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package builtin

import (
	"fmt"
	"github.com/jbrukh/ggit/api"
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
		b.WriteUsage(p.Werr)
		return
	}
	rev := args[0]
	o, err := api.ObjectFromRevision(p.Repo, rev)
	if err != nil {
		fmt.Fprintf(p.Wout, "%s\nfatal: ambiguous argument '%s': unknown revision or path not in the working tree.\n", rev, rev)
		fmt.Fprintln(p.Werr, err.Error())
		return
	}

	fmt.Fprintln(p.Wout, o.ObjectId())
}
