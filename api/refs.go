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
