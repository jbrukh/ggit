package api

import (
	"fmt"
	"strings"
)

// ================================================================= //
// REF OBJECTS
// ================================================================= //

// Ref is a representation of a ggit reference. A ref is a nice
// name for an ObjectId. 
type Ref interface {
	ObjectId() *ObjectId
	Name() string
	String() string
}

// sort interface for sorting refs
type refByName []Ref

func (s refByName) Len() int           { return len(s) }
func (s refByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s refByName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }

type NamedRef struct {
	oid  *ObjectId
	name string
}

func (r *NamedRef) ObjectId() *ObjectId {
	return r.oid
}

func (r *NamedRef) Name() string {
	return r.name
}

func (p *NamedRef) String() string {
	const format = "%s %s"
	return fmt.Sprintf(format, p.oid, p.name)
}

type PackedRef struct {
	NamedRef

	// if this is an annotated tag
	// we may have the pointed to commit here
	// as an optimization
	cid *ObjectId
}

func (p *PackedRef) String() string {
	return p.NamedRef.String()
}

// TODO: do we need this?
type PackedRefs []*PackedRef

// ================================================================= //
// REF FILTERING
// ================================================================= //

type RefFilter func(Ref) bool

func FilterRefs(refs []Ref, filters []RefFilter) []Ref {
	r := make([]Ref, 0, len(refs))
	for _, v := range refs {
		keep := true
		for _, f := range filters {
			if !f(v) {
				keep = false
				break
			}
		}
		if keep {
			r = append(r, v)
		}
	}
	return r
}

func RefFilterOr(filters []RefFilter) RefFilter {
	return func(ref Ref) bool {
		for _, f := range filters {
			if f(ref) {
				return true
			}
		}
		return false
	}
}

func RefFilterAnd(filters []RefFilter) RefFilter {
	return func(ref Ref) bool {
		for _, f := range filters {
			if !f(ref) {
				return false
			}
		}
		return true
	}
}

func RefFilterPattern(pattern string) RefFilter {
	return func(ref Ref) bool {
		return matchRefs(ref.Name(), pattern)
	}
}

func refFilterPrefix(prefix string) RefFilter {
	return func(ref Ref) bool {
		return strings.HasPrefix(ref.Name(), prefix)
	}
}

func RefFilterPrefix(prefix ...string) RefFilter {
	f := make([]RefFilter, 0, len(prefix))
	for _, p := range prefix {
		f = append(f, refFilterPrefix(p))
	}
	return RefFilterOr(f)
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

// ================================================================= //
// REF PARSING
// ================================================================= //

func (p *refParser) ParsePackedRefs() (PackedRefs, error) {
	r := make(PackedRefs, 0)
	err := safeParse(func() {
		for !p.EOF() {
			c := p.PeekByte()
			switch c {
			case '#':
				// if this is the first line, then it should be a comment
				// that says '# pack-refs with: <extention>' and <extention>
				// is exactly one of the items in this set: { 'peeled' }.
				// currently, we are just ignoring all comments.
				p.ReadString(LF)
			case '^':
				// this means the previous line is an annotated tag and the the current
				// line contains the commit that tag points to
				p.ConsumeByte('^')
				cid := p.ParseObjectId()
				p.ConsumeByte(LF)

				if l := len(r); l > 0 {
					r[l-1].cid = cid
				}
			default:
				re := new(NamedRef)
				re.oid = p.ParseObjectId()
				p.ConsumeByte(SP)
				re.name = p.ReadString(LF)

				r = append(r, &PackedRef{
					NamedRef: *re,
				})
			}
		}
	})
	return r, err
}
