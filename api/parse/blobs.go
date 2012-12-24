//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
blobs.go implements ggit Blob objects and their parsing and formatting.
*/
package parse

import (
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/util"
)

// ================================================================= //
// OBJECT PARSER
// ================================================================= //

// parseBlob parses the payload of a binary blob object
// and converts it to Blob. If there are parsing errors,
// it panics with parseErr, so this method should be
// called as a parameter a safeParse().
func (p *objectParser) parseBlob() *objects.Blob {

	p.ResetCount()
	data := p.Bytes()
	b := objects.NewBlob(p.oid, p.hdr, data)

	if p.Count() != p.hdr.Size() {
		util.PanicErr("payload doesn't match prescibed size")
	}

	return b
}
