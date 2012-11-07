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
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/util"
)

// ================================================================= //
// GGIT OBJECT PARSER
// ================================================================= //

type objectParser struct {
	objectIdParser
	oid *objects.ObjectId
	hdr *ObjectHeader
}

func newObjectParser(buf *bufio.Reader, oid *objects.ObjectId) *objectParser {
	op := &objectParser{
		objectIdParser: objectIdParser{
			*util.NewDataParser(buf),
		},
		oid: oid,
	}
	return op
}

func (p *objectParser) ParseHeader() (*ObjectHeader, error) {
	err := util.SafeParse(func() {
		p.hdr = new(ObjectHeader)
		p.hdr.otype = objects.ObjectType(p.ConsumeStrings(objectTypes))
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

	err = util.SafeParse(func() {
		switch p.hdr.otype {
		case objects.ObjectBlob:
			obj = p.parseBlob()
		case objects.ObjectTree:
			obj = p.parseTree()
		case objects.ObjectCommit:
			obj = p.parseCommit()
		case objects.ObjectTag:
			obj = p.parseTag()
		default:
			panic("unsupported type")
		}
	})
	return obj, err
}
