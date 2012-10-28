package api

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	PackSignature    = "PACK"    //0x5041434b
	PackIdxSignature = "\377tOc" //0xff744f63
	PackVersion      = 2
)

type PackedObjectType byte

const (
	PackedCommit      PackedObjectType = 1
	PackedTree        PackedObjectType = 2
	PackedBlob        PackedObjectType = 3
	PackedTag         PackedObjectType = 4
	ObjectOffsetDelta PackedObjectType = 6
	ObjectRefDelta    PackedObjectType = 7
)

type PackedObject struct {
	Object
	bytes   []byte
	DeltaOf *ObjectId
}

type Pack struct {
	// Git currently accepts version number 2 or 3 but
	// generates version 2 only.
	version int32
	// the unpacked objects sorted by offset
	content []*PackedObject
	idx     *Idx
	name    string
}

type Idx struct {
	// the object ids sorted by offset
	entries []*PackedObjectId
	// the object ids mapped by oid
	entriesById map[string]*PackedObjectId
	// number of objects contained in the pack (network
	// byte order)
	count int64
	// copy of the checksum for this idx file's
	// corresponding pack file.
	packChecksum *ObjectId
	// checksum for this idx file.
	idxChecksum *ObjectId
}

type PackedObjectId struct {
	*ObjectId
	offset int64
	crc32  int64
	index  int
}

// Returns the one Object in this pack with the given ObjectId,
// or nil, NoSuchObject if no such Object is in this pack.
func (pack *Pack) unpack(oid *ObjectId) (obj Object, result packSearch) {
	s := oid.String()
	if entry := pack.idx.entriesById[s]; entry != nil {
		obj, result = pack.content[entry.index].Object, OneSuchObject
	}
	return
}

type packSearch byte

const (
	NoSuchObject    packSearch = 0
	OneSuchObject   packSearch = 1
	MultipleObjects packSearch = 2
)

func (pack *Pack) unpackFromShortOid(short string) (obj Object, result packSearch) {
	var already bool
	for _, oid := range pack.idx.entries {
		if s := oid.String(); strings.HasPrefix(s, short) {
			if already {
				return nil, MultipleObjects
			}
			obj, result = pack.unpack(oid.ObjectId)
			already = true
		}
	}
	return
}

// Returns the object for the given short oid, if exactly one such object exists.
// Otherwise returns nil, false.
func unpackFromShortOid(packs []*Pack, short string) (obj Object, ok bool) {
	//this function could be O(1) if we... tried... hard enough. tried. get it?
	//awwwwwwwwwwwwwwwww yeh
	var already bool
	for _, pack := range packs {
		var result packSearch
		if obj, result = pack.unpackFromShortOid(short); result != NoSuchObject {
			if already || result == MultipleObjects {
				return nil, false
			}
			already = true
		}
	}
	return obj, already
}

func unpack(packs []*Pack, oid *ObjectId) (obj Object, ok bool) {
	var result packSearch
	for _, pack := range packs {
		if obj, result = pack.unpack(oid); result == OneSuchObject {
			//trust for now that there will only be one matching object among the packs.
			return obj, true
		}
	}
	return
}

func objectIdsFromPacks(packs []*Pack) (ids []*ObjectId) {
	var count int64
	for _, pack := range packs {
		count += pack.idx.count
	}
	ids = make([]*ObjectId, count, count)
	i := 0
	for _, pack := range packs {
		for _, id := range pack.idx.entries {
			ids[i] = id.ObjectId
			i++
		}
	}
	return ids
}

func objectsFromPacks(packs []*Pack) (objects []*PackedObject) {
	var count int64
	for _, pack := range packs {
		count += pack.idx.count
	}
	objects = make([]*PackedObject, count, count)
	i := 0
	for _, pack := range packs {
		for _, entry := range pack.content {
			objects[i] = entry
			i++
		}
	}
	return objects
}

// ================================================================= //
// GGIT PACK PARSER
// ================================================================= //

type packIdxParser struct {
	idxParser  objectIdParser
	packParser dataParser
	name       string
	packFile   *os.File
}

