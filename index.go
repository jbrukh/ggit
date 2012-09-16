package ggit

import (
    "bufio"
    "bytes"
    "encoding/binary"
    "errors"
    "fmt"
    "time"
)

// all serialized data is stored in network
// byte order as per the specification
var ord binary.ByteOrder = binary.BigEndian

// ================================================================= //
// CONSTANTS
// ================================================================= //

// index files use 4-letter codes as signatures
// to identify the format and extention
type signature string

// the Git index must begin with this signature code
const SIG_INDEX_FILE signature = "DIRC"

// extention signatures
const (
    SIG_CACHED_TREE  signature = "TREE"
    SIG_RESOLVE_UNDO signature = "REUC"
)

// ================================================================= //
// PUBLIC GGIT INDEX OBJECTS
// ================================================================= //

// the ggit Index object
type Index struct {
    version    int32
    entries    []*IndexEntry
    extentions []*IndexExtention
}

// Version returns the version of the index
// file. Currently I am supporting version 2.
// TODO: support 3, 4, ...
func (inx *Index) Version() int32 {
    return inx.version
}

// TODO: how shall we traverse entries? We
// can possibly use a visitor, depending on
// what makes sense.
func (inx *Index) Entries() []*IndexEntry {
    return inx.entries
}

// Extentions will visit and/or return the
// index file extentions.
func (inx *Index) Extentions() []*IndexExtention {
    return inx.extentions
}

func (inx *Index) String() string {
    buf := bytes.NewBufferString("")
    fmt.Fprintf(buf, "Index (v.%d)\n", inx.version)

    if inx.entries != nil {
        for _, entry := range inx.entries {
            buf.WriteString(entry.String())
            buf.WriteString("\n")
        }
    } else {
        buf.WriteString("(no entries)\n")
    }
    if inx.extentions != nil {
        for _, ext := range inx.extentions {
            buf.WriteString(ext.String())
            buf.WriteString("\n")
        }
    } else {
        buf.WriteString("(no extentions)\n")
    }
    return buf.String()
}

type IndexEntry struct {
    eid   *ObjectId // TODO: is this an object id, or just a SHA??
    flags EntryFlagsV2
    name  string
    info  *statInfo
}

func (entry *IndexEntry) String() string {
    return fmt.Sprint(entry.eid.String(), " ", entry.info.String(), " ", entry.name)
}

type IndexExtention struct {
    etype signature
}

func (ext *IndexExtention) String() string {
    return ""
}

// EntryFlags for version 2
type EntryFlagsV2 int16

// TODO: document
func (f *EntryFlagsV2) AssumeValid() bool {
    return false
}

func (f *EntryFlagsV2) Extended() bool {
    return false // should be false for version 2
}

// TODO: this returns some two-bit result
// not yet clear what it is for
func (f *EntryFlagsV2) Stage() {
    // TODO
}

// 12-bit name length if less than 0xFFF, and
// 0xFFF otherwise
func (f *EntryFlagsV2) NameLength() int {
    // TODO
    return 0
}

// ================================================================= //
// INTERNAL REPRESENTATIONS OF BINARY DATA
// ================================================================= //

//
// For reference, see:
// 
// http://git.kernel.org/?p=git/git.git;a=blob;f=Documentation/technical/index-format.txt;hb=HEAD
//
type indexHeader struct {
    Sig     [4]byte
    Version int32
    Count   int32
}

func (hdr *indexHeader) String() string {
    const HEADER_FMT = "IndexHeader{Sig=%q, Version=%d, Count=%d}"
    return fmt.Sprintf(HEADER_FMT, string(hdr.Sig[:]), hdr.Version, hdr.Count)
}

// index entry version 2
type indexEntry struct {
    Info  statInfo
    Sha1  [20]byte
    Flags int16 // TODO: this should be wrapped in the appropriate EntryFlags
}

// the header of an index extention
type indexExtentionHeader struct {
    Sig   [4]byte
    Count int32
}

