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
	"github.com/jbrukh/ggit/api/format"
	"github.com/jbrukh/ggit/api/objects"
	"github.com/mikebosw/gdiff"
	"os"
)

func init() {
	// add to command list
	Add(Diff)
}

type DiffBuiltin struct {
	HelpInfo
}

var Diff = &DiffBuiltin{
	HelpInfo: HelpInfo{
		Name:        "diff",
		Description: "Describe the difference between two files", //TODO: between commits
		UsageLine:   "",
		ManPage:     "TODO",
	},
}

func (db *DiffBuiltin) Execute(p *Params, args []string) {
	if len(args) < 2 {
		fmt.Printf("expected two arguments; got %d\n", len(args))
		return
	}
	var e error
	var oa, ob objects.Object
	if oa, e = api.ObjectFromRevision(p.Repo, args[0]); e != nil {
		fmt.Printf(e.Error())
		return
	}
	if ob, e = api.ObjectFromRevision(p.Repo, args[1]); e != nil {
		fmt.Printf(e.Error())
		return
	}
	fa := format.NewStrFormat()
	fa.Object(oa)
	a := fa.String()
	fb := format.NewStrFormat()
	fb.Object(ob)
	b := fb.String()
	differ := gdiff.MyersDiffer()
	d := differ.Diff(a, b, gdiff.LINE_SPLIT)
	d.Unified(os.Stdout)
}
