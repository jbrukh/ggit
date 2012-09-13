package ggit

import (
    "bufio"
    "encoding/binary"
    "errors"
    "fmt"
    "os"
    "time"
)

// all serialized data is stored in network
// byte order as per the specification
var ord binary.ByteOrder = binary.BigEndian

// the deserialized version of the 
// git index file
type Index struct {
    version int32
    entries []*IndexEntry
    // TODO: extentions
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
func (inx *Index) Entries() {
    // TODO
}

// Extentions will visit and/or return the
// index file extentions.
func (inx *Index) Extentions() {
    // TODO
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

type IndexEntry struct {
    eid   *[20]byte // TODO: is this an object id, or just a SHA??
    flags EntryFlagsV2
    name  string
}

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
func (i *statInfo) CTime() time.Time {
    return time.Unix(int64(i.CTimeSecs), int64(i.CTimeNanos))
}

// MTime returns the last time file metadata has changed
func (i *statInfo) MTime() time.Time {
    return time.Unix(int64(i.MTimeSecs), int64(i.MTimeNanos))
}

// index entry version 2
type indexEntry struct {
    StatInfo statInfo
    Sha1     [20]byte
    Flags    int16 // TODO: this should be wrapped in the appropriate EntryFlags
}

func (i *indexEntry) String() string {
    const INDEX_ENTRY_FMT = "IndexEntry{" +
        "CTime=%v, " +
        "MTime=%v, " +
        "Dev=%d, " +
        "Ino=%d, " +
        "Mode=%o, " +
        "Uid=%d, " +
        "Gid=%d, " +
        "Size=%d, " +
        "SHA1=%s, " +
        "Flags=%d}"
    sha := NewObjectIdFromArray(i.Sha1)
    return fmt.Sprintf(
        INDEX_ENTRY_FMT,
        i.StatInfo.CTime(),
        i.StatInfo.MTime(),
        i.StatInfo.Dev,
        i.StatInfo.Ino,
        i.StatInfo.Mode,
        i.StatInfo.Uid,
        i.StatInfo.Gid,
        i.StatInfo.Size,
        sha,
        i.Flags,
    )

}

func ParseIndexFile(repo Repository) (err error) {
    file, err := repo.IndexFile()
    if err != nil {
        return
    }
    defer file.Close()
    return parseIndex(file)
}

func parseIndex(f *os.File) (err error) {
    const SIGNATURE = "DIRC"

    file := bufio.NewReader(f)

    var hdr indexHeader
    if err = binary.Read(file, ord, &hdr); err != nil {
        return
    }
    sig := string(hdr.Sig[:])
    if string(sig) != SIGNATURE {
        return errors.New("wrong signature")
    }
    if hdr.Version != 2 {
        return errors.New("unsupported index file format version")
    }
    if hdr.Count < 0 {
        return errors.New("header count is off")
    }
    fmt.Printf("%s\n", hdr.String())

    // read the entries
    var i int32
    for i = 0; i < hdr.Count; i++ {
        var entry indexEntry
        err = binary.Read(file, ord, &entry)
        if err != nil {
            return err
        }
        fmt.Println(entry.String())
        // TODO: what if it is corrupted and too long?
        line, err := file.ReadBytes('\000')
        if err != nil {
            return err
        }

        line = line[:len(line)-1] // get rid of NUL
        fmt.Println(string(line))

        // don't ask me how I figured this out after
        // a 14 hour workday
        leftOver := 7 - (len(line)+6)%8
        for j := 0; j < leftOver; j++ {
            if _, err = file.ReadByte(); err != nil {
                return err
            }
        }
    }
    return
}
