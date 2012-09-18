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

// ================================================================= //
// DATA PARSER
// ================================================================= //

type dataParser struct {
    buf *bufio.Reader
}

func (p *dataParser) panicErrn(items ...string) {
    panic(parseErrn(items...))
}

func (p *dataParser) panicErrf(format string, items ...string) {
    panic(parseErrf(format, items))
}

func (p *dataParser) tokenString(delim byte) string {
    return "" // TODO
}

func (p *dataParser) tokenBytes(delim byte) []byte {
    return []byte{} // TODO
}
