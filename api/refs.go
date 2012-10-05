package api

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	regexpCaret    *regexp.Regexp
	regexpTilde    *regexp.Regexp
	regexpHex      *regexp.Regexp
	regexpShortHex *regexp.Regexp
)

func init() {
	regexpCaret, _ = regexp.Compile("^[^\\^]+\\^+$")
	regexpTilde, _ = regexp.Compile("^.+~[1-9][0-9]*$")
	regexpHex, _ = regexp.Compile("^[0-9a-f]{40}$") // TODO: replace with const
	regexpShortHex, _ = regexp.Compile("^[0-9a-f]{3,39}$")
}

// ================================================================= //
// REF OBJECTS
// ================================================================= //

// Ref is a representation of a ggit reference. A ref is a nice
// name for an ObjectId. 
type Ref interface {
	ObjectId() *ObjectId
	Target() *ObjectId
	Name() string
}

// sort interface for sorting refs
type refByName []Ref

func (s refByName) Len() int           { return len(s) }
func (s refByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s refByName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }

type ref struct {
	oid       *ObjectId
	name      string
	targetOid *ObjectId // if tag, this is the commit the tag points to
}

func (r *ref) ObjectId() *ObjectId {
	return r.oid
}

func (r *ref) Name() string {
	return r.name
}

func (r *ref) Target() *ObjectId {
	return r.targetOid
}

// ================================================================= //
// TAG FORMATTING
// ================================================================= //

func (f *Format) Ref(r Ref) (int, error) {
	return fmt.Fprintf(f.Writer, "%s %s", r.ObjectId(), r.Name())
}

func (f *Format) Deref(r Ref) (int, error) {
	return fmt.Fprintf(f.Writer, "%s %s^{}", r.Target(), r.Name())
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
				targetOid := p.ParseObjectId()
				p.ConsumeByte(LF)

				if l := len(r); l > 0 {
					r[l-1].(*ref).targetOid = targetOid
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

// ================================================================= //
// REF RESOLUTION
// ================================================================= //

func ResolveRef(repo Repository, refstr string) (*ObjectId, error) {
	if regexpCaret.MatchString(refstr) {
		oid, err := ResolveRef(repo, trimLast(refstr))
		if err != nil {
			return nil, err
		}
		obj, e := repo.ReadObject(oid)
		if e != nil {
			return nil, e
		}
		t := obj.Type()
		if t == ObjectCommit {
			oid = obj.(*Commit).FirstParent()
			if oid == nil {
				return nil, errors.New("no parent")
			}
		} else if t == ObjectTag {
			oid = obj.(*Tag).Object()
		}
		return oid, nil
	} else if regexpHex.MatchString(refstr) {
		oid, _ := NewObjectIdFromString(refstr)
		return oid, nil
	}
	return nil, errors.New("unknown reference")
}
