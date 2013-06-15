package diff

import (
	"fmt"
	"github.com/jbrukh/ggit/api"
	"github.com/jbrukh/ggit/api/objects"
	"github.com/mikebosw/gdiff"
	"sort"
)

type TreeDiffer interface {
	Diff(ta, tb *objects.Tree) (*TreeDiff, error)
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
	return fmt.Sprintf("%s %s", prefix, te.Name())
}

func (te *treeEdit) String() string {
	switch te.action {
	case Insert:
		return treeFormat("+", te.after)
	case Delete:
		return treeFormat("-", te.before)
	case Modify:
		return treeFormat("M", te.before)
	case Rename:
		return treeFormat(treeFormat("R", te.before), te.after)
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
	Modify editType = 'm'
	Rename editType = 'r'
)

type treeDiffer struct {
	repository api.Repository
}

func NewTreeDiffer(r api.Repository) TreeDiffer {
	return &treeDiffer{r}
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

func diffEntries(entriesA, entriesB []*objects.TreeEntry) *TreeDiff {
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
	result.checkForModified()
	return result
}

func (td *TreeDiff) checkForModified() {
	var conflated []*treeEdit
	paths := make(map[string]*treeEdit)
	for _, edit := range td.edits {
		var blob *objects.TreeEntry
		switch edit.action {
		case Insert:
			blob = edit.after
		case Delete:
			blob = edit.before
		default:
			return
		}
		if existing := paths[blob.Name()]; existing == nil {
			conflated = append(conflated, edit)
			paths[blob.Name()] = edit
		} else {
			switch existing.action {
			case Insert:
				existing.before = blob
			case Delete:
				existing.after = blob
			}
			existing.action = Modify
		}
	}
	td.edits = conflated
}

func (d *treeDiffer) Diff(ta, tb *objects.Tree) (result *TreeDiff, err error) {
	entriesA, entriesB := ta.Entries(), tb.Entries()
	if entriesA, err = flatten(d.repository, "", entriesA); err != nil {
		return nil, err
	}
	if entriesB, err = flatten(d.repository, "", entriesB); err != nil {
		return nil, err
	}
	return diffEntries(entriesA, entriesB), nil
}

func flatten(r api.Repository, base string, entries []*objects.TreeEntry) (result []*objects.TreeEntry, err error) {
	result = make([]*objects.TreeEntry, 0)
	for _, entry := range entries {
		switch entry.ObjectType() {
		case objects.ObjectBlob:
			result = append(result, objects.NewTreeEntry(entry.Mode(), entry.ObjectType(), base+entry.Name(), entry.ObjectId()))
		case objects.ObjectTree:
			var object objects.Object
			object, err = r.ObjectFromOid(entry.ObjectId())
			if err != nil {
				return nil, err
			}
			tree, _ := object.(*objects.Tree)
			var blobs []*objects.TreeEntry
			blobs, err = flatten(r, base+entry.Name()+"/", tree.Entries())
			result = append(result, blobs...)
		}
	}
	return result, nil
}
