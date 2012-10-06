package api

import (
	"sort"
	"testing"
)

func Test_Sorting(t *testing.T) {
	refs := make([]Ref, 3)
	refs[0] = &ref{nil, "zoo", nil}
	refs[1] = &ref{nil, "yogurt", nil}
	refs[2] = &ref{nil, "xavier", nil}
	sort.Sort(refByName(refs))
	assert(t, refs[0].Name() == "xavier")
	assert(t, refs[1].Name() == "yogurt")
	assert(t, refs[2].Name() == "zoo")
}
