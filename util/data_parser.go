//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package util

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
)

// ================================================================= //
// PARSE ERROR TYPE & PANICS
// ================================================================= //

// ParserDebug toggles parser debugging in safeParse().
var ParserDebug bool = true

// ================================================================= //
// PARSE ERROR TYPE & PANICS
// ================================================================= //

// ParseErr is a common error that occurs when ggit is 
// parsing binary objects
type ParseErr struct {
	msg   string
	stack []byte
}

// ParseErr is an error
func (p *ParseErr) Error() string {
	return p.msg
}

func (p *ParseErr) Stack() string {
	return string(p.stack)
}

// ParseErrf allows convenience formatting for ParseErrors
func ParseErrf(format string, items ...interface{}) *ParseErr {
	return &ParseErr{
		msg:   fmt.Sprintf(format, items...),
		stack: debug.Stack(),
	}
}

// ParseErrn concatenates a bunch of items together to 
// form the error string
func ParseErrn(items ...string) *ParseErr {
	return ParseErrf("%s", strings.Join(items, ""))
}

func PanicErr(msg string) {
	panic(ParseErrn(msg))
}

func PanicErrn(items ...string) {
	panic(ParseErrn(items...))
}

func PanicErrf(format string, items ...interface{}) {
	panic(ParseErrf(format, items...))
}

// ================================================================= //
// UTIL
// ================================================================= //

func ParserForBytes(b []byte) *DataParser {
	return &DataParser{
		buf: bufio.NewReader(bytes.NewBuffer(b)),
	}
}

func ParserForString(str string) *DataParser {
	return ParserForBytes([]byte(str))
}

// ================================================================= //
// DATA PARSER
// ================================================================= //

type DataParser struct {
	buf   *bufio.Reader
	count int64
}

func NewDataParser(rd *bufio.Reader) *DataParser {
	return &DataParser{
		buf: rd,
	}
}

// safeParse allows you to call a number of parsing functions on your
// parser at once, without having to handle errors explicitly. If an
// error occurs, the parser commands will panic with ParseErr, which
// this method will recover and return
func SafeParse(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(*ParseErr); ok {
				err = e
				if ParserDebug {
					fmt.Fprintln(os.Stderr, "------------ ParserDebug Output ------------")
					fmt.Fprint(os.Stderr, e.Stack())
				}
			}
		}
	}()
	f()
	return
}

// ================================================================= //
// DATA PARSING API
// ================================================================= //

func (p *DataParser) consume(n int) []byte {
	b := make([]byte, n)
	if rd, e := p.buf.Read(b); e != nil {
		PanicErrf("expected: %d byte(s), read %d, values %x", n, rd, b[0:rd])
	} else if rd != n {
		more := p.consume(n - rd)
		for i, v := range more {
			b[rd+i] = v
		}
	}
	p.count += int64(n)
	return b
}

func (p *DataParser) peek(n int) (b []byte) {
	var e error
	if b, e = p.buf.Peek(n); e != nil || len(b) != n {
		PanicErrf("expected: %d byte(s)", n)
	}
	return b
}

func (p *DataParser) consumeUntil(delim byte) []byte {
	b, e := p.buf.ReadBytes(delim)
	if e != nil {
		PanicErrf("expected delimiter: %v", delim)
	}
	p.count += int64(len(b))
	return TrimLastByte(b)
}

// ResetCount resets the read byte count
func (p *DataParser) ResetCount() {
	p.count = 0
}

// Count returns the number of bytes read since
// the parser was initialized or ResetCount() was
// called, whichever came last.
func (p *DataParser) Count() int64 {
	return p.count
}

func (p *DataParser) EOF() bool {
	if _, e := p.buf.Peek(1); e != nil {
		if e == io.EOF {
			return true
		} else {
			PanicErr("reading error")
		}
	}
	return false
}

// Consume will consume n bytes without regard for what the underlying
// data might be. If it is unable to consume, then a panic is raised
// with ParseErr.
func (p *DataParser) Consume(n int) []byte {
	return p.consume(n)
}

// ConsumeByte will consume a single byte and compare it to b. If it
// does not match, or cannot be read, then a panic is raised with ParseErr.
func (p *DataParser) ConsumeByte(b byte) {
	if p.consume(1)[0] != b {
		PanicErrf("expected byte: %v", b)
	}
}

