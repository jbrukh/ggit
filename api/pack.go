package api

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"sort"
)

const (
	PackSignature    = "PACK"    //0x5041434b
	PackIdxSignature = "\377tOc" //0xff744f63
	PackVersion      = 2
)

type PackedObjectType byte

const (
	COMMIT              PackedObjectType = 1
	TREE                PackedObjectType = 2
	BLOB                PackedObjectType = 3
	TAG                 PackedObjectType = 4
	OBJECT_OFFSET_DELTA PackedObjectType = 6
	OBJECT_REF_DELTA    PackedObjectType = 7
)

type packedObject struct {
	Object
	bytes *[]byte
}

type Pack struct {
	// GIT currently accepts version number 2 or 3 but
	// generates version 2 only.
	version int32
	// the unpacked objects
	content []*packedObject
	*Idx
}

type Idx struct {
	// the object ids sorted by offset
	entries []*PackedObjectId
	// the object ids mapped by oid
	entriesById map[string]*PackedObjectId
	// number of objects contained in the pack (network
	// byte order)
	count int64
	// copy of the checksum for this idx file's corresponding pack file.
	packChecksum *ObjectId
	// checksum for this idx file.
	idxChecksum *ObjectId
}

type PackedObjectId struct {
	ObjectId
	offset int64
	crc32  int64
}

// ================================================================= //
// GGIT PACK PARSER
// ================================================================= //

type packIdxParser struct {
	idxParser  objectIdParser
	packParser dataParser
	name       string
}

// ================================================================= //
// Sorting of PackedObjectIds by offset.
// ================================================================= //

type packedObjectIds []*PackedObjectId

func (e packedObjectIds) Less(i, j int) bool {
	s := []*PackedObjectId(e)
	a, b := s[i], s[j]
	return a.offset < b.offset
}

func (e packedObjectIds) Swap(i, j int) {
	s := []*PackedObjectId(e)
	s[i], s[j] = s[j], s[i]
}

func (e packedObjectIds) Len() int {
	s := []*PackedObjectId(e)
	return len(s)
}

// ================================================================= //
// .idx parsing.
// ================================================================= //

func (p *packIdxParser) parseIdx() *Idx {
	p.idxParser.ConsumeString(PackIdxSignature)
	p.idxParser.ConsumeBytes([]byte{0, 0, 0, PackVersion})
	var counts [256]int64
	for i := range counts {
		counts[i] = p.idxParser.ParseIntBigEndian(4)
	}
	//discard the fan-out values, just use the largest value, which is the total # of objects:
	count := counts[255]
	entries := make([]*PackedObjectId, count, count)
	entriesByOid := make(map[string]*PackedObjectId)
	for i := int64(0); i < count; i++ {
		b := p.idxParser.ReadNBytes(20)
		representation := fmt.Sprintf("%x", b)
		entries[i] = &PackedObjectId{
			ObjectId: ObjectId{
				b,
				representation,
			},
		}
		entriesByOid[representation] = entries[i]
	}
	for i := int64(0); i < count; i++ {
		entries[i].crc32 = int64(p.idxParser.ParseIntBigEndian(4))
	}
	for i := int64(0); i < count; i++ {
		//TODO: 8-byte #'s for some offsets for some pack files (packs > 2gb)
		entries[i].offset = p.idxParser.ParseIntBigEndian(4)
	}
	checksumPack := p.idxParser.ReadNBytes(20)
	packChecksum := &ObjectId{
		bytes: checksumPack,
		repr:  fmt.Sprintf("%x", checksumPack),
	}
	//TODO: check the checksum
	checksumIdx := p.idxParser.ReadNBytes(20)
	idxChecksum := &ObjectId{
		bytes: checksumIdx,
		repr:  fmt.Sprintf("%x", checksumIdx),
	}
	if !p.idxParser.EOF() {
		panicErrf("Found extraneous bytes! %x", p.idxParser.Bytes())
	}
	//order by offset
	sort.Sort(packedObjectIds(entries))
	return &Idx{
		entries,
		entriesByOid,
		count,
		packChecksum,
		idxChecksum,
	}
}

// ================================================================= //
// .pack and pack entry parsing.
// ================================================================= //

type packedObjectParser struct {
	*objectParser
	bytes *[]byte
}

func newPackedObjectParser(data *[]byte, oid *ObjectId) (p *packedObjectParser, e error) {
	compressedReader := bytes.NewReader(*data)
	var zr io.ReadCloser
	if zr, e = zlib.NewReader(compressedReader); e == nil {
		exploder := &dataParser{
			bufio.NewReader(zr),
			0,
		}
		exploded := exploder.Bytes()
		explodedReader := bufio.NewReader(bytes.NewReader(exploded))
		op := newObjectParser(explodedReader, oid)
		pop := packedObjectParser{
			op,
			&exploded,
		}
		p = &pop
	}
	return
}

func (p *packIdxParser) parsePack() *Pack {
	idx := p.parseIdx()
	objects := make([]*packedObject, idx.count)
	p.packParser.ConsumeString(PackSignature)
	p.packParser.ConsumeBytes([]byte{0, 0, 0, PackVersion})
	count := p.packParser.ParseIntBigEndian(4)
	if count != idx.count {
		panicErrf("Pack file count doesn't match idx file count for pack-%s!", p.name) //todo: don't panic.
	}
	entries := &idx.entries
	data := p.packParser.Bytes()
	for i := range *entries {
		objects[i] = parseEntry(&data, i, idx, &objects)
	}
	return &Pack{
		PackVersion,
		objects,
		idx,
	}
}

func parseEntry(packedData *[]byte, i int, idx *Idx, packedObjects *[]*packedObject) (object *packedObject) {
	return nil
}

func parseNonDeltaEntry(bytes *[]byte, pot PackedObjectType, oid *ObjectId, size int64) (object *packedObject) {
	var (
		dp  *packedObjectParser
		err error
	)
	if dp, err = newPackedObjectParser(bytes, oid); err != nil {
		panicErr(err.Error())
	}
	switch pot {
	case BLOB:
		object = dp.parseBlob(size)
	case COMMIT:
		object = dp.parseCommit(size)
	case TREE:
		object = dp.parseTree(size)
	case TAG:
		object = dp.parseTag(size)
	}
	return
}

func (dp *packedObjectParser) parseCommit(size int64) *packedObject {
	dp.hdr = &objectHeader{
		ObjectCommit,
		int(size),
	}
	commit := dp.objectParser.parseCommit()

	return &packedObject{
		commit,
		dp.bytes,
	}
}
func (dp *packedObjectParser) parseTag(size int64) *packedObject {
	dp.hdr = &objectHeader{
		ObjectTag,
		int(size),
	}
	tag := dp.objectParser.parseTag()
	return &packedObject{
		tag,
		dp.bytes,
	}
}

func (dp *packedObjectParser) parseBlob(size int64) *packedObject {
	blob := new(Blob)
	blob.data = dp.Bytes()
	blob.oid = dp.objectParser.oid
	blob.hdr = &objectHeader{
		ObjectBlob,
		int(size),
	}
	return &packedObject{
		blob,
		&blob.data,
	}
}

func (dp *packedObjectParser) parseTree(size int64) *packedObject {
	dp.hdr = &objectHeader{
		ObjectTree,
		int(size),
	}
	tree := dp.objectParser.parseTree()
	return &packedObject{
		tree,
		dp.bytes,
	}
}
