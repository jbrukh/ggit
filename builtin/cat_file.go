package builtin

import (
	"errors"
	"fmt"
	"github.com/jbrukh/ggit/api"
	"os"
)

var (
	fType, fPrint, fSize bool
)

var catFileInfo = &Builtin{
	Execute:     catFile,
	Name:        "cat-file",
	Description: "Provide content or type and size information for repository objects",
	UsageLine:   "(-t|-s|-p) <object>",
	ManPage:     "TODO",
}

func init() {
	catFileInfo.FlagSet.BoolVar(&fType, "t", false, "show object type")
	catFileInfo.FlagSet.BoolVar(&fPrint, "p", false, "pretty-print object's contents")
	catFileInfo.FlagSet.BoolVar(&fSize, "s", false, "show object size")

	// add to command list
	Add(catFileInfo)
}

func catFile(p *Params, b *Builtin, args []string) {
	if len(args) != 1 {
		b.Usage(p.Werr)
		return
	}
	id := args[0]
	oid, err := api.NewObjectIdFromString(id)
	if err != nil {
		// TODO
		fmt.Fprintln(p.Werr, "unknown object")
		return
	}

	switch {
	case fPrint:
		err = doPrint(p.Repo, oid)
	case fType:
		err = doType(p.Repo, oid)
	case fSize:
		err = doSize(p.Repo, oid)
	default:
		panic("should not get here")
	}

	if err != nil {
		fmt.Fprintln(p.Werr, err.Error())
	}
}

func doPrint(repo api.Repository, oid *api.ObjectId) error {
	if obj, err := repo.ReadObject(oid); err != nil {
		return errors.New(err.Error())
	} else {
		obj.WriteTo(os.Stdout)
		return err
	}
	return nil
}

func doType(repo api.Repository, oid *api.ObjectId) (err error) {
	var obj api.Object
	if obj, err = repo.ReadObject(oid); err != nil {
		return err
	}
	fmt.Println(obj.Type())
	return
}

// commenting until I figure out what size means in this context
func doSize(repo api.Repository, oid *api.ObjectId) (err error) {
	var obj api.Object
	if obj, err = repo.ReadObject(oid); err != nil {
		return err
	}
	fmt.Println(obj.Size())
	return
}
