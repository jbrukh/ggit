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
// PARSING
// ================================================================= //

// objectIdParser is a dataparser that supports parsing of oids.
type objectIdParser struct {
	util.DataParser
}

func newObjectIdParser(rd *bufio.Reader) *objectIdParser {
	return &objectIdParser{
		*util.NewDataParser(rd),
	}
}

// ParseOid reads the next objects.OidHexSize bytes from the
// Reader and places the resulting object id in oid.
func (p *objectIdParser) ParseOid() *objects.ObjectId {
	hex := string(p.Consume(objects.OidHexSize))
	oid, e := objects.OidFromString(hex)
	if e != nil {
		util.PanicErrf("expected: hex string of size %d", objects.OidHexSize)
	}
	return oid
}

// ParseOidBytes reads the next objects.OidSize bytes from
// the Reader and generates an ObjectId.
func (p *objectIdParser) ParseOidBytes() *objects.ObjectId {
	b := p.Consume(objects.OidSize)
	oid, e := objects.OidFromBytes(b)
	if e != nil {
		util.PanicErrf("expected: hash bytes %d long", objects.OidSize)
	}
	return oid
}
