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
}

type PackedRefs struct {
	refs []*PackedRef
}

func (p *refParser) parsePackedRefs() *PackedRefs {
	r := new(PackedRefs)
	r.refs = make([]*PackedRef, 0)
	for !p.EOF() {
		r.refs = append(r.refs, &PackedRef{
			oid:     p.ParseObjectId(),
			refPath: p.ReadString(LF),
		})
	}
	return r
}
