//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package parse

import (
	"bufio"
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/api/token"
	"github.com/jbrukh/ggit/util"
)

// ================================================================= //
// REF PARSING
// ================================================================= //

const (
	markerRef = "ref:"
)

// refParser implements functions for parsing refs and packed refs.
type refParser struct {
	objectIdParser
	name string // the name of the ref file
}

// NewRefParser creates a new Ref parser for a ref with the
// given name.
func NewRefParser(buf *bufio.Reader, name string) *refParser {
	return &refParser{
		*NewObjectIdParser(buf),
		name,
	}
}

func (p *refParser) ParsePackedRefs() ([]objects.Ref, error) {
	r := make([]objects.Ref, 0)
	err := util.SafeParse(func() {
		for !p.EOF() {
			c := p.PeekByte()
			switch c {
			case '#':
				// if this is the first line, then it should be a comment
				// that says '# pack-refs with: <extention>' and <extention>
				// is exactly one of the items in this set: { 'peeled' }.
				// currently, we are just ignoring all comments.
				p.ReadString(token.LF)
			case '^':
				// this means the previous line is an annotated tag and the the current
				// line contains the commit that tag points to
				p.ConsumeByte('^')
				commit := p.ParseOid()
				p.ConsumeByte(token.LF)

				if l := len(r); l > 0 {
					_, oid := r[l-1].Target()
					//TODO: inefficient (copying):
					r[l-1] = objects.NewRef(r[l-1].Name(), "", oid.(*objects.ObjectId), commit)
				}
			default:
				oid := p.ParseOid()
				p.ConsumeByte(token.SP)
				name := p.ReadString(token.LF)

				r = append(r, objects.NewRef(name, "", oid, nil))
			}
		}
	})
	return r, err
}

func (p *refParser) ParseRef() (r objects.Ref, err error) {
	err = util.SafeParse(func() {
		// is it a symbolic ref?
		if p.PeekString(len(markerRef)) == markerRef {
			p.ConsumeString(markerRef)
			p.ConsumeByte(token.SP)
			spec := p.ReadString(token.LF)
			r = objects.NewRef(p.name, spec, nil, nil)
		} else {
			oid := p.ParseOid()
			p.ConsumeByte(token.LF)
			r = objects.NewRef(p.name, "", oid, nil)
		}
	})
	return r, err
}
