//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import (
	"bufio"
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

	// ObjectId returns the object id that this ref references
	// provided this ref is not symbolic, and otherwise panics.
	ObjectId() *ObjectId

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
// ERRORS
// ================================================================= //

type (
	noSuchRef    error
	ambiguousRef error
)

func noSuchRefErrf(ref string) noSuchRef {
	return noSuchRef(fmt.Errorf("no such ref: %s", ref))
}

func ambiguousRefErrf(ref string) noSuchRef {
	return ambiguousRef(fmt.Errorf("ambiguous ref: %s", ref))
}

func IsNoSuchRef(e error) bool {
	switch t := e.(type) {
	case noSuchRef:
		e = t // must use t
		return true
	}
	return false
}

func IsAmbiguousRef(e error) bool {
	switch t := e.(type) {
	case ambiguousRef:
		e = t // must use t
		return true
	}
	return false
}

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

func (r *ref) ObjectId() *ObjectId {
	symbolic, oid := r.Target()
	if !symbolic {
		return oid.(*ObjectId)
	}
	panic("cannot return oid: this is a symbolic ref")
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
	var r []Ref
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

// refParser implements functions for parsing refs and packed refs.
type refParser struct {
	objectIdParser
	name string // the name of the ref file
}

// newRefParser creates a new Ref parser 
func newRefParser(buf *bufio.Reader, name string) *refParser {
	return &refParser{
		objectIdParser: objectIdParser{
			dataParser{
				buf: buf,
			},
		},
		name: name,
	}
}

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
				commit := p.ParseOid()
				p.ConsumeByte(LF)

				if l := len(r); l > 0 {
					r[l-1].(*ref).commit = commit
				}
			default:
				re := new(ref)
				re.oid = p.ParseOid()
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
			oid := p.ParseOid()
			p.ConsumeByte(LF)
			r = &ref{name: p.name, oid: oid}
		}
	})
	return r, err
}

// ================================================================= //
// OPERATIONS
// ================================================================= //

var refSearchPath = []string{
	"%s",
	"refs/%s",
	"refs/tags/%s",
	"refs/heads/%s",
	"refs/remotes/%s",
	"refs/remotes/%s/HEAD",
}

// RefFromSpec delivers a ref from the repository that is in
// the form of a Ref object. The ref may be symbolic or peeled.
// 
// This method resolves "shorthand" refs (e.g. when one writes 
// "master" for "refs/head/master". It follows the rules specified 
// here:
//
//    http://www.kernel.org/pub/software/scm/git/docs/gitrevisions.html
//
// under "Specifying Revisions/<refname>". If the ref does not exist
// then you can check the returned error with api.IsNoSuchRef.
// TODO: what about ambiguous refs?
func RefFromSpec(repo Repository, spec string) (ref Ref, err error) {
	for _, prefix := range refSearchPath {
		refPath := fmt.Sprintf(prefix, spec)
		if ref, err = repo.Ref(refPath); err == nil {
			return ref, nil
		} else if !IsNoSuchRef(err) {
			return nil, err // something went wrong
		}
	}
	return nil, err // no such ref
}

// PeeledRefFromSpec takes the same arguments as RefFromSpec, but
// peels the ref before sending it back.
func PeeledRefFromSpec(repo Repository, spec string) (ref Ref, err error) {
	if ref, err = RefFromSpec(repo, spec); err != nil {
		return nil, err
	}
	return PeelRef(repo, ref)
}

// PeelRef resolves the final target oid of the ref and returns
// a peeled ref for this target. It examines the target of the
// given ref and if the target is symbolic, it is followed and
// resolved. This process repeats as many times as necessary to
// obtain a peeled ref.
func PeelRef(repo Repository, r Ref) (Ref, error) {
	var (
		err      error
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
	return r, nil
}

func expandHeadRef(short string) string {
	return "refs/heads/" + short
}

func expandTagRef(short string) string {
	return "refs/tags/" + short
}
