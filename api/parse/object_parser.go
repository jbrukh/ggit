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
// GGIT OBJECT PARSER
// ================================================================= //

type ObjectParser interface {
	ParsePayload() (objects.Object, error)
	ParseHeader() (*objects.ObjectHeader, error)
}

type objectParser struct {
	objectIdParser
	oid *objects.ObjectId
	hdr *objects.ObjectHeader
}

func NewObjectParser(buf *bufio.Reader, oid *objects.ObjectId) *objectParser {
	op := &objectParser{
		*NewObjectIdParser(buf),
		oid,
		nil,
	}
	return op
}

func (p *objectParser) ParseHeader() (*objects.ObjectHeader, error) {
	err := util.SafeParse(func() {
		ot := objects.ObjectType(p.ConsumeStrings(token.ObjectTypes))
		p.ConsumeByte(token.SP)
		size := p.ParseAtoi(token.NUL)
		p.hdr = objects.NewObjectHeader(ot, size)
	})
	if err != nil {
		return nil, err
	}
	return p.hdr, nil
}

func (p *objectParser) ParsePayload() (objects.Object, error) {
	// parse header if it wasn't parsed already
	if p.hdr == nil {
		if _, e := p.ParseHeader(); e != nil {
			return nil, e
		}
	}
	var (
		obj objects.Object
		err error
	)

	err = util.SafeParse(func() {
		switch p.hdr.Type() {
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
