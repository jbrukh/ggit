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
	"github.com/jbrukh/ggit/api/format"
	"github.com/jbrukh/ggit/api/objects"
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
	flag.FlagSet
	flagQuiet bool
	flagWhich bool
	flagHeads bool
	flagTags  bool
	flagHelp  bool
	flagDeref bool
	flagHead  bool
}

var ShowRef = &ShowRefBuiltin{
	HelpInfo: HelpInfo{
		Name:        "show-ref",
		Description: "List references in a local repository",
		UsageLine:   "[--which] [-d] [--head] [--heads] [--tags] [<pattern>]",
		ManPage:     "TODO",
	},
}

// ================================================================= //
// SHOW-REF FLAGS
// ================================================================= //

func init() {
	ShowRef.BoolVar(&ShowRef.flagQuiet, "q", false, "Do not print any results to stdout.")
	ShowRef.BoolVar(&ShowRef.flagWhich, "which", false, "Show which refs are loose and which are packed.")
	ShowRef.BoolVar(&ShowRef.flagHeads, "heads", false, "Show only heads.")
	ShowRef.BoolVar(&ShowRef.flagTags, "tags", false, "Show only tags.")
	ShowRef.BoolVar(&ShowRef.flagHelp, "help", false, "Show help.")
	ShowRef.BoolVar(&ShowRef.flagDeref, "d", false, "Dereference tags into object IDs as well.")
	ShowRef.BoolVar(&ShowRef.flagHead, "head", false, "Show the head reference.")

	ShowRef.Usage = func() {}

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
	b.Parse(args)
	args = b.Args()

	if b.flagWhich {
		b.Which(p)
		return
	}

	if b.flagHelp {
		b.WriteUsage(p.Wout)
		return
	}

	var f []api.Filter
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
		return
	}
	f := api.FilterAnd(filters...)
	filtered := api.FilterRefs(refs, f)

	if b.flagHead {
		if r, err := api.PeeledRefFromSpec(p.Repo, "HEAD"); err == nil {
			filtered = append([]objects.Ref{r}, filtered...)
		}
	}
	// formatter
	fmtr := format.Format{p.Wout}

	if b.flagQuiet {
		return
	}

	if b.flagDeref {
		for _, r := range filtered {
			fmtr.Ref(r)
			fmtr.Lf()
			if r.Commit() != nil {
				fmtr.Deref(r)
				fmtr.Lf()
			} else {
				_, oid := r.Target() // better not be symbolic
				o, err := p.Repo.ObjectFromOid(oid.(*objects.ObjectId))
				if err == nil {
					if o.Header().Type() == objects.ObjectTag {
						tag := o.(*objects.Tag)
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

/*
TODO: remove this method, it is mainly for debugging
*/
func (b *ShowRefBuiltin) Which(p *Params) {
	repo := p.Repo.(*api.DiskRepository)
	fmtr := format.Format{p.Wout}

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
