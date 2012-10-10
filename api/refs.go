package api

import (
	"fmt"
	"strings"
)

const (
	markerRef = "ref:"
)

// ================================================================= //
// REF OBJECTS
// ================================================================= //

// Ref is a representation of a ggit reference. A ref is a nice
// name for an ObjectId. More precisely, a ref is a path relative
// to the git directory (without duplicate path separators, ".", or
// "..").
type Ref interface {

	// Name returns the string name of this ref. This is
	// a simple path relative to the git directory, which
	// may or may not be HEAD, MERGE_HEAD, etc.
	Name() string

	// Target returns the target reference, whether an oid
	// or another string ref. If the ref is symbolic then
	// "symbolic" is true.
	Target() (symbolic bool, o interface{})

	// If this ref is a tag, then this field may contain
	// the target commit of the tag, if such an optimization
	// is available. Otherwise, this field is nil.
	Commit() *ObjectId
}

// sort interface for sorting refs
type refByName []Ref

func (s refByName) Len() int           { return len(s) }
func (s refByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s refByName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }

// ================================================================= //
// REF IMPLEMENTATION
// ================================================================= //

type ref struct {
	name   string
	oid    *ObjectId
	spec   string
	commit *ObjectId // if tag, this is the commit the tag points to
}

func (r *ref) Target() (bool, interface{}) {
	if r.oid != nil {
		return false, r.oid
	}
	if r.spec != "" {
		return true, r.spec
	}
	panic("does not have an object reference")
}

func (r *ref) Name() string {
	return r.name
}

func (r *ref) Commit() *ObjectId {
	return r.commit
}

// ================================================================= //
// REF FORMATTING
// ================================================================= //

func (f *Format) Ref(r Ref) (int, error) {
	_, rf := r.Target() // symbolic or oid
	return fmt.Fprintf(f.Writer, "%s %s", rf, r.Name())
}

// TODO: come up with a better name for this
func (f *Format) Deref(r Ref) (int, error) {
	return fmt.Fprintf(f.Writer, "%s %s^{}", r.Commit(), r.Name())
}

// ================================================================= //
// REF FILTERING
// ================================================================= //

func FilterRefs(refs []Ref, f Filter) []Ref {
	r := make([]Ref, 0)
	for _, v := range refs {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

func FilterRefPattern(pattern string) Filter {
	return func(ref interface{}) bool {
		return matchRefs(ref.(Ref).Name(), pattern)
	}
}

func FilterRefPrefix(prefix string) Filter {
	return func(ref interface{}) bool {
		return strings.HasPrefix(ref.(Ref).Name(), prefix)
	}
}

// matchRefs performs the matching of a partial ref with a full (or longer)
// ref. Matching occurs from the end and matches on completed parts of the
// name. So for instance, refs/heads/master and master would match, but "ter"
// would not match the former.
func matchRefs(full, partial string) bool {
	const SL = "/"
	if full == "" || partial == "" {
		return false
	}

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

func (p *refParser) ParsePackedRefs() ([]Ref, error) {
	r := make([]Ref, 0)
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
				commit := p.ParseObjectId()
				p.ConsumeByte(LF)

				if l := len(r); l > 0 {
					r[l-1].(*ref).commit = commit
				}
			default:
				re := new(ref)
				re.oid = p.ParseObjectId()
				p.ConsumeByte(SP)
				re.name = p.ReadString(LF)

				r = append(r, re)
			}
		}
	})
	return r, err
}

func (p *refParser) parseRef() (r Ref, err error) {
	err = safeParse(func() {
		// is it a symbolic ref?
		if p.PeekString(len(markerRef)) == markerRef {
			p.ConsumeString(markerRef)
			p.ConsumeByte(SP)
			spec := p.ReadString(LF)
			r = &ref{name: p.name, spec: spec}
		} else {
			oid := p.ParseObjectId()
			p.ConsumeByte(LF)
			r = &ref{name: p.name, oid: oid}
		}
	})
	return r, err
}

// ================================================================= //
// OPERATIONS
// ================================================================= //

// OidFromRef turns a ref specification into the ObjectId that ref
// points to, by peeling away any symbolic refs that might stand
// in its way.
func OidFromRef(repo Repository, spec string) (*ObjectId, error) {
	r, err := OidRefFromRef(repo, spec)
	if err != nil {
		return nil, err
	}
	_, oid := r.Target()
	return oid.(*ObjectId), nil
}

// ObjectFromRef is similar to OidFromRef, except it derefernces the
// target ObjectId into an actual Object.
func ObjectFromRef(repo Repository, spec string) (Object, error) {
	oid, err := OidFromRef(repo, spec)
	if err != nil {
		return nil, err
	}
	return repo.ObjectFromOid(oid)
}

// OidRefFromRef returns a Ref object representing the target of
// the ref, by peeling away any symbolic refs that might stand
// in its way.
func OidRefFromRef(repo Repository, spec string) (Ref, error) {
	r, err := repo.Ref(spec)
	if err != nil {
		return nil, err
	}

	var (
		symbolic bool
		target   interface{}
	)

	// TODO: make a limit
	for {
		symbolic, target = r.Target()
		if symbolic {
			r, err = repo.Ref(target.(string))
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return &ref{name: spec, oid: target.(*ObjectId)}, nil
}
