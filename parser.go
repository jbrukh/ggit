package ggit

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ================================================================= //
// PARSE ERROR TYPE & PANICS
// ================================================================= //

// ParseErr is a common error that occurs when ggit is 
// parsing binary objects
type parseErr string

// ParseErr is an error
func (p parseErr) Error() string {
	return string(p)
}

// parseErrf allows convenience formatting for ParseErrors
func parseErrf(format string, items ...interface{}) parseErr {
	return parseErr(fmt.Sprintf(format, items))
}

// parseErrn concatenates a bunch of items together to 
// form the error string
func parseErrn(items ...string) parseErr {
	return parseErr(strings.Join(items, ""))
}

func panicErr(msg string) {
	panic(parseErr(msg))
}

func panicErrn(items ...string) {
	panic(parseErrn(items...))
}

func panicErrf(format string, items ...interface{}) {
	panic(parseErrf(format, items))
}

// ================================================================= //
// HELPERS
// ================================================================= //

// trimLast throws away the last character of a byte slice
func trimLast(b []byte) []byte {
	if b == nil || len(b) == 0 {
		return b
	}
	return b[:len(b)-1]
}

func trimLastStr(b []byte) string {
	return string(trimLast(b))
}

// ================================================================= //
// DATA PARSER
// ================================================================= //

type dataParser struct {
	buf  *bufio.Reader
	read int // bytes read TODO make 64
}

// dataParse allows you to call a number of parsing functions on your
// parser at once, without having to handle errors explicitly. If an
// error occurs, the parser commands will panic with parseErr, which
// this method will recover and return
func dataParse(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(parseErr); ok {
				err = e
			}
		}
	}()
	f()
	return
}

// ================================================================= //
// DATA PARSING API
// ================================================================= //

func (p *dataParser) consume(n int) []byte {
	b := make([]byte, n)
	if rd, e := p.buf.Read(b); e != nil || rd != n {
		panicErrf("expected: %d byte(s)", n)
	}
	p.read += n
	return b
}

func (p *dataParser) peek(n int) (b []byte) {
	var e error
	if b, e = p.buf.Peek(n); e != nil || len(b) != n {
		panicErrf("expected: %d byte(s)", n)
	}
	return b
}

func (p *dataParser) consumeUntil(delim byte) []byte {
	b, e := p.buf.ReadBytes(delim)
	if e != nil {
		panicErrf("expected delimiter: %v", delim)
	}
	return trimLast(b)
}

// ResetRead resets the read byte count
func (p *dataParser) ResetRead() {
	p.read = 0
}

func (p *dataParser) EOF() bool {
	if _, e := p.buf.Peek(1); e != nil {
		if e == io.EOF {
			return true
		} else {
			panicErr("reading error")
		}
	}
	return false
}

// Consume will consume n bytes without regard for what the underlying
// data might be. If it is unable to consume, then a panic is raised
// with parseErr.
func (p *dataParser) Consume(n int) {
	p.consume(n)
}

// ConsumeByte will consume a single byte and compare it to b. If it
// does not match, or cannot be read, then a panic is raised with parseErr.
func (p *dataParser) ConsumeByte(b byte) {
	if p.consume(1)[0] != b {
		panicErrf("expected byte: %v", b)
	}
}

// PeekBute will return the next byte without advancing the reader. If
// it cannot be read, then a panic is raised with parseErr.
func (p *dataParser) PeekByte() byte {
	return p.peek(1)[0]
}

// PeekBytes
func (p *dataParser) PeekBytes(n int) []byte {
	return p.peek(n)
}

// ConsumeBytes will consume len(b) bytes and compare them to b. If they
// do not match, or cannot be read, then a panic is raised with parseErr.
func (p *dataParser) ConsumeBytes(b []byte) {
	d := p.consume(len(b))
	for inx, v := range d {
		if d[inx] != v {
			panicErrf("expected bytes: %v", b)
		}
	}
}

// ConsumeString will consume len(b) bytes and compare them to the string s. If they
// do not match, or cannot be read, then a panic is raised with parseErr.
func (p *dataParser) ConsumeString(s string) {
	b := p.consume(len(s))
	if string(b) != s {
		panicErrf("expected string: %s", s)
	}
}

// ConsumeStrings will check if the Reader contains any one of the provided strings
// and if so, consumes it and returns it. If a string match cannot be found, or 
// cannot be read, than a panic is raised with parseErr. (The first string matched
// is the one returned, such that if strings are substrings of one another, only
// the first match matters.)
func (p *dataParser) ConsumeStrings(s []string) string {
	for _, str := range s {
		l := len(str)
		pk := p.PeekString(l)
		if l != 0 && pk == str {
			return string(p.consume(l))
		}
	}
	panicErrf("expected one of: %v", s)
	return ""
}

// PeekString returns true if and only if the next bytes
// in the buffer match the given input string (the string
// in the buffer is NOT consumed)
func (p *dataParser) PeekString(n int) string {
	pk, e := p.buf.Peek(n)
	if e != nil {
		panicErrf("expected: %d byte(s)", n)
	}
	return string(pk)
}

// ReadBytesUntil
func (p *dataParser) ReadBytes(delim byte) []byte {
	return p.consumeUntil(delim)
}

// ReadStringUtil
func (p *dataParser) ReadString(delim byte) string {
	return string(p.consumeUntil(delim))
}

// String returns the entirety of the remaining data
// in the buffer, up to the EOF, as a string
func (p *dataParser) String() string {
	return string(p.Bytes())
}

// Bytes returns the entirety of the remaining data
// in the buffer, up to the EOF, as bytes
func (p *dataParser) Bytes() []byte {
	b := new(bytes.Buffer)
	_, e := io.Copy(b, p.buf)
	if e != nil {
		panicErr(e.Error())
	}
	bts := b.Bytes()
	p.read += len(bts)
	return bts
}

// ================================================================= //
// SPECIALIZED PARSING FUNCTIONS
// ================================================================= //

func parseInt(str string, base int, bitSize int) (i64 int64) {
	var e error
	if i64, e = strconv.ParseInt(str, base, bitSize); e != nil {
		panicErrf("cannot convert integer (base $d): %s", base, str)
	}
	return i64
}

// TODO: this should be smarter than delimiter
func (p *dataParser) ParseAtoi(delim byte) (n int) {
	return int(p.ParseInt(delim, 10, 0))
}

func (p *dataParser) ParseIntN(n int, base int, bitSize int) (i64 int64) {
	return parseInt(string(p.consume(n)), base, bitSize)
}

func (p *dataParser) ParseInt(delim byte, base int, bitSize int) (i64 int64) {
	return parseInt(p.ReadString(delim), base, bitSize)
}
