package builtin

import (
	"flag"
	"fmt"
	"github.com/jbrukh/ggit/api"
)

// ================================================================= //
// SHOW-REF
// ================================================================= //

// ShowRefBuiltin implements a command very similar
// to git-show-ref. This allows one to display the
// refs that are both available in the refs/ directory
// as well as the packed refs file.
type ShowRefBuiltin struct {
	HelpInfo
	flags     flag.FlagSet
	flagQuiet bool
	flagWhich bool
	flagHeads bool
	flagTags  bool
	flagHelp  bool
	flagDeref bool
}

var ShowRef = &ShowRefBuiltin{
	HelpInfo: HelpInfo{
		Name:        "show-ref",
		Description: "List references in a local repository",
		UsageLine:   "[--which] [-d] [--heads] [--tags] [<pattern>]",
		ManPage:     "TODO",
	},
}

// ================================================================= //
// SHOW-REF FLAGS
// ================================================================= //

func init() {
	ShowRef.flags.BoolVar(&ShowRef.flagQuiet, "q", false, "Do not print any results to stdout.")
	ShowRef.flags.BoolVar(&ShowRef.flagWhich, "which", false, "Show which refs are loose and which are packed.")
	ShowRef.flags.BoolVar(&ShowRef.flagHeads, "heads", false, "Show only heads.")
	ShowRef.flags.BoolVar(&ShowRef.flagTags, "tags", false, "Show only tags.")
	ShowRef.flags.BoolVar(&ShowRef.flagHelp, "help", false, "Show help.")
	ShowRef.flags.BoolVar(&ShowRef.flagDeref, "d", false, "Dereference tags into object IDs as well.")
	ShowRef.flags.Usage = func() {}

	// add to command list
	Add(ShowRef)
}

// ================================================================= //
// CONSTANTS
// ================================================================= //

const (
	prefixHeads = "refs/heads"
	prefixTags  = "refs/tags"
)

var HeadsFilter = api.FilterRefPrefix(prefixHeads)
var TagsFilter = api.FilterRefPrefix(prefixTags)

// ================================================================= //
// SHOW-REF FUNCTIONS
// ================================================================= //

func (b *ShowRefBuiltin) Execute(p *Params, args []string) {
	b.flags.Parse(args)
	args = b.flags.Args()

	if b.flagWhich {
		b.Which(p)
		return
	}

	if b.flagHelp {
		b.Usage(p.Wout)
		return
	}

	f := make([]api.Filter, 0)

	if b.flagHeads && b.flagTags {
		f = append(f, api.FilterOr(HeadsFilter, TagsFilter))
	} else if b.flagHeads {
		f = append(f, HeadsFilter)
	} else if b.flagTags {
		f = append(f, TagsFilter)
	}

	if len(args) > 0 {
		pattern := args[0]
		f = append(f, api.FilterRefPattern(pattern))
	}
	b.filterRefs(p, f)
}

func (b *ShowRefBuiltin) filterRefs(p *Params, filters []api.Filter) {
	refs, e := p.Repo.Refs()
	if e != nil {
		fmt.Fprintln(p.Werr, e.Error())
	}
	f := api.FilterAnd(filters...)
	filtered := api.FilterRefs(refs, f)

	// formatter
	fmtr := api.Format{p.Wout}

	if !b.flagQuiet {
		if b.flagDeref {
			for _, r := range filtered {
				fmtr.Ref(r)
				fmtr.Lf()
				if r.Target() != nil {
					fmtr.Deref(r)
					fmtr.Lf()
				} else {
					obj, err := api.ObjectFromOid(p.Repo, r.ObjectId())
					if err == nil {
						if obj.Type() == api.ObjectTag {
							tag := obj.(*api.Tag)
							fmtr.Printf("%s %s^{}\n", tag.Object(), r.Name()) // TODO
						}
					}
				}
			}
		} else { // just do the non-deref case separately, for performance
			for _, r := range filtered {
				fmtr.Ref(r)
				fmtr.Lf()
			}
		}
	}
}

/*
TODO: remove this method, it is mainly for debugging
*/
func (b *ShowRefBuiltin) Which(p *Params) {
	repo := p.Repo.(*api.DiskRepository)
	fmtr := api.Format{p.Wout}

	fmt.Fprintln(p.Wout, "Loose refs:")
	refs, e := repo.LooseRefs()
	if e != nil {
		fmt.Fprint(p.Werr, e.Error())
		return
	}
	for _, v := range refs {
		fmtr.Ref(v)
		fmtr.Lf()
	}

	fmt.Fprintln(p.Wout, "\nPacked refs:")
	prefs, e := repo.PackedRefs()
	if e != nil {
		fmt.Fprint(p.Werr, e.Error())
		return
	}
	for _, v := range prefs {
		fmtr.Ref(v)
		fmtr.Lf()
	}
}
