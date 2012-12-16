//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

package parse

import "github.com/jbrukh/ggit/util"

//TODO: separate test package?

// ================================================================= //
// METHODS FOR TESTING
// ================================================================= //

func ObjectParserForString(str string) *ObjectParser {
	return &ObjectParser{
		*NewObjectIdParser(util.ReaderForString(str)), nil, nil,
	}
}
