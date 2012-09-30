package api

import (
	"fmt"
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

type PackedRefs []*PackedRef

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

func (p *refParser) ParseRefFile() (oid *ObjectId, err error) {
	err = safeParse(func() {
		oid = p.ParseObjectId()
		p.ConsumeByte(LF)
	})
	return oid, err
}
