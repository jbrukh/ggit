package api

import (
	"fmt"
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
