//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package parse

import (
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/api/token"
	"github.com/jbrukh/ggit/util"
)

const (
	markerObject = "object"
	markerType   = "type"
	markerTag    = "tag"
	markerTagger = "tagger"
)

// ================================================================= //
// PARSING
// ================================================================= //

func (p *ObjectParser) parseTag() *objects.Tag {
	p.ResetCount()

	// read the object id
	p.ConsumeString(markerObject)
	p.ConsumeByte(token.SP)
	target := p.ParseOid()
	p.ConsumeByte(token.LF)

	// read object type
	p.ConsumeString(markerType)
	p.ConsumeByte(token.SP)
	t := objects.ObjectType(p.ConsumeStrings(token.ObjectTypes))
	p.ConsumeByte(token.LF)

	// read the tag name
	p.ConsumeString(markerTag)
	p.ConsumeByte(token.SP)
	name := p.ReadString(token.LF) // gets rid of the LF!

	// read the tagger
	tagger := p.parseWhoWhen(markerTagger)
	p.ConsumeByte(token.LF)

	// read the commit message
	p.ConsumeByte(token.LF)
	msg := p.String()

	if p.Count() != p.hdr.Size() {
		util.PanicErr("payload doesn't match prescibed size")
	}

	return objects.NewTag(p.oid, target, t, p.hdr, name, msg, tagger)
}
