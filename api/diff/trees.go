package diff

import (
	"github.com/jbrukh/ggit/api/objects"
	"sort"
)

type TreeDiffer interface {
	Diff(ta, tb *objects.Tree) *TreeDiff
}

type TreeDiff struct {
	edits []*treeEdit
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
	result := new(TreeDiff)
	var entryA, entryB *objects.TreeEntry
	for i, j := 0, 0; i < len(entriesA) || j < len(entriesB); {
		if entryA == nil {
			if len(entriesA) > i {
				entryA = entriesA[i]
			}
		}
		if entryB == nil {
			if len(entriesB) > i {
				entryB = entriesB[i]
			}
		}
		var op editType
		switch diff := compare(entryA, entryB); diff {
		case Less:
			op = Delete
			i++
		case More:
			op = Insert
			j++
		case Same:
			i, j = i+1, j+1
			if entryA.Name() != entryB.Name() {
				op = Rename
			} else {
				continue
			}
		}
		edit := &treeEdit{
			op,
			entryA,
			entryB,
		}
		result.edits = append(result.edits, edit)
	}
	return result
}
