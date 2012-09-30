package builtin

import (
	"fmt"
	"github.com/jbrukh/ggit/api"
	"io"
)

var lsObjectsBuiltin = &Builtin{
	Execute:     lsObjects,
	Name:        "ls-objects",
	Description: "Provide a debug dump of all loose object ids", //TODO all object ids
	UsageLine:   "ls-objects",
	ManPage:     "TODO",
}

func init() {
	// add to command list
	Add(lsObjectsBuiltin)
}

func lsObjects(b *Builtin, args []string, repo api.Repository, w io.Writer) {
	oids, e := repo.ObjectIds()
	if e != nil {
		println("Error:", e.Error())
		return
	}
	for i := range oids {
		oid := oids[i]
		fmt.Fprintln(w, oid.String())
	}
}
