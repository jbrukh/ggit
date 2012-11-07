//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import "github.com/jbrukh/ggit/api/objects"

// ================================================================= //
// CONSTANTS RELATED TO TYPES
// ================================================================= //

var objectTypes []string = []string{
	string(objects.ObjectBlob),
	string(objects.ObjectTree),
	string(objects.ObjectCommit),
	string(objects.ObjectTag),
}
