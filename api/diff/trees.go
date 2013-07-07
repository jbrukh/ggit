package diff

import (
	"fmt"
	"github.com/jbrukh/ggit/api"
	"github.com/jbrukh/ggit/api/objects"
	"sort"
)

type TreeDiffer interface {
	Diff(ta, tb *objects.Tree) (*TreeDiff, error)
}

type TreeDiff struct {
	edits       []*TreeEdit
	modified    []*TreeEdit
	renamed     []*TreeEdit
	insertEdits []*TreeEdit
	deleteEdits []*TreeEdit
	inserted    []*objects.TreeEntry
	deleted     []*objects.TreeEntry
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
		return treeFormat(treeFormat(fmt.Sprintf("R%d", int(te.score)), te.Before), te.After)
	}
	return ""
}

type TreeEdit struct {
	action editType
	//non-nil for delete and rename
	Before *objects.TreeEntry
	//non-nil for insert and rename
	After *objects.TreeEntry
	//for renames
	score float64
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
	blobDiffer BlobMatcher
}

func NewTreeDiffer(r api.Repository, blobDiffer BlobMatcher) TreeDiffer {
	return &treeDiffer{r, blobDiffer}
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
	result.categorizeEdits(d.repository, d.blobDiffer)
	return result, nil
}

func entryKey(entry *objects.TreeEntry, usingMode bool) string {
	key := entry.ObjectId().String() + entry.Name()
	if usingMode {
		key = fmt.Sprintf("%s%d", key, uint16(entry.Mode()))
	}
	return key
}

//find the blobs that differ between two TreeEntry slices
func findBlobDiffs(r api.Repository, a, b []*objects.TreeEntry) (blobsA, blobsB []*objects.TreeEntry, err error) {
	onlyInA, onlyInB := make(map[string]*objects.TreeEntry), make(map[string]*objects.TreeEntry)

	//determine the entries that only exist on side a or only exist on side b

	for _, entry := range a {
		onlyInA[entryKey(entry, true)] = entry
	}

	for _, entry := range b {
		key := entryKey(entry, true)
		if onlyInA[key] != nil {
			delete(onlyInA, key)
		} else {
			onlyInB[key] = entry
		}
	}

	//separate out the blobs, and explode any trees

	var explodedA, explodedB []*objects.TreeEntry

	for _, entry := range onlyInA {
		switch entry.ObjectType() {
		case objects.ObjectBlob:
			blobsA = append(blobsA, entry)
		case objects.ObjectTree:
			if exploded, err := flatten(r, entry.Name()+"/", entry); err != nil {
				return nil, nil, err
			} else {
				explodedA = append(explodedA, exploded...)
			}
		}
	}
	for _, entry := range onlyInB {
		switch entry.ObjectType() {
		case objects.ObjectBlob:
			blobsB = append(blobsB, entry)
		case objects.ObjectTree:
			if exploded, err := flatten(r, entry.Name()+"/", entry); err != nil {
				return nil, nil, err
			} else {
				explodedB = append(explodedB, exploded...)
			}
		}
	}

	//when there are no more exploded tree entries to explore, we are at the deepest level and can terminate the search
	if len(explodedA) == 0 && len(explodedB) == 0 {
		return
	}
	if moreBlobsA, moreBlobsB, err := findBlobDiffs(r, explodedA, explodedB); err != nil {
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

func (result *TreeDiff) categorizeEdits(r api.Repository, blobMatcher BlobMatcher) error {
	result.detectModified()
	if err := result.detectRenamed(r, blobMatcher); err != nil {
		return err
	}
	return nil
}

func (result *TreeDiff) detectModified() {
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
			result.insertEdits = append(result.insertEdits, edit)
		case Delete:
			result.deleteEdits = append(result.deleteEdits, edit)
		}
	}
	result.edits = conflated
}

func (result *TreeDiff) detectRenamed(r api.Repository, blobMatcher BlobMatcher) (err error) {
	//build a matrix of similarity scores
	deletes, inserts := result.deleteEdits, result.insertEdits
	scores := make([]float64, len(deletes)*len(inserts))
	for di, delete := range deletes {
		for ii, insert := range inserts {
			var objA, objB objects.Object
			if objA, err = r.ObjectFromOid(delete.Before.ObjectId()); err != nil {
				return
			}
			if objB, err = r.ObjectFromOid(insert.After.ObjectId()); err != nil {
				return err
			}
			a, _ := objA.(*objects.Blob)
			b, _ := objB.(*objects.Blob)
			index := di*(len(deletes)-1) + ii
			scores[index] = blobMatcher.Match(a, b)
		}
	}
	//get the indices of the score matrix in sorted order (i.e., from greatest score to least)
	sorted := indexedInOrder(scores)
	deletesDone, insertsDone := make([]bool, len(deletes)), make([]bool, len(inserts))
	renameCount := 0
	for _, index := range sorted {
		score := scores[index]
		if score < 60 {
			break
		}
		di := int(index) / (len(deletes) - 1)
		ii := int(index) - (di * (len(deletes) - 1))
		if deletesDone[di] || insertsDone[ii] {
			continue
		}
		renameCount++
		rename, insert := result.deleteEdits[di], result.insertEdits[ii]
		result.insertEdits[ii], result.deleteEdits[di] = nil, nil
		rename.After = insert.After
		rename.action = Rename
		rename.score = score
		result.renamed = append(result.renamed, rename)
		deletesDone[di], insertsDone[ii] = true, true
	}
	//prune the removed (nil) edits and populate the \deleted\ []*objects.TreeEntry and \inserted\ []*objects.TreeEntry
	deleteEdits, insertEdits := result.deleteEdits, result.insertEdits
	result.edits, result.deleteEdits, result.insertEdits = nil, nil, nil
	for _, v := range insertEdits {
		if v == nil {
			continue
		}
		result.edits = append(result.edits, v)
		result.insertEdits = append(result.insertEdits, v)
		result.inserted = append(result.inserted, v.After)
	}
	for _, v := range deleteEdits {
		if v == nil {
			continue
		}
		result.edits = append(result.edits, v)
		result.deleteEdits = append(result.deleteEdits, v)
		result.deleted = append(result.deleted, v.Before)
	}
	result.edits = append(result.edits, result.modified...)
	result.edits = append(result.edits, result.renamed...)
	return
}

type byValueInReverse struct {
	values  []float64
	indices []int
}

func indexedInOrder(scores []float64) []int {
	//1. create a slice of score indices with initial values from 0 to len(scores) - 1
	//2. sort that slice as though its values were the scores
	bv := &byValueInReverse{scores, make([]int, len(scores))}
	for i := 0; i < len(bv.values); i++ {
		bv.indices[i] = i
	}
	sort.Sort(bv)
	return bv.indices
}

func (bv *byValueInReverse) Len() int {
	return len(bv.values)
}

func (bv *byValueInReverse) Swap(i, j int) {
	bv.indices[i], bv.indices[j] = bv.indices[j], bv.indices[i]
}

func (bv *byValueInReverse) Less(i, j int) bool {
	a, b := bv.values[bv.indices[i]], bv.values[bv.indices[j]]
	return int(a) > int(b)
}
