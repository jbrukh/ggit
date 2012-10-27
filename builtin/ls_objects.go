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

type LsObjectsBuiltin struct {
	HelpInfo
	flag.FlagSet
	flagLoose  bool
	flagPacked bool
}

var LsObjects = &LsObjectsBuiltin{
	HelpInfo: HelpInfo{
		Name:        "ls-objects",
		Description: "Provide a debug dump of all loose object ids", //TODO all object ids
		UsageLine:   "[--loose] [--packed]",
		ManPage:     "TODO",
	},
}

func init() {
	LsObjects.BoolVar(&LsObjects.flagLoose, "loose", false, "Print loose objects.")
	LsObjects.BoolVar(&LsObjects.flagPacked, "packed", false, "Print packed objects.")
	// add to command list
	Add(LsObjects)
}

func (b *LsObjectsBuiltin) Execute(p *Params, args []string) {
	b.Parse(args)
	args = b.Args()

	var (
		diskRepo *api.DiskRepository
		err      error
	)

	// make sure this is a disk repo
	if b.flagLoose || b.flagPacked {
		println("haha")
		if diskRepo, err = api.AssertDiskRepo(p.Repo); err != nil {
			fmt.Fprintf(p.Werr, err.Error())
			return
		}
	}

	var (
		oids []*api.ObjectId
		e    error
	)
	if b.flagLoose {
		println("loose")
		oids, e = diskRepo.LooseObjectIds()
	} else if b.flagPacked {
		println("packed")
		oids, e = diskRepo.PackedObjectIds()
	} else {
		println("default")
		oids, e = diskRepo.ObjectIds()
	}

	if e != nil {
		fmt.Fprintf(p.Werr, "Error: %s\n", e.Error())
		return
	}
	printAll(p, oids)
}

func printAll(p *Params, oids []*api.ObjectId) {
	for _, oid := range oids {
		fmt.Fprintln(p.Wout, oid.String())
	}
}
