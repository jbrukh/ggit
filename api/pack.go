//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//

package api

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/api/parse"
	"github.com/jbrukh/ggit/util"
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
	object  objects.Object
	bytes   []byte
	DeltaOf *objects.ObjectId
	//the length of this object's delta chain. 0 for non-delta objects.
	Depth int
}

func (p *PackedObject) Object() objects.Object {
	return p.object
}

type Pack struct {
	// Git currently accepts version number 2 or 3 but
	// generates version 2 only.
	version int32
	// the unpacked objects sorted by offset
	content []*PackedObject
	idx     *Idx
	name    string
	opener  opener
	file    *os.File
}

type Idx struct {
	// the object ids sorted by offset
	entries []*PackedObjectId
	// the object ids in order of oid
	entriesById []*PackedObjectId
	// caches oid lookup results
	idToEntry map[string]*PackedObjectId
	// the fan-out counts - each value represents the
	// number of objects in this pack whose 1st byte
	// is >= the index of that value.
	counts *[256]int
	// number of objects contained in the pack). equal to
	// counts[255].
	count int64
	// copy of the checksum for this idx file's
	// corresponding pack file.
	packChecksum *objects.ObjectId
	// checksum for this idx file.
	idxChecksum *objects.ObjectId
}

type PackedObjectId struct {
	*objects.ObjectId
	offset int64
	crc32  int64
	index  int
}

type opener func() (*os.File, error)

func (pack *Pack) open() error {
	if pack.file != nil {
		//already open
		return nil
	}
	if file, err := pack.opener(); err != nil {
		return err
	} else {
		pack.file = file
	}
	return nil
}

// close will nil-ify and close the 
// pack file resource, but not in that
// order
func (pack *Pack) close() (err error) {
	file := pack.file
	pack.file = nil
	if file != nil {
		err = file.Close()
	}
	return
}

// ================================================================= //
// OBJECT RETRIEVAL
// ================================================================= //

func (idx *Idx) entriesWithPrefix(prefix byte) []*PackedObjectId {
	//the object we want is somewhere in the range between from (inclusive) and to (exclusive).
	var from, to int
	if prefix == 0 {
		from = 0
	} else {
		from = idx.counts[prefix-1]
	}
	to = idx.counts[prefix]
	if from == to {
		return nil
	}
	return idx.entriesById[from:to]
}

func (idx *Idx) entryById(oid *objects.ObjectId) *PackedObjectId {
	trimmed := idx.entriesWithPrefix(oid.Bytes()[0])
	if trimmed == nil {
		return nil
	}
	id := oid.String()
	if idx.idToEntry[id] != nil {
		return idx.idToEntry[id]
	}
	gte := func(i int) bool {
		var oid *objects.ObjectId
		oid = trimmed[i].ObjectId
		return oid.String() >= id
	}
	i := sort.Search(len(trimmed), gte)
	if i >= len(trimmed) {
		return nil
	}
	result := trimmed[i]
	if result.ObjectId.String() != id {
		return nil
	}
	idx.idToEntry[id] = result
	return result
}

// Returns the one Object in this pack with the given ObjectId,
// or nil, NoSuchObject if no such Object is in this pack.
func (pack *Pack) unpack(oid *objects.ObjectId) (obj objects.Object, result packSearch) {
	defer pack.close()
	if entry := pack.idx.entryById(oid); entry != nil {
		if pack.content[entry.index] == nil {
			pack.content[entry.index] = pack.parseEntry(entry.index)
		}
		obj, result = pack.content[entry.index].object, OneSuchObject
	}
	return
}

type packSearch byte

const (
	NoSuchObject    packSearch = 0
	OneSuchObject   packSearch = 1
	MultipleObjects packSearch = 2
)

