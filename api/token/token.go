//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors

package token

import "github.com/jbrukh/ggit/api/objects"

const (
	NUL = '\000'
	SP  = ' '
	LF  = '\n'
	LT  = '<'
	GT  = '>'
	TAB = '\t'
)

const (
	PLUS  = "+"
	MINUS = "-"
)

var ObjectTypes []string = []string{
	string(objects.ObjectBlob),
	string(objects.ObjectTree),
	string(objects.ObjectCommit),
	string(objects.ObjectTag),
}
