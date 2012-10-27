//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import (
	"errors"
	//"fmt"
	//"os"
	"github.com/jbrukh/ggit/util"
	"regexp"
	"strconv"
)

// ================================================================= //
// CONSTANTS
// ================================================================= //

// regular expression for text-based hexadecimal strings
// that signify oid's or short oids
var hexRegex *regexp.Regexp

// parentFunc specifies a strategy for selecting a parent
// or an ancestor of a commit
type parentFunc func(Repository, *Commit, int) (*Commit, error)

func init() {
	hexRegex, _ = regexp.Compile("[0-9a-fA-F]{4,40}")
}

// ================================================================= //
// REV PARSER
// ================================================================= //

// revParser is a parser for revision specs.
type revParser struct {
	repo Repository

	inx int
	rev string

	o Object
}

// Object returns the object that the rev spec
// refers to after (and during) parsing.
func (p *revParser) Object() Object {
	return p.o
}

// more returns true if and only if the index
// of the parser has more characters to read.
func (p *revParser) more() bool {
	return p.inx < len(p.rev)
}

// next increments the index of the parser.
func (p *revParser) next() {
	p.inx++
}

// curr returns the current character in the
// revision that the index is pointing to.
func (p *revParser) curr() byte {
	return p.rev[p.inx]
}

// symbol returns the revision string until,
// but not including, the index.
func (p *revParser) symbol() string {
	return p.rev[:p.inx]
}

// peek peeks ahead a number of characters
// from the index.
func (p *revParser) peek(n int) string {
	return p.rev[p.inx : p.inx+n]
}

// number parses and returns the number that
// is represented by the string immediately
// following the index. If the string in question
// is empty, or we are at the end of the revision
// sting, then 1 is returned by default.
func (p *revParser) number() (n int) {
	p.next()
	start := p.inx
	if !p.more() {
		return 1 // 1 by default
	}
	for p.more() && util.IsDigit(p.curr()) {
		p.next()
	}
	strNum := p.rev[start:p.inx]
	if strNum == "" {
		return 1
	}

	// no error is possible, since we have
	// verified digits ahead of time
	// TODO: what about digit strings that are too long?
	n, _ = strconv.Atoi(strNum)
	return n
}

func (p *revParser) revParse() error {
	if p.rev == "" {
		return errors.New("revision spec is empty")
	}

	if p.peek(1) == ":" {
		return errors.New(": syntaxes not supported") // TODO
	}

	// read until modifier or end
	for p.more() {
		if !isModifier(p.curr()) {
			p.next()
		} else {
			break
		}
	}

	rev := p.symbol()
	if rev == "" {
		return errors.New("revision is empty")
	}
	err := p.findCommit(rev)
	if err != nil {
		return err
	}

	for p.more() {
		var parent *Commit
		var err error
		if p.curr() == '^' {
			parent, err = applyParentFunc(p, CommitNthParent)
		} else if p.curr() == '~' {
			parent, err = applyParentFunc(p, CommitNthAncestor)
		} else {
			panic("unknown modifier, shouldn't get here")
		}

		if err != nil {
			return err
		}
		p.o = parent
	}

	return nil
}

func applyParentFunc(p *revParser, f parentFunc) (*Commit, error) {
	n := p.number()

	c, err := CommitFromObject(p.repo, p.o)
	if err != nil {
		return nil, err
	}
	return f(p.repo, c, n)
}

func (p *revParser) findCommit(simpleRev string) (err error) {
	// oid or short oid
	var o Object
	switch {
	case hexRegex.MatchString(simpleRev):
		o, err = ObjectFromShortOid(p.repo, simpleRev)
		if err != nil {
			return err
		}
	default:
		ref, err := OidRefFromShortRef(p.repo, simpleRev)
		if err != nil {
			return err
		}
		o, err = ObjectFromOid(p.repo, ref.ObjectId())
		if err != nil {
			return err
		}
	}
	p.o = o
	return nil
}

func (r *revParser) parseNumber() (int, error) {
	c := r.rev[r.inx]
	if c != '^' && c != '~' {
		return 0, errors.New("not expecting a number")
	}
	i := r.inx + 1
	for i < len(r.rev) {
		if !util.IsDigit(r.rev[i]) {
			break
		}
		i++
	}
	n := r.rev[r.inx+1 : i]
	if n == "" {
		return 1, nil
	}
	num, err := strconv.Atoi(n)
	if err != nil {
		return 0, err
	}
	r.inx = i
	return num, nil
}

// ================================================================= //
// UTILITY METHODS
// ================================================================= //

// isModifier returns true if and only if the parameter
// is a supported modifier that may appear in rev parsing.
// The modifier usually comes after the rev spec, signifying
// a path to take from the commit object referred to by
// the spec.
func isModifier(c byte) bool {
	switch c {
	case '^', '~', '@':
		return true
	}
	return false
}
