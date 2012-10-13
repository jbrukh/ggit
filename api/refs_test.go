package api

import (
	"github.com/jbrukh/ggit/util"
	"sort"
	"testing"
)

func Test_Sorting(t *testing.T) {
	refs := make([]Ref, 3)
	refs[0] = &ref{name: "zoo"}
	refs[1] = &ref{name: "yogurt"}
	refs[2] = &ref{name: "xavier"}
	sort.Sort(refByName(refs))
	util.Assert(t, refs[0].Name() == "xavier")
	util.Assert(t, refs[1].Name() == "yogurt")
	util.Assert(t, refs[2].Name() == "zoo")
}
