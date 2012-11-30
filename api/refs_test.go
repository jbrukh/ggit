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
	"sort"
	"testing"
)

func Test_Sorting(t *testing.T) {
	refs := make([]objects.Ref, 3)
	refs[0] = objects.NewRef("zoo", "", nil, nil)
	refs[1] = objects.NewRef("yogurt", "", nil, nil)
	refs[2] = objects.NewRef("xavier", "", nil, nil)
	sort.Sort(refByName(refs))
	util.Assert(t, refs[0].Name() == "xavier")
	util.Assert(t, refs[1].Name() == "yogurt")
	util.Assert(t, refs[2].Name() == "zoo")
}
