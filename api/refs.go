package api

import (
	"fmt"
)

// ================================================================= //
// REF OBJECT
// ================================================================= //

type Ref struct {
	oid *ObjectId
}

// ================================================================= //
// OBJECT PARSER METHODS OBJECT
// ================================================================= //

func (p *objectParser) parseRef() *Ref {
	r := &Ref{
		oid: p.ParseObjectId(),
	}
	p.ConsumeByte(LF)
	return r
}

type RefEntry struct {
	oid     *ObjectId
	refPath string
}

type PackedRef struct {
	RefEntry

	// if this is an annotated tag
	// we may have the pointed to commit here
	// as an optimization
	cid *ObjectId
}

func (p *RefEntry) String() string {
	const format = "%s %s"
	return fmt.Sprintf(format, p.oid, p.refPath)
}

func (p *PackedRef) String() string {
	const format = "%s %s%s"
	if p.cid != nil {
		return fmt.Sprintf(format, p.cid, p.refPath, "^{}")
	}
	return fmt.Sprintf(format, p.oid, p.refPath, "")
}

type PackedRefs []*PackedRef

func (p *refParser) ParsePackedRefs() (PackedRefs, error) {
	r := make(PackedRefs, 0)
	err := dataParse(func() {
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
				r = append(r, &PackedRef{
					RefEntry: RefEntry{
						oid:     p.ParseObjectId(),
						refPath: p.ReadString(LF),
					},
				})
			}
		}
	})
	return r, err
}

func (p *refParser) ParseRefFile() (oid *ObjectId, err error) {
	err = dataParse(func() {
		oid = p.ParseObjectId()
		p.ConsumeByte(LF)
	})
	return oid, err
}