func (pack *Pack) unpackFromShortOid(short string) (obj objects.Object, result packSearch) {
	prefix, err := strconv.ParseUint(short[0:2], 16, 8)
	if err != nil {
		util.PanicErrf("invalid short oid; non-hex characters: %s. %s", short, err.Error())
	}
	entries := pack.idx.entriesWithPrefix(byte(prefix))
	if entries == nil {
		return
	}
	var already bool
	for _, oid := range entries {
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
func unpackFromShortOid(packs []*Pack, short string) (obj objects.Object, ok bool) {
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

func unpack(packs []*Pack, oid *objects.ObjectId) (obj objects.Object, ok bool) {
	var result packSearch
	for _, pack := range packs {
		if obj, result = pack.unpack(oid); result == OneSuchObject {
			//trust for now that there will only be one matching object among the packs.
			return obj, true
		}
	}
	return
}

func objectIdsFromPacks(packs []*Pack) (ids []*objects.ObjectId) {
	var count int64
	for _, pack := range packs {
		count += pack.idx.count
	}
	ids = make([]*objects.ObjectId, count, count)
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
		defer pack.close()
		for j := range pack.idx.entries {
			pack.content[j] = pack.parseEntry(j)
			objects[i] = pack.content[j]
			i++
		}
	}
	return objects
}

// ================================================================= //
// GGIT PACK PARSER
// ================================================================= //

type packIdxParser struct {
	idxParser  *parse.ObjectIdParser
	name       string
	packOpener opener
}

func newPackIdxParser(idx *bufio.Reader, packOpener opener, name string) *packIdxParser {
	oidParser := parse.NewObjectIdParser(idx)
	return &packIdxParser{
		idxParser:  oidParser,
		name:       name,
		packOpener: packOpener,
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
	var counts [256]int
	for i := range counts {
		counts[i] = int(p.idxParser.ParseIntBigEndian(4))
	}
	//discard the fan-out values, just use the largest value,
	//which is the total # of objects:
	count := counts[255]
	idToEntry := make(map[string]*PackedObjectId)
	entries := make([]*PackedObjectId, count, count)
	entriesByOid := make([]*PackedObjectId, count, count)
	for i := 0; i < count; i++ {
		b := p.idxParser.ReadNBytes(20)
		oid, _ := objects.OidFromBytes(b)
		entries[i] = &PackedObjectId{
			ObjectId: oid,
		}
		entriesByOid[i] = entries[i]
	}
	for i := 0; i < count; i++ {
		entries[i].crc32 = int64(p.idxParser.ParseIntBigEndian(4))
	}
	for i := 0; i < count; i++ {
		//TODO: 8-byte #'s for some offsets for some pack files (packs > 2gb)
		entries[i].offset = p.idxParser.ParseIntBigEndian(4)
	}
	checksumPack := p.idxParser.ReadNBytes(20)
	checksumIdx := p.idxParser.ReadNBytes(20)
	if !p.idxParser.EOF() {
		util.PanicErrf("Found extraneous bytes! %x", p.idxParser.Bytes())
	}
	//order by offset
	sort.Sort(packedObjectIds(entries))
	for i, v := range entries {
		v.index = i
	}
	packChecksum, _ := objects.OidFromBytes(checksumPack)
	idxChecksum, _ := objects.OidFromBytes(checksumIdx)
	return &Idx{
		entries,
		entriesByOid,
		idToEntry,
		&counts,
		int64(count),
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

func newPackedObjectParser(data []byte, oid *objects.ObjectId) (p *packedObjectParser, e error) {
	compressedReader := bytes.NewReader(data)
	var zr io.ReadCloser
	if zr, e = zlib.NewReader(compressedReader); e == nil {
		defer zr.Close()
		exploder := util.NewDataParser(bufio.NewReader(zr))
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

//parse the pack's meta data and close it
func (p *packIdxParser) parsePack() *Pack {
	//parse the index and construct the pack
	idx := p.parseIdx()
	objects := make([]*PackedObject, idx.count)
	pack := &Pack{
		PackVersion,
		objects,
		idx,
		p.name,
		p.packOpener,
		nil,
	}
	//verify the pack file
	if err := pack.open(); err != nil {
		util.PanicErrf("Could not open pack file %s: %s", pack.name, err)
	}
	dataParser := util.NewDataParser(bufio.NewReader(pack.file))
	dataParser.ConsumeString(PackSignature)
	dataParser.ConsumeBytes([]byte{0, 0, 0, PackVersion})
	count := dataParser.ParseIntBigEndian(4)
	if count != idx.count {
		util.PanicErrf("Pack file count doesn't match idx file count for pack-%s!", p.name) //todo: don't panic.
	}
	pack.close()
	return pack
}

// parse the ith entry of this pack, opening the pack resource if necessary
func (p *Pack) parseEntry(i int) (obj *PackedObject) {
	if len(p.content) > i && p.content[i] != nil {
		return p.content[i] //already parsed
	}
	if err := p.open(); err != nil {
		util.PanicErrf("Could not open pack file %s: %s", p.name, err.Error())
	}
	size, pot, bytes := p.entrySizeTypeData(i)
	e := p.idx.entries[i]
	switch {
	case pot == PackedBlob || pot == PackedCommit || pot == PackedTree || pot == PackedTag:
		obj = parseNonDeltaEntry(bytes, pot, e.ObjectId, int64(size))
	case pot == ObjectOffsetDelta || pot == ObjectRefDelta:
		obj = p.parseDeltaEntry(bytes, pot, e.ObjectId, i)
	default:
		util.PanicErrf("Unrecognized object type %d in pack %s for entry with id %s", pot, p.name, e.ObjectId)
	}
	return
}

// extract the size, type, and compressed data of the ith object of the pack file
func (p *Pack) entrySizeTypeData(i int) (uint64, PackedObjectType, []byte) {
	data := p.readEntry(i)
	// keep track of bytes read so that, in conjunction with the next entry's offset, we can know where the next
	// object in the pack begins.
	headerHeader := data[0]
	read := 1
	typeBits := (headerHeader & 127) >> 4
	sizeBits := (headerHeader & 15)
	//collect remaining size bytes, if any.
	size := uint64(0)
	for s := headerHeader; isSetMSB(s); {
		s = data[read]
		size |= uint64(s&127) << uint64((read-1)*7)
		read++
	}
	size = (size << 4) + uint64(sizeBits)
	return size, PackedObjectType(typeBits), data[read:]
}

// read all the bytes of the ith object of the pack file
func (p *Pack) readEntry(i int) []byte {
	e := p.idx.entries[i]
	var size int64
	if i+1 < len(p.idx.entries) {
		size = p.idx.entries[i+1].offset - e.offset
	} else {
		if info, err := p.file.Stat(); err != nil {
			util.PanicErrf("Could not determine size of pack file %s: %s", p.file.Name(), err)
		} else {
			size = info.Size() - e.offset
		}
	}
	data := make([]byte, size, size)
	if _, err := p.file.ReadAt(data, e.offset); err != nil {
		util.PanicErrf("Could not read %d bytes from %d of pack file %s: %s", len(data), e.offset, p.file.Name(), err)
	}
	return data
}

func parseNonDeltaEntry(bytes []byte, pot PackedObjectType, oid *objects.ObjectId, size int64) (po *PackedObject) {
	var (
		dp  *packedObjectParser
		err error
	)
	if dp, err = newPackedObjectParser(bytes, oid); err != nil {
		util.PanicErr(err.Error())
	} else if int64(len(dp.bytes)) != size {
		util.PanicErrf("Expected object of %d bytes but found %d bytes", size, len(dp.bytes))
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
	dp.hdr = objects.NewObjectHeader(objects.ObjectCommit, size)
	commit := dp.objectParser.parseCommit()

	return &PackedObject{
		object: commit,
		bytes:  dp.bytes,
	}
}
func (dp *packedObjectParser) parseTag(size int64) *PackedObject {
	dp.hdr = objects.NewObjectHeader(objects.ObjectTag, size)
	tag := dp.objectParser.parseTag()
	return &PackedObject{
		object: tag,
		bytes:  dp.bytes,
	}
}

func (dp *packedObjectParser) parseBlob(size int64) *PackedObject {
	data := dp.Bytes()
	oid := dp.objectParser.oid
	hdr := objects.NewObjectHeader(objects.ObjectBlob, size)
	blob := objects.NewBlob(oid, hdr, data)
	return &PackedObject{
		object: blob,
		bytes:  data,
	}
}

func (dp *packedObjectParser) parseTree(size int64) *PackedObject {
	dp.hdr = objects.NewObjectHeader(objects.ObjectTree, size)
	tree := dp.objectParser.parseTree()
	return &PackedObject{
		object: tree,
		bytes:  dp.bytes,
	}
}

// ================================================================= //
// Delta parsing.
// ================================================================= //

type packedDelta []byte

func (p *Pack) parseDeltaEntry(bytes []byte, pot PackedObjectType, oid *objects.ObjectId, i int) *PackedObject {
	var (
		deltaDeflated packedDelta
		baseOffset    int64
		dp            *packedObjectParser
		err           error
	)
	e := p.idx.entries[i]
	switch pot {
	case ObjectRefDelta:
		var oid *objects.ObjectId
		deltaDeflated, oid = readPackedRefDelta(bytes)
		e := p.idx.entryById(oid)
		if e == nil {
			util.PanicErrf("nil entry for base object with id %s", oid.String())
		}
		baseOffset = e.offset
	case ObjectOffsetDelta:
		if deltaDeflated, baseOffset, err = readPackedOffsetDelta(bytes); err != nil {
			util.PanicErrf("Err parsing size: %v. Could not determine size for %s", err, e.String())
		}
		baseOffset = e.offset - baseOffset
	}
	base := p.findObjectByOffset(baseOffset)
	bytes = []byte(deltaDeflated)
	if dp, err = newPackedObjectParser(bytes, oid); err != nil {
		util.PanicErr(err.Error())
	}
	return dp.applyDelta(base, oid)
}

func readPackedRefDelta(bytes []byte) (delta packedDelta, oid *objects.ObjectId) {
	baseOidBytes := bytes[0:20]
	deltaBytes := bytes[20:]
	delta = packedDelta(deltaBytes)
	oid, _ = objects.OidFromBytes(baseOidBytes)
	return
}

func readPackedOffsetDelta(bytes []byte) (delta packedDelta, offset int64, err error) {
	//first the offset to the base object earlier in the pack
	var i int
	offset, err, i = parseOffset(bytes)
	//now the rest of the bytes - the compressed delta
	deltaBytes := bytes[i:]
	delta = packedDelta(deltaBytes)
	return
}

func parseOffset(bytes []byte) (offset int64, err error, index int) {
	offsetBits := ""
	var base int64
	for i := 0; ; {
		v := bytes[i]
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

func (p *Pack) findObjectByOffset(offset int64) *PackedObject {
	i := sort.Search(len(p.idx.entries), func(j int) bool {
		return p.idx.entries[j].offset >= int64(offset)
	})
	if p.idx.entries[i].offset != offset {
		util.PanicErrf("Could not find object with offset %d. Closest match was %d.", offset, i)
	}
	if p.content[i] == nil {
		p.content[i] = p.parseEntry(i)
	}
	if p.content[i] == nil {
		util.PanicErrf("Could not find or parse object with offset %d", offset)
	}
	return p.content[i]
}

func (p *objectParser) readByteAsInt() int64 {
	return int64(p.ReadByte())
}

func (dp *packedObjectParser) applyDelta(base *PackedObject, id *objects.ObjectId) (object *PackedObject) {
	p := dp.objectParser

	baseSize := p.parseIntWhileMSB()
	outputSize := p.parseIntWhileMSB()

	src := base.bytes

	if int(baseSize) != len(src) {
		util.PanicErrf("Expected size of base object is %d, but actual size is %d")
	}

	out := make([]byte, outputSize, outputSize)
	var appended int64
	cmd := p.ReadByte()
	for {
		if cmd == 0 {
			util.PanicErrf("Invalid delta! Byte 0 is not a supported delta code.")
		}
		var offset, len int64
		if cmd&0x80 != 0 {
			//copy from base to output
			offset, len = dp.parseCopyCmd(cmd)
			for i := offset; i < offset+len; i++ {
				out[appended+(i-offset)] = src[i]
			}
			if offset+len > baseSize {
				util.PanicErrf("Bad delta - references byte %d of a %d-byte source", offset+len, baseSize)
				break
			}
		} else {
			//copy from delta to output
			offset, len = 0, int64(cmd)
			for i := offset; i < offset+len; i++ {
				out[appended+(i-offset)] = p.ReadByte()
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
		util.PanicErrf("Expected output of size %d, got %d. \n", outputSize, appended)
	}
	if outputSize != int64(len(out)) {
		util.PanicErrf("Expected output of len %d, got %d. \n", outputSize, len(out))
	}
	outputType := base.object.Header().Type()
	outputParser := newObjectParser(bufio.NewReader(bytes.NewReader(out)), id)
	outputParser.hdr = objects.NewObjectHeader(outputType, outputSize)
	var obj objects.Object
	switch outputType {
	case objects.ObjectBlob:
		obj = outputParser.parseBlob()
	case objects.ObjectTree:
		obj = outputParser.parseTree()
	case objects.ObjectCommit:
		obj = outputParser.parseCommit()
	case objects.ObjectTag:
		obj = outputParser.parseTag()
	}
	return &PackedObject{
		obj,
		out,
		base.object.ObjectId(),
		base.Depth + 1,
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
