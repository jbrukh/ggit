package ggit

import (
	"bufio"
)

// ================================================================= //
// OBJECT HEADER PARSING
// ================================================================= //

// ObjectHeader is the deserialized (and more efficiently stored)
// version of a git object header
type objectHeader struct {
	Type ObjectType
	Size int
}

// ================================================================= //
// GGIT OBJECT PARSER
// ================================================================= //

type objectParser struct {
	dataParser
	hdr *objectHeader
}

func newObjectParser(buf *bufio.Reader) *objectParser {
	op := &objectParser{
		hdr: nil,
		dataParser: dataParser{
			buf: buf,
		},
	}
	return op
}

func (p *objectParser) ParseHeader() (*objectHeader, error) {
	err := dataParse(func() {
		p.hdr = new(objectHeader)
		p.hdr.Type = ObjectType(p.ConsumeStrings(objectTypes))
		p.ConsumeByte(SP)
		p.hdr.Size = p.ParseAtoi(NUL)
	})
	return p.hdr, err
}

func (p *objectParser) ParsePayload() (Object, error) {
	// parse header if it wasn't parsed already
	if _, e := p.ParseHeader(); e != nil {
		return nil, e
	}
	var (
		obj Object
		err error
	)

	err = dataParse(func() {
		switch p.hdr.Type {
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

// ================================================================= //
// GGIT REF PARSER
// ================================================================= //

type refParser dataParser

// ================================================================= //
// GGIT INDEX PARSER
// ================================================================= //

type indexParser dataParser
