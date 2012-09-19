package ggit

import (
    "bufio"
    "bytes"
    "fmt"
    "io"
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
    return b
}

func (p *dataParser) peek(n int) []byte {
    b := make([]byte, n)
    if pk, e := p.buf.Peek(n); e != nil || len(pk) != n {
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
    return b.Bytes()
}

// ================================================================= //
// SPECIALIZED PARSING FUNCTIONS
// ================================================================= //

// ParseObjectId reads the next OID_HEXSZ bytes from the
// Reader and places the resulting object id in oid.
func (p *dataParser) ParseObjectId(oid **ObjectId) {
    hex := string(p.consume(OID_HEXSZ))
    id, e := NewObjectIdFromString(hex)
    if e != nil {
        panicErrf("expected: hex string of size %d", OID_HEXSZ)
    }
    *oid = id
}

// func (p *dataParser) NextObjectIdString() *ObjectId {
// 	hex := p.NextString(OID_HEXSZ)
// 	oid, e := NewObjectIdFromString(hex)
// 	if e != nil {
// 		panicErr(e.Error())
// 	}
// 	return oid
// }

// // NextObjectIdString reads the next OID_SZ bytes from the Reader
// // and interprets them as an ObjectId
// func (p *dataParser) NextObjectIdBytes() *ObjectId {
// 	b := p.ReadBytes(OID_SZ)
// 	oid, e := NewObjectIdFromBytes(b)
// 	if e != nil {
// 		panicErr(e.Error())
// 	}
// 	return oid
// }

// //
// func (p *dataParser) TokenStringInt(delim byte) (n int) {
// 	str := p.TokenString(delim)
// 	var e error
// 	if n, e = strconv.Atoi(str); e != nil {
// 		panicErrn("cannot convert integer: %s", str)
// 	}
// 	return n
// }