// PeekBute will return the next byte without advancing the reader. If
// it cannot be read, then a panic is raised with ParseErr.
func (p *DataParser) PeekByte() byte {
	return p.peek(1)[0]
}

// PeekBytes
func (p *DataParser) PeekBytes(n int) []byte {
	return p.peek(n)
}

// ConsumeBytes will consume len(b) bytes and compare them to b. If they
// do not match, or cannot be read, then a panic is raised with ParseErr.
func (p *DataParser) ConsumeBytes(b []byte) {
	d := p.consume(len(b))
	for inx, v := range d {
		if b[inx] != v {
			PanicErrf("expected bytes: 0x%x, found: 0x%x", b, d)
		}
	}
}

// ConsumeString will consume len(b) bytes and compare them to the string s. If they
// do not match, or cannot be read, then a panic is raised with ParseErr.
func (p *DataParser) ConsumeString(s string) {
	b := p.consume(len(s))
	if string(b) != s {
		PanicErrf("expected string: %s", s)
	}
}

// ConsumeStrings will check if the Reader contains any one of the provided strings
// and if so, consumes it and returns it. If a string match cannot be found, or 
// cannot be read, than a panic is raised with ParseErr. (The first string matched
// is the one returned, such that if strings are substrings of one another, only
// the first match matters.)
func (p *DataParser) ConsumeStrings(s []string) string {
	for _, str := range s {
		l := len(str)
		pk := p.PeekString(l)
		if l != 0 && pk == str {
			return string(p.consume(l))
		}
	}
	PanicErrf("expected one of: %v", s)
	return ""
}

// PeekString returns true if and only if the next bytes
// in the buffer match the given input string (the string
// in the buffer is NOT consumed)
func (p *DataParser) PeekString(n int) string {
	pk, e := p.buf.Peek(n)
	if e != nil {
		PanicErrf("expected: %d byte(s); got: %s", n, e.Error())
	}
	return string(pk)
}

func (p *DataParser) ReadByte() byte {
	b := p.PeekBytes(1)
	p.Consume(1)
	return b[0]
}

// ReadBytesUntil
func (p *DataParser) ReadNBytes(n int) []byte {
	return p.consume(n)
}

// ReadBytesUntil
func (p *DataParser) ReadBytes(delim byte) []byte {
	return p.consumeUntil(delim)
}

// ReadStringUtil
func (p *DataParser) ReadString(delim byte) string {
	return string(p.consumeUntil(delim))
}

// String returns the entirety of the remaining data
// in the buffer, up to the EOF, as a string
func (p *DataParser) String() string {
	return string(p.Bytes())
}

// Bytes returns the entirety of the remaining data
// in the buffer, up to the EOF, as bytes
func (p *DataParser) Bytes() []byte {
	b := new(bytes.Buffer)
	_, e := io.Copy(b, p.buf)
	if e != nil {
		PanicErr(e.Error())
	}
	bts := b.Bytes()
	p.count += int64(len(bts))
	return bts
}

// ================================================================= //
// SPECIALIZED PARSING FUNCTIONS
// ================================================================= //

func parseInt(str string, base int, bitSize int) (i64 int64) {
	var e error
	if i64, e = strconv.ParseInt(str, base, bitSize); e != nil {
		PanicErrf("cannot convert integer (base $d): %s", base, str)
	}
	return i64
}

// TODO: this should be smarter than delimiter
func (p *DataParser) ParseAtoi(delim byte) (n int64) {
	return p.ParseInt(delim, 10, 64)
}

// Returns the int64 represented by the next n bytes in network byte order (most significant first).
func (p *DataParser) ParseIntBigEndian(n int) (i64 int64) {
	switch {
	case n == 4:
		return int64(binary.BigEndian.Uint32(p.consume(4)))
	case n == 8:
		return int64(binary.BigEndian.Uint64(p.consume(8)))
	}
	bytes := p.consume(n)
	value := fmt.Sprintf("%x", bytes)
	return parseInt(value, 16, 64)
}

func (p *DataParser) ParseIntN(n int, base int, bitSize int) (i64 int64) {
	bytes := p.consume(n)
	value := string(bytes)
	return parseInt(value, base, bitSize)
}

func (p *DataParser) ParseInt(delim byte, base int, bitSize int) (i64 int64) {
	return parseInt(p.ReadString(delim), base, bitSize)
}
