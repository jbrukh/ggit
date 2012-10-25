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
)

// ================================================================= //
// OBJECT HEADER PARSING
// ================================================================= //

// ObjectHeader is the deserialized (and more efficiently stored)
// version of a git object header
type objectHeader struct {
	otype ObjectType
	size  int
}

func (h *objectHeader) Type() ObjectType {
	return h.otype
}

func (h *objectHeader) Size() int {
	return h.size
}

// ================================================================= //
// GGIT OBJECT PARSER
// ================================================================= //

type objectParser struct {
	objectIdParser
	oid *ObjectId
	hdr *objectHeader
}

func newObjectParser(buf *bufio.Reader, oid *ObjectId) *objectParser {
	op := &objectParser{
		objectIdParser: objectIdParser{
			dataParser{
				buf: buf,
			},
		},
		oid: oid,
	}
	return op
}

func (p *objectParser) ParseHeader() (*objectHeader, error) {
	err := safeParse(func() {
		p.hdr = new(objectHeader)
		p.hdr.otype = ObjectType(p.ConsumeStrings(objectTypes))
		p.ConsumeByte(SP)
		p.hdr.size = p.ParseAtoi(NUL)
	})
	if err != nil {
		return nil, err
	}
	return p.hdr, nil
}

func (p *objectParser) ParsePayload() (Object, error) {
	// parse header if it wasn't parsed already
	if p.hdr == nil {
		if _, e := p.ParseHeader(); e != nil {
			return nil, e
		}
	}
	var (
		obj Object
		err error
	)

	err = safeParse(func() {
		switch p.hdr.otype {
		case ObjectBlob:
			obj = p.parseBlob()
		case ObjectTree:
			obj = p.parseTree()
		case ObjectCommit:
			obj = p.parseCommit()
		case ObjectTag:
			obj = p.parseTag()
		default:
			panic("unsupported type")
		}
	})
	return obj, err
}
