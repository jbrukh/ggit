package builtin

import (
	"errors"
	"fmt"
	"github.com/jbrukh/ggit/api"
	"io"
	"os"
)

var (
	fType, fPrint, fSize bool
)

var catFileBuiltin = &Builtin{
	Execute:     catFile,
	Name:        "cat-file",
	Description: "Provide content or type and size information for repository objects",
	UsageLine:   "(-t|-s|-p) <object>",
	ManPage:     "TODO",
}

func init() {
	catFileBuiltin.FlagSet.BoolVar(&fType, "t", false, "show object type")
	catFileBuiltin.FlagSet.BoolVar(&fPrint, "p", false, "pretty-print object's contents")
	catFileBuiltin.FlagSet.BoolVar(&fSize, "s", false, "show object size")

	// add to command list
	Add(catFileBuiltin)
}

func catFile(b *Builtin, args []string, repo api.Repository, w io.Writer) {
	if len(args) != 1 {
		b.Usage(w)
		return
	}
	id := args[0]
	oid, err := api.NewObjectIdFromString(id)
	if err != nil {
		// TODO
		fmt.Fprintln(w, "unknown object")
		return
	}

	switch {
	case fPrint:
		err = doPrint(repo, oid)
	case fType:
		err = doType(repo, oid)
	case fSize:
		err = doSize(repo, oid)
	default:
		panic("should not get here")
	}

	if err != nil {
		fmt.Fprintln(w, err.Error())
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
