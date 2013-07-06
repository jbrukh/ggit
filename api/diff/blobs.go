package diff

import (
	"github.com/jbrukh/ggit/api/format"
	"github.com/jbrukh/ggit/api/objects"
	"github.com/mikebosw/gdiff"
)

type BlobDiffer interface {
	Diff(a, b *objects.Blob) gdiff.Diff
}

type BlobMatcher interface {
	Match(a, b *objects.Blob) float64
}

type blobComparator int

func NewBlobComparator() *blobComparator {
	differ := blobComparator(0)
	return &differ
}

func (differ *blobComparator) Match(a, b *objects.Blob) float64 {
	diff := differ.Diff(a, b, gdiff.CharSplit)
	return gdiff.SimpleComparator().Score(diff)
}

func (*blobComparator) Diff(a, b *objects.Blob, seq gdiff.Sequencer) gdiff.Diff {
	fa, fb := format.NewStrFormat(), format.NewStrFormat()
	fa.Object(a)
	fb.Object(b)
	s1, s2 := fa.String(), fb.String()
	differ := gdiff.MyersDiffer()
	return differ.Diff(s1, s2, seq)
}
