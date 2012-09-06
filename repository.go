package ggit

import (
	"path"
	"os"
	"io"
	"bytes"
	"strings"
	"errors"
	"strconv"
	"compress/zlib"
)

type Repository struct {
	path string
}

func OpenRepository(path string) (r *Repository, err error) {
	// check that repo is valid
	r = &Repository{
		path: path,
	}
	return
}

func (r *Repository) ReadBlob(oid *ObjectId) (b *Blob, err error) {
	var file *os.File
	filePath := path.Join(r.path, objectPath(oid))

	file, err = os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	
	buf := new(bytes.Buffer)
	zr, err := zlib.NewReader(file)
	defer zr.Close()
	
	if _, err = io.Copy(buf, zr); err != nil {
		return
	}
	header, err := buf.ReadString('\000')
	
	// remove last character
	header = header[:len(header)-1]
	splt := strings.Split(header, " ")
	var l int
	l, err = strconv.Atoi(splt[1])
	if splt[0] != "blob" || err != nil {
		return nil, errors.New("Not a blob")
	}
	
	bts := buf.Bytes()
	if len(bts) != l {
		return nil, errors.New("Size mismatch")
	}
	
	b = &Blob{
		bytes: bts,
	}
	return
}

func objectPath(oid *ObjectId) string {
	hex := oid.String()
	return path.Join("objects", hex[0:2], hex[2:])
}