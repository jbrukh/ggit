package ggit

import (
    "bufio"
    "encoding/binary"
    "errors"
    "fmt"
    "os"
)

// all serialized data is stored in network
// byte order as per the specification
var ord binary.ByteOrder = binary.BigEndian

// the deserialized version of the 
// git index file
type Index struct {
    version int32
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

// index entry version 2
type indexEntry struct {
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
    Sha1       [20]byte
    Flags      int16
    // bytes mod 8 == 4 here
}

func (hdr *indexEntry) String() string {
    const INDEX_ENTRY_FMT = "IndexEntry{CTimeSecs=%d, CTimeNanos=%d, MTimeSecs=%d, MTimeNanos=%d, Dev=%d, Ino=%d, Mode=%o, Uid=%d, Gid=%d, Size=%d, SHA1=%s, Flags=%d}"
    sha := NewObjectIdFromArray(hdr.Sha1)
    return fmt.Sprintf(INDEX_ENTRY_FMT, hdr.CTimeSecs, hdr.CTimeNanos, hdr.MTimeSecs, hdr.MTimeNanos, hdr.Dev, hdr.Ino, hdr.Mode, hdr.Uid, hdr.Gid, hdr.Size, sha, hdr.Flags)
}

func ParseIndexFile(repo *Repository) (err error) {
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
        fmt.Printf("read %d: %s\n", len(line), string(line))

        // don't ask me how I figured this out afte
        // a 14 hour workday
        leftOver := 8 - (len(line)+6)%8 - 1
        for j := 0; j < leftOver; j++ {
            if _, err = file.ReadByte(); err != nil {
                return err
            }
        }
    }
    return
}
