//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import (
	"github.com/jbrukh/ggit/api/objects"
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

func (p *objectParser) parseTag() *objects.Tag {
	p.ResetCount()

	// read the object id
	p.ConsumeString(markerObject)
	p.ConsumeByte(SP)
	target := p.ParseOid()
	p.ConsumeByte(LF)

	// read object type
	p.ConsumeString(markerType)
	p.ConsumeByte(SP)
	t := objects.ObjectType(p.ConsumeStrings(objectTypes))
	p.ConsumeByte(LF)

	// read the tag name
	p.ConsumeString(markerTag)
	p.ConsumeByte(SP)
	name := p.ReadString(LF) // gets rid of the LF!

	// read the tagger
	tagger := p.parseWhoWhen(markerTagger)
	p.ConsumeByte(LF)

	// read the commit message
	p.ConsumeByte(LF)
	msg := p.String()

	if p.Count() != p.hdr.Size() {
		util.PanicErr("payload doesn't match prescibed size")
	}

	return objects.NewTag(p.oid, target, t, p.hdr, name, msg, tagger)
}
