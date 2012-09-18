package ggit

import (
    "bufio"
    "fmt"
    "strings"
)

// ================================================================= //
// PARSE ERROR TYPE
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

// parseErrn concatenates a bunch of items together
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
// DATA PARSER
// ================================================================= //

type dataParser struct {
    buf *bufio.Reader
}

// ================================================================= //
// DATA PARSING PANICS
// 
// The idea is that the parser will panic on any error and the user
// will catch the panic. The type of object thrown during a panic
// is parseErr.
// ================================================================= //

func (p *dataParser) TokenString(delim byte) string {
    str, e := p.buf.ReadString(delim)
    if e != nil {
        panicErrf(e.Error())
    }
    return str
}

func (p *dataParser) TokenBytes(delim byte) []byte {
    b, e := p.buf.ReadBytes(delim)
    if e != nil {
        panicErr(e.Error())
    }
    return b
}

func (p *dataParser) NextString(n int) string {
    return string(p.NextBytes(n))
}

func (p *dataParser) NextBytes(n int) []byte {
    b := make([]byte, n)
    if rd, e := p.buf.Read(b); e != nil || rd != n {
        panicErrf("couldn't read %d bytes", n)
    }
    return b
}

// nextObjectIdString reads the next OID_HEXSZ bytes from the Reader
// and interprets them as an ObjectId
func (p *dataParser) NextObjectIdString() *ObjectId {
    hex := p.NextString(OID_HEXSZ)
    oid, e := NewObjectIdFromString(hex)
    if e != nil {
        panicErr(e.Error())
    }
    return oid
}

// nextObjectIdString reads the next OID_SZ bytes from the Reader
// and interprets them as an ObjectId
func (p *dataParser) NextObjectIdBytes() *ObjectId {
    b := p.NextBytes(OID_SZ)
    oid, e := NewObjectIdFromBytes(b)
    if e != nil {
        panicErr(e.Error())
    }
    return oid
}

func (p *dataParser) NextIntString(delim byte) int {
    return 0 // TODO
}

func (p *dataParser) NextInt32String(delim byte) int32 {
    return 0 // TODO
}

// FlushString returns the entirety of the remaining data
// in the buffer, up to the EOF, as a string
func (p *dataParser) FlushString() string {
    return "" // TODO
}

func (p *dataParser) FlushBytes() []byte {
    return nil // TODO
}

func (p *dataParser) VerifyString(str string) bool {
    return p.NextString(len(str)) == str
}
