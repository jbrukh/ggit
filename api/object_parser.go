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
		hdr: nil,
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

// ================================================================= //
// GGIT OBJECT ID PARSER
// ================================================================= //

type objectIdParser struct {
	dataParser
}

// ================================================================= //
// GGIT REF PARSER
// ================================================================= //

type refParser struct {
	objectIdParser
	name string
}

func newRefParser(buf *bufio.Reader, name string) *refParser {
	return &refParser{
		objectIdParser: objectIdParser{
			dataParser{
				buf: buf,
			},
		},
		name: name,
	}
}

// ================================================================= //
// GGIT INDEX PARSER
// ================================================================= //

type indexParser dataParser
