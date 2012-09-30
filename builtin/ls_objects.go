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
		UsageLine:   "ls-objects",
		ManPage:     "TODO",
	},
}

func (b *LsObjectsBuiltin) Info() *HelpInfo {
	return &b.HelpInfo
}

func (b *LsObjectsBuiltin) Execute(p *Params, args []string) {
	oids, e := p.Repo.ObjectIds()
	if e != nil {
		println("Error:", e.Error())
		return
	}
	for i := range oids {
		oid := oids[i]
		fmt.Fprintln(p.Wout, oid.String())
	}
}