func newPackIdxParser(idx *bufio.Reader, pack *os.File, name string) *packIdxParser {
	return &packIdxParser{
		idxParser: objectIdParser{
			dataParser{
				buf: idx,
			},
		},
		packParser: dataParser{
			buf: bufio.NewReader(pack),
		},
		name:     name,
		packFile: pack,
	}
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
	//discard the fan-out values, just use the largest value,
	//which is the total # of objects:
	count := counts[255]
	entries := make([]*PackedObjectId, count, count)
	entriesByOid := make(map[string]*PackedObjectId)
	for i := int64(0); i < count; i++ {
		b := p.idxParser.ReadNBytes(20)
		representation := fmt.Sprintf("%x", b)
		entries[i] = &PackedObjectId{
			ObjectId: &ObjectId{
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
	for i, v := range entries {
		v.index = i
	}
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
	bytes []byte
}

func newPackedObjectParser(data *[]byte, oid *ObjectId) (p *packedObjectParser, e error) {
	compressedReader := bytes.NewReader(*data)
	var zr io.ReadCloser
	if zr, e = zlib.NewReader(compressedReader); e == nil {
		defer zr.Close()
		exploder := &dataParser{
			bufio.NewReader(zr),
			0,
		}
		exploded := exploder.Bytes()
		explodedReader := bufio.NewReader(bytes.NewReader(exploded))
		op := newObjectParser(explodedReader, oid)
		pop := packedObjectParser{
			op,
			exploded,
		}
		p = &pop
	}
	return
}

func (p *packIdxParser) parsePack() *Pack {
	idx := p.parseIdx()
	objects := make([]*PackedObject, idx.count)
	p.packParser.ConsumeString(PackSignature)
	p.packParser.ConsumeBytes([]byte{0, 0, 0, PackVersion})
	count := p.packParser.ParseIntBigEndian(4)
	if count != idx.count {
		panicErrf("Pack file count doesn't match idx file count for pack-%s!", p.name) //todo: don't panic.
	}
	for i := range idx.entries {
		objects[i] = parseEntry(p.packFile, i, idx, objects)
	}
	return &Pack{
		PackVersion,
		objects,
		idx,
		p.name,
	}
}

func parseEntry(file *os.File, i int, idx *Idx, objects []*PackedObject) (obj *PackedObject) {
	//TODO: break this function up... or at least segment it more clearly
	if len(objects) > i && objects[i] != nil {
		//sometimes (for ref delta objects) we jump ahead in the []byte
		return objects[i]
	}
	entries := idx.entries
	v := entries[i]
	var entryLen int64
	if i+1 < len(entries) {
		entryLen = entries[i+1].offset - v.offset
	} else {
		if info, err := file.Stat(); err != nil {
			panicErrf("Could not determine size of pack file %s: %s", file.Name(), err)
		} else {
			entryLen = info.Size() - v.offset
		}
	}
	data := make([]byte, entryLen, entryLen)
	if _, err := file.ReadAt(data, v.offset); err != nil {
		fmt.Printf("Could not read %d bytes from %d of pack file %s: %s", len(data), v.offset, file.Name(), err)
	}
	var (
		size int64
		err  error
	)
	// keep track of bytes read so that, in conjunction with the next entry's offset, we can know where the next
	// object in the pack begins.
	var cursor int
	headerHeader := data[cursor]
	cursor++
	typeBits := (headerHeader & 127) >> 4
	sizeBits := (headerHeader & 15)
	//collect remaining size bytes, if any.
	sizeBytes := fmt.Sprintf("%04b", sizeBits)
	for s := headerHeader; isSetMSB(s); {
		s = data[cursor]
		cursor++
		sizeBytes = fmt.Sprintf("%07b", s&127) + sizeBytes
	}
	if size, err = strconv.ParseInt(sizeBytes, 2, 64); err != nil {
		panicErrf("Err parsing size: %v. Could not determine size for %s", err, v.repr)
	}
	pot := PackedObjectType(typeBits)
	var bytes []byte
	if i+1 < len(entries) {
		n := int(entryLen)
		bytes = data[cursor:n]
		cursor = n
	} else {
		bytes = data[cursor:]
	}
	switch {
	case pot == PackedBlob || pot == PackedCommit || pot == PackedTree || pot == PackedTag:
		obj = parseNonDeltaEntry(&bytes, pot, v.ObjectId, size)
	default:
		var (
			deltaDeflated packedDelta
			base          *PackedObject
			baseOffset    int64
			dp            *packedObjectParser
		)
		switch pot {
		case ObjectRefDelta:
			var oid *ObjectId
			deltaDeflated, oid = readPackedRefDelta(&bytes)
			baseOffset = idx.entriesById[oid.String()].offset
		case ObjectOffsetDelta:
			if deltaDeflated, baseOffset, err = readPackedOffsetDelta(&bytes); err != nil {
				panicErrf("Err parsing size: %v. Could not determine size for %s", err, v.repr)
			}
			baseOffset = v.offset - baseOffset
		default:
			fmt.Printf("Unrecognized pack object type: %b ", pot)
		}
		objectIndex := sort.Search(len(idx.entries), func(i int) bool {
			return idx.entries[i].offset >= int64(baseOffset)
		})
		if idx.entries[objectIndex].offset != baseOffset {
			panicErrf("Could not find object with offset %d (%d - %d). Closest match was %d.", baseOffset,
				v.offset+baseOffset, v.offset, idx.entries[i].offset)
		}
		if objects[objectIndex] == nil {
			objects[objectIndex] = parseEntry(file, objectIndex, idx, objects)
			if objects[objectIndex] == nil {
				panicErrf("Ref deltas not yet implemented!")
			}
		}
		base = objects[objectIndex]
		bytes = *((*[]byte)(deltaDeflated))
		if dp, err = newPackedObjectParser(&bytes, v.ObjectId); err != nil {
			panicErr(err.Error())
		}
		obj = dp.parseDelta(base, v.ObjectId)
	}
	return
}

func parseNonDeltaEntry(bytes *[]byte, pot PackedObjectType, oid *ObjectId, size int64) (po *PackedObject) {
	var (
		dp  *packedObjectParser
		err error
	)
	if dp, err = newPackedObjectParser(bytes, oid); err != nil {
		panicErr(err.Error())
	}
	switch pot {
	case PackedBlob:
		po = dp.parseBlob(size)
	case PackedCommit:
		po = dp.parseCommit(size)
	case PackedTree:
		po = dp.parseTree(size)
	case PackedTag:
		po = dp.parseTag(size)
	}
	return
}

func (dp *packedObjectParser) parseCommit(size int64) *PackedObject {
	dp.hdr = &objectHeader{
		ObjectCommit,
		size,
	}
	commit := dp.objectParser.parseCommit()

	return &PackedObject{
		commit,
		dp.bytes,
		nil,
	}
}
func (dp *packedObjectParser) parseTag(size int64) *PackedObject {
	dp.hdr = &objectHeader{
		ObjectTag,
		size,
	}
	tag := dp.objectParser.parseTag()
	return &PackedObject{
		tag,
		dp.bytes,
		nil,
	}
}

func (dp *packedObjectParser) parseBlob(size int64) *PackedObject {
	blob := new(Blob)
	blob.data = dp.Bytes()
	blob.oid = dp.objectParser.oid
	blob.hdr = &objectHeader{
		ObjectBlob,
		size,
	}
	return &PackedObject{
		blob,
		blob.data,
		nil,
	}
}

func (dp *packedObjectParser) parseTree(size int64) *PackedObject {
	dp.hdr = &objectHeader{
		ObjectTree,
		size,
	}
	tree := dp.objectParser.parseTree()
	return &PackedObject{
		tree,
		dp.bytes,
		nil,
	}
}

// ================================================================= //
// Delta parsing.
// ================================================================= //

type packedDelta *[]byte

func readPackedRefDelta(bytes *[]byte) (delta packedDelta, oid *ObjectId) {
	baseOidBytes := (*bytes)[0:20]
	deltaBytes := (*bytes)[20:]
	delta = packedDelta(&deltaBytes)
	oid, _ = OidFromBytes(baseOidBytes)
	return packedDelta(&deltaBytes), oid
}

func readPackedOffsetDelta(bytes *[]byte) (delta packedDelta, offset int64, err error) {
	//first the offset to the base object earlier in the pack
	var i int
	offset, err, i = parseOffset(bytes)
	//now the rest of the bytes - the compressed delta
	deltaBytes := (*bytes)[i:]
	delta = packedDelta(&deltaBytes)
	return
}

func parseOffset(bytes *[]byte) (offset int64, err error, index int) {
	offsetBits := ""
	var base int64
	for i := 0; ; {
		v := (*bytes)[i]
		offsetBits = offsetBits + fmt.Sprintf("%07b", v&127)
		if i >= 1 {
			base += int64(1 << (7 * uint(i)))
		}
		if !isSetMSB(v) {
			if offset, err = strconv.ParseInt(offsetBits, 2, 64); err != nil {
				return
			}
			offset += base
			index = i + 1
			break
		}
		i++
	}
	return
}

func (p *objectParser) readByteAsInt() int64 {
	return int64(p.ReadByte())
}

func (dp *packedObjectParser) parseDelta(base *PackedObject, id *ObjectId) (object *PackedObject) {
	p := dp.objectParser

	baseSize := p.parseIntWhileMSB()
	outputSize := p.parseIntWhileMSB()

	src := base.bytes

	if int(baseSize) != len(src) {
		panicErrf("Expected size of base object is %d, but actual size is %d")
	}

	out := make([]byte, 0)
	var appended int64
	cmd := p.ReadByte()
	for {
		if cmd == 0 {
			panicErrf("Invalid delta! Byte 0 is not a supported delta code.")
		}
		var offset, len int64
		if cmd&0x80 != 0 {
			//copy from base to output
			offset, len = dp.parseCopyCmd(cmd)
			for i := offset; i < offset+len; i++ {
				out = append(out, (src)[i])
			}
			if offset+len > baseSize {
				panicErrf("Bad delta - references byte %d of a %d-byte source", offset+len, baseSize)
				break
			}
		} else {
			//copy from delta to output
			offset, len = 0, int64(cmd)
			for i := offset; i < offset+len; i++ {
				out = append(out, p.ReadByte())
			}
		}
		appended += len
		if appended < outputSize {
			cmd = p.ReadByte()
		} else {
			break
		}
	}
	if appended != outputSize {
		panicErrf("Expected output of size %d, got %d. \n", outputSize, appended)
	}
	if outputSize != int64(len(out)) {
		panicErrf("Expected output of len %d, got %d. \n", outputSize, len(out))
	}
	outputType := base.Object.Header().Type()
	outputParser := newObjectParser(bufio.NewReader(bytes.NewReader(out)), id)
	outputParser.hdr = &objectHeader{
		outputType,
		outputSize,
	}
	var obj Object
	switch outputType {
	case ObjectBlob:
		obj = outputParser.parseBlob()
	case ObjectTree:
		obj = outputParser.parseTree()
	case ObjectCommit:
		obj = outputParser.parseCommit()
	case ObjectTag:
		obj = outputParser.parseTag()
	}
	return &PackedObject{
		obj,
		out,
		base.Object.ObjectId(),
	}
}

// Given a copy command, return the offset and length it represents. None or all of
// the next seven bytes may be read, as determined by the seven least significant
// bits of the copy command.
func (dp *packedObjectParser) parseCopyCmd(cmd byte) (offset int64, len int64) {
	p := dp.objectParser
	offset, len = 0, 0
	if cmd&0x01 != 0 {
		offset = p.readByteAsInt()
	}
	if cmd&0x02 != 0 {
		offset |= (p.readByteAsInt() << 8)
	}
	if cmd&0x04 != 0 {
		offset |= (p.readByteAsInt() << 16)
	}
	if cmd&0x08 != 0 {
		offset |= (p.readByteAsInt() << 24)
	}
	if cmd&0x10 != 0 {
		len = p.readByteAsInt()
	}
	if cmd&0x20 != 0 {
		len |= (p.readByteAsInt() << 8)
	}
	if cmd&0x40 != 0 {
		len |= (p.readByteAsInt() << 16)
	}
	if len == 0 {
		len = 0x10000
	}
	return
}

// Compute an integer value in the format that pack files use for a delta's base
// size and output size. The function is named after the decoding mechanism:
// bytes are read and computed until a byte is found whose most significant
// bit is not set.
func (p *objectParser) parseIntWhileMSB() (i int64) {
	n := 0
	for {
		v := p.ReadByte()
		i |= (int64(v&127) << (uint(n) * 7))
		if !isSetMSB(v) {
			break
		}
		n++
	}
	return i
}

// ================================================================= //
// UTIL METHODS
// ================================================================= //

//return true if the most significant bit is set, false otherwise
func isSetMSB(b byte) bool {
	return b > 127
}

func packName(fileName string) string {
	return fileName[5 : len(fileName)-4]
}
