package diff

import (
	"fmt"
	"github.com/jbrukh/ggit/api"
	"github.com/jbrukh/ggit/api/objects"
)

type TreeDiffer interface {
	Diff(ta, tb *objects.Tree) (*TreeDiff, error)
}

type TreeDiff struct {
	edits    []*TreeEdit
	modified []*TreeEdit
	inserted []*objects.TreeEntry
	deleted  []*objects.TreeEntry
}

func (td *TreeDiff) Modified() []*TreeEdit {
	return td.modified
}

func (td *TreeDiff) Inserted() []*objects.TreeEntry {
	return td.inserted
}

func (td *TreeDiff) Deleted() []*objects.TreeEntry {
	return td.deleted
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

func (te *TreeEdit) String() string {
	switch te.action {
	case Insert:
		return treeFormat("A", te.After)
	case Delete:
		return treeFormat("D", te.Before)
	case Modify:
		return treeFormat("M", te.Before)
	case Rename:
		return treeFormat(treeFormat("R", te.Before), te.After)
	}
	return ""
}

type TreeEdit struct {
	action editType
	//non-nil for delete and rename
	Before *objects.TreeEntry
	//non-nil for insert and rename
	After *objects.TreeEntry
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

func (d *treeDiffer) Diff(ta, tb *objects.Tree) (*TreeDiff, error) {
	result := new(TreeDiff)
	entriesA, entriesB := ta.Entries(), tb.Entries()
	if blobsA, blobsB, err := findBlobDiffs(d.repository, entriesA, entriesB); err != nil {
		return nil, err
	} else {
		for _, blob := range blobsA {
			result.makeDeletionEdit(blob)
		}
		for _, blob := range blobsB {
			result.makeInsertionEdit(blob)
		}
	}
	result.categorizeEdits()
	return result, nil
}

func findBlobDiffs(r api.Repository, a, b []*objects.TreeEntry) (blobsA, blobsB []*objects.TreeEntry, err error) {
	var mixedA, mixedB []*objects.TreeEntry

	onlyInA := make(map[string]*objects.TreeEntry)

	for _, entry := range a {
		onlyInA[entry.ObjectId().String()] = entry
	}

	for _, entry := range b {
		id := entry.ObjectId().String()
		if onlyInA[id] != nil {
			delete(onlyInA, id)
		} else {
			switch entry.ObjectType() {
			case objects.ObjectBlob:
				blobsB = append(blobsB, entry)
			case objects.ObjectTree:
				if exploded, err := flatten(r, entry.Name()+"/", entry); err != nil {
					return nil, nil, err
				} else {
					mixedB = append(mixedB, exploded...)
				}
			}
		}
	}

	for _, entry := range onlyInA {
		switch entry.ObjectType() {
		case objects.ObjectBlob:
			blobsA = append(blobsA, entry)
		case objects.ObjectTree:
			if exploded, err := flatten(r, entry.Name()+"/", entry); err != nil {
				return nil, nil, err
			} else {
				mixedA = append(mixedA, exploded...)
			}
		}
	}

	if len(mixedA) == 0 && len(mixedB) == 0 {
		return
	}
	if moreBlobsA, moreBlobsB, err := findBlobDiffs(r, mixedA, mixedB); err != nil {
		return nil, nil, err
	} else {
		blobsA = append(blobsA, moreBlobsA...)
		blobsB = append(blobsB, moreBlobsB...)
	}
	return
}

func (result *TreeDiff) makeInsertionEdit(entry *objects.TreeEntry) {
	result.edits = append(result.edits, &TreeEdit{
		action: Insert,
		Before: nil,
		After:  entry,
	})
}

func (result *TreeDiff) makeDeletionEdit(entry *objects.TreeEntry) {
	result.edits = append(result.edits, &TreeEdit{
		action: Delete,
		Before: entry,
		After:  nil,
	})
}

func flatten(r api.Repository, base string, treeEntry *objects.TreeEntry) (result []*objects.TreeEntry, err error) {
	result = make([]*objects.TreeEntry, 0)
	var object objects.Object
	object, err = r.ObjectFromOid(treeEntry.ObjectId())
	if err != nil {
		return nil, err
	}
	tree, _ := object.(*objects.Tree)
	for _, entry := range tree.Entries() {
		result = append(result, objects.NewTreeEntry(entry.Mode(), entry.ObjectType(), base+entry.Name(), entry.ObjectId()))
	}
	return result, nil
}

func (result *TreeDiff) categorizeEdits() {
	var conflated []*TreeEdit
	paths := make(map[string]*TreeEdit)
	uniques := make(map[string]*TreeEdit)
	for _, edit := range result.edits {
		var blob *objects.TreeEntry
		switch edit.action {
		case Insert:
			blob = edit.After
		case Delete:
			blob = edit.Before
		default:
			return
		}
		if existing := paths[blob.Name()]; existing == nil {
			conflated = append(conflated, edit)
			paths[blob.Name()] = edit
			uniques[blob.Name()] = edit
		} else {
			delete(uniques, blob.Name())
			switch existing.action {
			case Insert:
				existing.Before = blob
			case Delete:
				existing.After = blob
			}
			existing.action = Modify
			result.modified = append(result.modified, existing)
		}
	}
	for _, edit := range uniques {
		switch edit.action {
		case Insert:
			result.inserted = append(result.inserted, edit.After)
		case Delete:
			result.deleted = append(result.deleted, edit.Before)
		}
	}
	result.edits = conflated
}
