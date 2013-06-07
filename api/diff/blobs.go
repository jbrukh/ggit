package diff

import (
	"github.com/jbrukh/ggit/api/objects"
	"github.com/mikebosw/gdiff"
)

func BlobDiff(a, b *objects.Blob) *gdiff.Diff {
	s1, s2 := string(a.Data()), string(b.Data())
	differ := gdiff.MyersDiffer()
	return differ.Diff(s1, s2, gdiff.LINE_SPLIT)
}
