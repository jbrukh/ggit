package api

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

type PackedRef struct {
	oid     *ObjectId
	refPath string

	// if this is an annotated tag
	// we may have the pointed to commit here
	// as an optimization
	cid *ObjectId
}

type PackedRefs struct {
	refs []*PackedRef
}

func (p *refParser) parsePackedRefs() *PackedRefs {
	r := new(PackedRefs)
	r.refs = make([]*PackedRef, 0)
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
			if len(r.refs) > 0 {
				r.refs[len(r.refs)-1].cid = cid
			}
		default:
			r.refs = append(r.refs, &PackedRef{
				oid:     p.ParseObjectId(),
				refPath: p.ReadString(LF),
			})
		}

	}
	return r
}
