package builtin

import (
	"fmt"
	"github.com/jbrukh/ggit/api"
)

func init() {
	// add to command list
	Add(VerifyPackedObjects)
}

type VerifyPackedObjectsBuiltin struct {
	HelpInfo
}

var VerifyPackedObjects = &VerifyPackedObjectsBuiltin{
	HelpInfo: HelpInfo{
		Name:        "verify-packed-objects",
		Description: "Debug command for listing packed objects by id",
		UsageLine:   "",
		ManPage:     "TODO",
	},
}

func (b *VerifyPackedObjectsBuiltin) Execute(p *Params, args []string) {
	var repo *api.DiskRepository
	switch t := p.Repo.(type) {
	case *api.DiskRepository:
		repo = t
	default:
		fmt.Fprintf(p.Werr, "verify-packed-objects applies only to DiskRepository; found: %s", t)
		return
	}
	objects, e := repo.PackedObjects()
	if e != nil {
		fmt.Fprintf(p.Werr, "Error: %s\n", e.Error())
		return
	}
	for i := range objects {
		object := objects[i]
		fmt.Fprintln(p.Wout, object.ObjectId(), object.Header().Type())
	}
}
