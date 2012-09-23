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

func catFile(b *Builtin, args []string, w io.Writer) {
	if len(args) != 1 {
		b.Usage(w)
		return
	}

	id := args[0]
	var (
		repo *api.DiskRepository
		oid  *api.ObjectId
		err  error
	)

	// get a proper id
	if oid, err = api.NewObjectIdFromString(id); err != nil {
		return
	}

	// TODO: perhaps not open the repo before parsing args?
	if repo, err = api.Open(api.DEFAULT_GIT_DIR); err != nil {
		return
	}
	defer repo.Close()

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

func doPrint(repo *api.DiskRepository, oid *api.ObjectId) error {
	if obj, err := repo.ReadObject(oid); err != nil {
		return errors.New(err.Error())
	} else {
		obj.WriteTo(os.Stdout)
		return err
	}
	return nil
}

func doType(repo *api.DiskRepository, oid *api.ObjectId) (err error) {
	var obj api.Object
	if obj, err = repo.ReadObject(oid); err != nil {
		return err
	}
	fmt.Println(obj.Type())
	return
}

// commenting until I figure out what size means in this context
func doSize(repo *api.DiskRepository, oid *api.ObjectId) (err error) {
	var obj api.Object
	if obj, err = repo.ReadObject(oid); err != nil {
		return err
	}
	fmt.Println(obj.Size())
	return
}
