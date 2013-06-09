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
	"github.com/jbrukh/ggit/api/diff"
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
		Description: "Describe the difference between two files or trees", //TODO: between commits
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
	var ot objects.ObjectType
	if oah, obh := oa.Header().Type(), ob.Header().Type(); oah != obh {
		fmt.Printf("objects are not the same type: \n%s [%s]\n%s [%s]\n", args[0], oah.String(), args[1], obh.String())
		return
	} else {
		ot = oah
	}
	switch ot {
	case objects.ObjectBlob:
		oab, _ := oa.(*objects.Blob)
		obb, _ := ob.(*objects.Blob)
		d := diff.BlobDiff(oab, obb)
		gdiff.Unified().Print(d, os.Stdout)
	case objects.ObjectTree:
		/*oat, _ := oa.(*objects.Tree)
		obt, _ := oa.(*objects.Tree)
		td := diff.TreeDiff(oat, obt)*/
		//TODO
	default:
		fmt.Printf("objects are of type %s; only blobs (files) and trees are currently supported", ot)
	}
}
