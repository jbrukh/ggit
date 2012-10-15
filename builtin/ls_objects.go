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
)

func init() {
	// add to command list
	Add(LsObjects)
}

type LsObjectsBuiltin struct {
	HelpInfo
}

var LsObjects = &LsObjectsBuiltin{
	HelpInfo: HelpInfo{
		Name:        "ls-objects",
		Description: "Provide a debug dump of all loose object ids", //TODO all object ids
		UsageLine:   "",
		ManPage:     "TODO",
	},
}

func (b *LsObjectsBuiltin) Execute(p *Params, args []string) {
	oids, e := p.Repo.ObjectIds()
	if e != nil {
		fmt.Fprintf(p.Werr, "Error: %s\n", e.Error())
		return
	}
	for i := range oids {
		oid := oids[i]
		fmt.Fprintln(p.Wout, oid.String())
	}
}
