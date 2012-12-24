//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

/*
file_modes.go implements the git-supported file modes used mainly in trees.
*/
package parse

import (
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/util"
)

// ================================================================= //
// FILE MODE
// ================================================================= //

func assertFileMode(u uint16) (objects.FileMode, bool) {
	m := objects.FileMode(u)
	switch m {
	case objects.ModeNew,
		objects.ModeTree,
		objects.ModeBlob,
		objects.ModeBlobExec,
		objects.ModeLink,
		objects.ModeCommit:
		return m, true
	}
	return 0, false
}

// ================================================================= //
// OBJECT PARSER FUNCTIONS FOR FILE MODE
// ================================================================= //

func (p *objectParser) ParseFileMode(delim byte) (mode objects.FileMode) {
	var ok bool
	if mode, ok = assertFileMode(uint16(p.ParseInt(delim, 8, 32))); !ok {
		util.PanicErrf("expected: filemode")
	}
	return
}
