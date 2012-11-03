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
type parseErr struct {
	msg   string
	stack []byte
}

// ParseErr is an error
func (p *parseErr) Error() string {
	return p.msg
}

func (p *parseErr) Stack() string {
	return string(p.stack)
}

// parseErrf allows convenience formatting for ParseErrors
func parseErrf(format string, items ...interface{}) *parseErr {
	return &parseErr{
		msg:   fmt.Sprintf(format, items...),
		stack: debug.Stack(),
	}
}

// parseErrn concatenates a bunch of items together to 
// form the error string
func parseErrn(items ...string) *parseErr {
	return parseErrf("%s", strings.Join(items, ""))
}

func panicErr(msg string) {
	panic(parseErrn(msg))
}

func panicErrn(items ...string) {
	panic(parseErrn(items...))
}

func panicErrf(format string, items ...interface{}) {
	panic(parseErrf(format, items...))
}

// ================================================================= //
// DATA PARSER
// ================================================================= //

type dataParser struct {
	buf   *bufio.Reader
	count int64
}

// safeParse allows you to call a number of parsing functions on your
// parser at once, without having to handle errors explicitly. If an
// error occurs, the parser commands will panic with parseErr, which
// this method will recover and return
func safeParse(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(*parseErr); ok {
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

func (p *dataParser) consume(n int) []byte {
	b := make([]byte, n)
	if rd, e := p.buf.Read(b); e != nil {
		panicErrf("expected: %d byte(s), read %d, values %x", n, rd, b[0:rd])
	} else if rd != n {
		more := p.consume(n - rd)
		for i, v := range more {
			b[rd+i] = v
		}
	}
	p.count += int64(n)
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
	p.count += int64(len(b))
	return trimLastByte(b)
}

// ResetCount resets the read byte count
func (p *dataParser) ResetCount() {
	p.count = 0
}

// Count returns the number of bytes read since
// the parser was initialized or ResetCount() was
// called, whichever came last.
func (p *dataParser) Count() int64 {
	return p.count
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
		if b[inx] != v {
			panicErrf("expected bytes: 0x%x, found: 0x%x", b, d)
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
		panicErrf("expected: %d byte(s); got: %s", n, e.Error())
	}
	return string(pk)
}

func (p *dataParser) ReadByte() byte {
	b := p.PeekBytes(1)
	p.Consume(1)
	return b[0]
}

// ReadBytesUntil
func (p *dataParser) ReadNBytes(n int) []byte {
	return p.consume(n)
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
	p.count += int64(len(bts))
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
func (p *dataParser) ParseAtoi(delim byte) (n int64) {
	return p.ParseInt(delim, 10, 64)
}

// Returns the int64 represented by the next n bytes in network byte order (most significant first).
func (p *dataParser) ParseIntBigEndian(n int) (i64 int64) {
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

func (p *dataParser) ParseIntN(n int, base int, bitSize int) (i64 int64) {
	bytes := p.consume(n)
	value := string(bytes)
	return parseInt(value, base, bitSize)
}

func (p *dataParser) ParseInt(delim byte, base int, bitSize int) (i64 int64) {
	return parseInt(p.ReadString(delim), base, bitSize)
}
