package ggit

import (
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
    const HEADER_FMT = "IndexHeader{Sig=%s, Version=%d, Count=%d}"
    return fmt.Sprintf(HEADER_FMT, hdr.Sig, hdr.Version, hdr.Count)
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
}

func ParseIndexFile(repo *Repository) (err error) {
    file, err := repo.IndexFile()
    if err != nil {
        return
    }
    return parseIndex(file)
}

func parseIndex(file *os.File) (err error) {
    const SIGNATURE = "DIRC"

    var hdr indexHeader
    if err = binary.Read(file, ord, hdr); err != nil {
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
    fmt.Println(hdr)

    // read the entries
    var i int32
    for i = 0; i < hdr.Count; i++ {
    }
    return
}
