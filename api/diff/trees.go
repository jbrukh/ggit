package diff

import (
	"fmt"
	"github.com/jbrukh/ggit/api/objects"
	"github.com/mikebosw/gdiff"
	"sort"
)

type TreeDiffer interface {
	Diff(ta, tb *objects.Tree) *TreeDiff
}

type TreeDiff struct {
	edits []*treeEdit
}

func (td *TreeDiff) String() string {
	result := ""
	for _, v := range td.edits {
		result += v.String() + "\n"
	}
	return result
}

func treeFormat(prefix string, te *objects.TreeEntry) string {
	return fmt.Sprintf("%s %s %s", prefix, te.ObjectType().String(), te.ObjectId().String())
}

func (te *treeEdit) String() string {
	switch te.action {
	case Insert:
		return treeFormat("+", te.after)
	case Delete:
		return treeFormat("-", te.before)
	case Rename:
		return treeFormat(treeFormat("<>", te.before), te.after)
	}
	return ""
}

type treeEdit struct {
	action editType
	//non-nil for delete and rename
	before *objects.TreeEntry
	//non-nil for insert and rename
	after *objects.TreeEntry
}

type editType rune

const (
	Insert editType = 'i'
	Delete editType = 'd'
	Rename editType = 'm'
)

type treeDiffer struct{}

func NewTreeDiffer() TreeDiffer {
	return &treeDiffer{}
}

type byOid []*objects.TreeEntry

func (e byOid) Len() int {
	return len(([]*objects.TreeEntry(e)))
}

func (e byOid) Swap(i, j int) {
	entries := ([]*objects.TreeEntry(e))
	entries[i], entries[j] = entries[j], entries[i]
}

func (e byOid) Less(i, j int) bool {
	entries := []*objects.TreeEntry(e)
	var a, b *objects.TreeEntry = entries[i], entries[j]
	return compare(a, b) == Less
}

type order rune

const (
	Less order = iota
	More
	Same
)

func compare(a, b *objects.TreeEntry) order {
	if a == b {
		return Same
	}
	if a == nil {
		return Less
	}
	if b == nil {
		return More
	}
	aId, bId := a.ObjectId().String(), b.ObjectId().String()
	if aId < bId {
		return Less
	}
	if aId > bId {
		return More
	}
	return Same
}

func (d *treeDiffer) Diff(ta, tb *objects.Tree) *TreeDiff {
	entriesA, entriesB := ta.Entries(), tb.Entries()
	sort.Sort(byOid(entriesA))
	sort.Sort(byOid(entriesB))
	idsA := ""
	for _, v := range entriesA {
		idsA += fmt.Sprintln(v.ObjectId().String())
	}
	idsB := ""
	for _, v := range entriesB {
		idsB += fmt.Sprintln(v.ObjectId().String())
	}
	result := new(TreeDiff)

	tDiff := gdiff.MyersDiffer().Diff(idsA, idsB, gdiff.LineSplit)
	for _, edit := range tDiff.Edits() {
		switch edit.Type {
		case gdiff.Insert:
			for i := edit.Start; i <= edit.End; i++ {
				result.edits = append(result.edits, &treeEdit{
					action: Insert,
					before: nil,
					after:  entriesB[i],
				})
			}
		case gdiff.Delete:
			for i := edit.Start; i <= edit.End; i++ {
				result.edits = append(result.edits, &treeEdit{
					action: Delete,
					before: entriesA[i],
					after:  nil,
				})
			}
		}
	}
	return result
}
