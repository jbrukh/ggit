package builtin

import (
	"errors"
	"fmt"
	"github.com/jbrukh/ggit/api"
	"os"
)

func init() {
	fs := CatFile.HelpInfo.FlagSet
	fs.BoolVar(&CatFile.ShowType, "t", false, "show object type")
	fs.BoolVar(&CatFile.ShowSize, "p", false, "pretty-print object's contents")
	fs.BoolVar(&CatFile.PrettyPrint, "s", false, "show object size")

	// add to command list
	Add(CatFile)
}

type CatFileBuiltin struct {
	HelpInfo
	ShowType    bool
	ShowSize    bool
	PrettyPrint bool
}

var CatFile = &CatFileBuiltin{
	HelpInfo: HelpInfo{
		Name:        "cat-file",
		Description: "Provide content or type and size information for repository objects",
		UsageLine:   "(-t|-s|-p) <object>",
		ManPage:     "TODO",
	},
}

func (b *CatFileBuiltin) Info() *HelpInfo {
	return &b.HelpInfo
}

func (b *CatFileBuiltin) Execute(p *Params, args []string) {
	if len(args) != 1 {
		b.HelpInfo.Usage(p.Werr)
		return
	}
	id := args[0]
	oid, err := api.NewObjectIdFromString(id)
	if err != nil {
		b.HelpInfo.Usage(p.Werr)
		return
	}

	switch {
	case b.PrettyPrint:
		err = b.doPrint(p.Repo, oid)
	case b.ShowType:
		err = b.doType(p.Repo, oid)
	case b.ShowSize:
		err = b.doSize(p.Repo, oid)
	default:
		panic("should not get here")
	}

	if err != nil {
		fmt.Fprintln(p.Werr, err.Error())
	}
}

func (b *CatFileBuiltin) doPrint(repo api.Repository, oid *api.ObjectId) error {
	if obj, err := repo.ReadObject(oid); err != nil {
		return errors.New(err.Error())
	} else {
		obj.WriteTo(os.Stdout)
		return err
	}
	return nil
}

func (b *CatFileBuiltin) doType(repo api.Repository, oid *api.ObjectId) (err error) {
	var obj api.Object
	if obj, err = repo.ReadObject(oid); err != nil {
		return err
	}
	fmt.Println(obj.Type())
	return
}

// commenting until I figure out what size means in this context
func (b *CatFileBuiltin) doSize(repo api.Repository, oid *api.ObjectId) (err error) {
	var obj api.Object
	if obj, err = repo.ReadObject(oid); err != nil {
		return err
	}
	fmt.Println(obj.Size())
	return
}