// data returned from stat, used by git
// to detect when a file is changed. It appears (according to some docs)
// that the particular kind of data is not as relevant as the fact that
// it changes.
type statInfo struct {
    CTimeSecs  int32
    CTimeNanos int32
    MTimeSecs  int32
    MTimeNanos int32
    Dev        int32
    Ino        int32
    Mode       int32
    Uid        int32
    Gid        int32
    Size       int32
}

// CTime returns the last time file metadata has changed
func (info *statInfo) CTime() time.Time {
    return time.Unix(int64(info.CTimeSecs), int64(info.CTimeNanos))
}

// MTime returns the last time file metadata has changed
func (info *statInfo) MTime() time.Time {
    return time.Unix(int64(info.MTimeSecs), int64(info.MTimeNanos))
}

func (info *statInfo) String() string {
    const FMT_STATINFO = "%v  %v  %d  %d  %o  %d  %d  %6d"
    return fmt.Sprintf(
        FMT_STATINFO,
        info.CTime(),
        info.MTime(),
        info.Dev,
        info.Ino,
        info.Mode,
        info.Uid,
        info.Gid,
        info.Size,
    )
}

// ================================================================= //
// PARSING FUNCTIONS
// ================================================================= //

func parseIndexHeader(r *bufio.Reader) (hdr *indexHeader, err error) {
    var h indexHeader
    if err = binary.Read(r, ord, &h); err != nil {
        return
    }
    sig := signature(h.Sig[:])
    if sig != SIG_INDEX_FILE || h.Version != 2 || h.Count < 0 {
        return nil, errors.New("bad header")
    }
    return &h, nil
}

func parseIndexEntry(r *bufio.Reader) (entry *IndexEntry, err error) {
    var binEntry indexEntry
    err = binary.Read(r, ord, &binEntry)
    if err != nil {
        return nil, err
    }

    // TODO: what if it is corrupted and too long?
    name, err := r.ReadBytes(NUL)
    if err != nil {
        return nil, err
    }

    name = name[:len(name)-1] // get rid of NUL

    // don't ask me how I figured this out after
    // a 14 hour workday
    leftOver := 7 - (len(name)+6)%8
    for j := 0; j < leftOver; j++ {
        // TODO: read the bytes at once somehow
        if _, err = r.ReadByte(); err != nil {
            return nil, err
        }
    }

    // record the entry
    return toIndexEntry(&binEntry, string(name)), nil
}

func parseIndexExtention(file *bufio.Reader) (ext *IndexExtention, err error) {
    return
}

// ================================================================= //
// CONVERSION FUNCTIONS
// ================================================================= //

// toIndex converts a reader pointing at a serialized
// index object into a ggit.Index object
func toIndex(file *bufio.Reader) (idx *Index, err error) {
    // first parse the header of the index and make
    // sure we are OK with the version and know the
    // index entry count
    hdr, e := parseIndexHeader(file)
    if e != nil {
        return nil, e
    }

    // initialize the index object and prepare for
    // populating the entries. note we may be doing
    // this in vain if entries are invalid, etc.
    idx = new(Index)
    idx.version = hdr.Version
    idx.entries = make([]*IndexEntry, 0, hdr.Count)

    // read the entries
    var i int32
    for i = 0; i < hdr.Count; i++ {
        entry, e := parseIndexEntry(file)
        if e != nil {
            return nil, e
        }
        idx.entries = append(idx.entries, entry)
    }

    /*  // read the extentions
        idx.extentions = make([]*IndexExtention, 4)
        for {
            ext, e := parseIndexExtention(file)
            if e != nil {
                return nil, e
            }
            idx.extentions = append(idx.extentions, ext)
        }
    */
    return
}

func toIndexEntry(entry *indexEntry, name string) *IndexEntry {
    return &IndexEntry{
        eid:   NewObjectIdFromArray(entry.Sha1),
        flags: EntryFlagsV2(entry.Flags),
        name:  name,
        info:  &entry.Info,
    }
}
