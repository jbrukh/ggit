package builtin

import (
	"flag"
	"fmt"
	"github.com/jbrukh/ggit/api"
	"strings"
)

type ShowRefBuiltin struct {
	HelpInfo
	flags     flag.FlagSet
	flagQuiet bool
	flagWhich bool
	flagHeads bool
}

var ShowRef = &ShowRefBuiltin{
	HelpInfo: HelpInfo{
		Name:        "show-ref",
		Description: "List references in a local repository",
		UsageLine:   "[--which] [<pattern>]",
		ManPage:     "TODO",
	},
}

//var flags flag.FlagSet

func init() {
	ShowRef.flags.BoolVar(&ShowRef.flagQuiet, "q", false, "Do not print any results to stdout.")
	ShowRef.flags.BoolVar(&ShowRef.flagWhich, "which", false, "Show which refs are loose and which are packed.")
	ShowRef.flags.BoolVar(&ShowRef.flagHeads, "heads", false, "Show only heads.")
	ShowRef.flags.Usage = func() {}

	// add to command list
	Add(ShowRef)
}

func (b *ShowRefBuiltin) Execute(p *Params, args []string) {
	b.flags.Parse(args)
	args = b.flags.Args()

	switch {
	case b.flagWhich:
		b.Which(p)

	case b.flagHeads:
		b.Heads(p)

	default:

	}

	//fmt.Println("getting ", args)
	if len(args) > 0 {
		// we want to see a particular ref
		pattern := args[0]
		b.WithPattern(p, pattern)
	} else {
		if refs, e := p.Repo.Refs(); e != nil {
			fmt.Fprintln(p.Werr, e.Error())
			return
		} else {
			for _, v := range refs {
				fmt.Fprintln(p.Wout, v.String())
			}
		}
	}
}

func filterRefs(refs []api.Ref, f func(api.Ref) bool) []api.Ref {
	r := make([]api.Ref, 0, len(refs))
	for _, v := range refs {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

// matchRefs performs the matching of a partial ref with a full (or longer)
// ref. Matching occurs from the end and matches on completed parts of the
// name. So for instance, refs/heads/master and master would match, but "ter"
// would not match the former.
func matchRefs(full, partial string) bool {
	const SL = "/"
	f, p := strings.Split(full, SL), strings.Split(partial, SL)
	i, j := len(f), len(p)
	if i == 0 || j == 0 || i < j { // partial must be shorter
		return false
	}
	for j > 0 {
		i--
		j--
		if f[i] != p[j] {
			return false
		}
	}
	return true
}

func refsWithFilter(p *Params, f func(ref api.Ref) bool) {
	refs, e := p.Repo.Refs()
	if e != nil {
		fmt.Fprintln(p.Werr, e.Error())
	}
	filtered := filterRefs(refs, f)
	for _, v := range filtered {
		fmt.Fprintln(p.Wout, v.String())
	}
}

func (b *ShowRefBuiltin) WithPattern(p *Params, pattern string) {
	refsWithFilter(p, func(ref api.Ref) bool {
		return matchRefs(ref.Name(), pattern)
	})
}

func (b *ShowRefBuiltin) Heads(p *Params) {
	refsWithFilter(p, func(ref api.Ref) bool {
		return strings.HasPrefix(ref.Name(), "refs/heads")
	})
}

func (b *ShowRefBuiltin) Which(p *Params) {
	repo := p.Repo.(*api.DiskRepository)

	fmt.Fprintln(p.Wout, "Loose refs:")
	refs, e := repo.LooseRefs()
	if e != nil {
		fmt.Fprint(p.Werr, e.Error())
		return
	}
	for _, v := range refs {
		fmt.Fprintln(p.Wout, v.String())
	}

	fmt.Fprintln(p.Wout, "\nPacked refs:")
	prefs, e := repo.PackedRefs()
	if e != nil {
		fmt.Fprint(p.Werr, e.Error())
		return
	}
	for _, v := range prefs {
		fmt.Fprintln(p.Wout, v.String())
	}
}
