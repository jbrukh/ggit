package diff

import (
	"github.com/jbrukh/ggit/api/format"
	"github.com/jbrukh/ggit/api/objects"
	"github.com/mikebosw/gdiff"
)

type BlobDiffer interface {
	Diff(a, b *objects.Blob) gdiff.Diff
}

func BlobDiff(a, b *objects.Blob) gdiff.Diff {
	fa, fb := format.NewStrFormat(), format.NewStrFormat()
	fa.Object(a)
	fb.Object(b)
	s1, s2 := fa.String(), fb.String()
	differ := gdiff.MyersDiffer()
	return differ.Diff(s1, s2, gdiff.LineSplit)
}
