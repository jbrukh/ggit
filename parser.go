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
    buf *bufio.Reader
}

// ================================================================= //
// DATA PARSING API
// ================================================================= //

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

// TokenBytes returns the next token of bytes delimited
// by the given byte, not including the delimiter
func (p *dataParser) TokenBytes(delim byte) []byte {
    b, e := p.buf.ReadBytes(delim)
    if e != nil {
        panicErr(e.Error())
    }
    return trimLast(b)
}

// TokenStringreturns the next token of bytes delimited
// by the given byte, not including the delimiter, as 
// a string
func (p *dataParser) TokenString(delim byte) string {
    return string(p.TokenBytes(delim))
}

// NextBytes returns the next n bytes of the Reader,
// or bust
func (p *dataParser) NextBytes(n int) []byte {
    b := make([]byte, n)
    if rd, e := p.buf.Read(b); e != nil || rd != n {
        panicErrf("couldn't read %d bytes", n)
    }
    return b
}

// NextBytes returns the next n bytes of the Reader,
// or bust, as a string
func (p *dataParser) NextString(n int) string {
    return string(p.NextBytes(n))
}

// NextObjectIdString reads the next OID_HEXSZ bytes from the Reader
// and interprets them as an ObjectId
func (p *dataParser) NextObjectIdString() *ObjectId {
    hex := p.NextString(OID_HEXSZ)
    oid, e := NewObjectIdFromString(hex)
    if e != nil {
        panicErr(e.Error())
    }
    return oid
}

// ParseObjectIdString reads the next OID_HEXSZ bytes from the
// Reader and places the resulting object id in oid
func (p *dataParser) ParseObjectIdString(oid **ObjectId) {
    *oid = p.NextObjectIdString()
}

// NextObjectIdString reads the next OID_SZ bytes from the Reader
// and interprets them as an ObjectId
func (p *dataParser) NextObjectIdBytes() *ObjectId {
    b := p.NextBytes(OID_SZ)
    oid, e := NewObjectIdFromBytes(b)
    if e != nil {
        panicErr(e.Error())
    }
    return oid
}

//
func (p *dataParser) TokenStringInt(delim byte) (n int) {
    str := p.TokenString(delim)
    var e error
    if n, e = strconv.Atoi(str); e != nil {
        panicErrn("cannot convert integer: %s", str)
    }
    return n
}

// FlushString returns the entirety of the remaining data
// in the buffer, up to the EOF, as a string
func (p *dataParser) FlushString() string {
    return string(p.FlushBytes())
}

// FlushBytes returns the entirety of the remaining data
// in the buffer, up to the EOF, as bytes
func (p *dataParser) FlushBytes() []byte {
    b := new(bytes.Buffer)
    _, e := io.Copy(b, p.buf)
    if e != nil {
        panicErr(e.Error())
    }
    return b.Bytes()
}

// VerifyString panics if and only if the next bytes
// in the buffer do not match the given input string (the string
// in the buffer is consumed)
func (p *dataParser) VerifyString(str string) {
    // TODO	: can implement this more efficiently with ReadByte()
    if p.NextString(len(str)) != str {
        panicErrn("data did not match: ", str)
    }
}

// PeekString returns true if and only if the next bytes
// in the buffer match the given input string (the string
// in the buffer is NOT consumed)
func (p *dataParser) PeekString(str string) bool {
    peek, e := p.buf.Peek(len(str))
    if e != nil {
        panicErr(e.Error())
    }
    return string(peek) == str
}

// VerifyRune consumes the next byte and panics with parseErr
// if it does not match the reference byte b.
func (p *dataParser) VerifyByte(b byte) {
    data, e := p.buf.ReadByte()
    if e != nil {
        panicErr(e.Error())
    }
    if b != data {
        panicErrf("data did not match: %s", b)
    }
}

func (p *dataParser) PeekByte(r byte) bool {
    return p.PeekString(string(r))
}
