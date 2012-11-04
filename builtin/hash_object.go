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

type HashObjectBuiltin struct {
	HelpInfo
	flag.FlagSet
	flagType string
}

var HashObject = &HashObjectBuiltin{
	HelpInfo: HelpInfo{
		Name:        "hash-object",
		Description: "Output the SHA1 hash of the specified Git object (blob, tag, commit, or tree)",
		UsageLine:   "<object>",
		ManPage:     "TODO",
	},
}

func init() {
	// add to command list
	Add(HashObject)
}

func (b *HashObjectBuiltin) Execute(p *Params, args []string) {
	HashObject.Parse(args)
	args = HashObject.Args()

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

	if h, err := api.MakeHash(o); err != nil {
		fmt.Fprintf(p.Werr, "fatal: could not get hash for %s: %s", o.ObjectId().String(), err.Error())
		fmt.Fprintln(p.Werr, err.Error())
	} else {
		fmt.Fprintln(p.Wout, api.OidFromHash(h))
	}
}
