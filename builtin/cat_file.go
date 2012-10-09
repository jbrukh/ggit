package builtin

import (
	"errors"
	"flag"
	"fmt"
	"github.com/jbrukh/ggit/api"
)

type CatFileBuiltin struct {
	HelpInfo
	flags           flag.FlagSet
	flagShowType    bool
	flagShowSize    bool
	flagPrettyPrint bool
}

var CatFile = &CatFileBuiltin{
	HelpInfo: HelpInfo{
		Name:        "cat-file",
		Description: "Provide content or type and size information for repository objects",
		UsageLine:   "(-t|-s|-p) <object>",
		ManPage:     "TODO",
	},
}

func init() {
	CatFile.flags.BoolVar(&CatFile.flagShowType, "t", false, "show object type")
	CatFile.flags.BoolVar(&CatFile.flagPrettyPrint, "p", false, "pretty-print object's contents")
	CatFile.flags.BoolVar(&CatFile.flagShowSize, "s", false, "show object size")

	// add to command list
	Add(CatFile)
}

func (b *CatFileBuiltin) Execute(p *Params, args []string) {
	CatFile.flags.Parse(args)
	args = CatFile.flags.Args()

	if len(args) != 1 {
		b.HelpInfo.Usage(p.Werr)
		return
	}
	id := args[0]
	oid, err := api.OidFromString(id)
	if err != nil {
		b.HelpInfo.Usage(p.Werr)
		return
	}

	switch {
	case b.flagPrettyPrint:
		err = b.PrettyPrint(p, oid)
	case b.flagShowType:
		err = b.ShowType(p, oid)
	case b.flagShowSize:
		err = b.ShowSize(p, oid)
	default:
		b.HelpInfo.Usage(p.Werr)
		return
	}

	if err != nil {
		fmt.Fprintln(p.Werr, err.Error())
	}
}

func (b *CatFileBuiltin) PrettyPrint(p *Params, oid *api.ObjectId) error {
	if o, err := api.ObjectFromOid(p.Repo, oid); err != nil {
		return errors.New(err.Error())
	} else {
		f := api.Format{p.Wout}
		f.Object(o)
		return err
	}
	return nil
}

func (b *CatFileBuiltin) ShowType(p *Params, oid *api.ObjectId) (err error) {
	var o api.Object
	if o, err = api.ObjectFromOid(p.Repo, oid); err != nil {
		return err
	}
	fmt.Fprintln(p.Wout, o.Type())
	return
}

func (b *CatFileBuiltin) ShowSize(p *Params, oid *api.ObjectId) (err error) {
	var o api.Object
	if o, err = api.ObjectFromOid(p.Repo, oid); err != nil {
		return err
	}
	fmt.Fprintln(p.Wout, o.Size())
	return
}
