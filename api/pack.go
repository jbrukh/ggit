package api

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
