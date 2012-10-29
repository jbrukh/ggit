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
	dataParser
	repo Repository

	inx int
	rev string

	o Object
}

func newRevParser(repo Repository, rev string) *revParser {
	return &revParser{
		dataParser: dataParser{
			buf: readerForString(rev),
		},
		repo: repo,
		rev:  rev,
	}
}

// Object returns the object that the rev spec
// refers to after (and during) parsing.
func (p *revParser) Object() Object {
	return p.o
}

// // more returns true if and only if the index
// // of the parser has more characters to read.
// func (p *revParser) more() bool {
// 	return p.inx < len(p.rev)
// }

// // next increments the index of the parser.
// func (p *revParser) next() {
// 	p.inx++
// }

// // curr returns the current character in the
// // revision that the index is pointing to.
// func (p *revParser) curr() byte {
// 	return p.rev[p.inx]
// }

// // symbol returns the revision string until,
// // but not including, the index.
// func (p *revParser) symbol() string {
// 	return p.rev[:p.inx]
// }

// // peek peeks ahead a number of characters
// // from the index.
// func (p *revParser) peek(n int) string {
// 	return p.rev[p.inx : p.inx+n]
// }

// number parses and returns the number that
// is represented by the string immediately
// following the index. If the string in question
// is empty, or we are at the end of the revision
// sting, then 1 is returned by default.
func (p *revParser) number() (n int) {
	start := p.Count()
	if p.EOF() {
		return 1 // 1 by default
	}
	for !p.EOF() && util.IsDigit(p.PeekByte()) {
		p.ReadByte()
	}
	end := p.Count()
	strNum := p.rev[start:end]
	if strNum == "" {
		return 1
	}

	// no error is possible, since we have
	// verified digits ahead of time
	// TODO: what about digit strings that are too long?
	n, _ = strconv.Atoi(strNum)
	return n
}

func (p *revParser) Parse() error {
	if p.rev == "" {
		return errors.New("revision spec is empty")
	}

	if p.PeekByte() == ':' {
		return errors.New(": syntaxes not supported") // TODO
	}

	start := p.Count()
	// read until modifier or end
	for !p.EOF() {
		if !isModifier(p.PeekByte()) {
			p.ReadByte()
		} else {
			break
		}
	}
	end := p.Count()

	rev := p.rev[start:end]
	if rev == "" {
		return errors.New("revision is empty")
	}

	err := p.findObject(rev)
	if err != nil {
		return err
	}

	for !p.EOF() {
		var parent *Commit
		var err error
		if p.PeekByte() == '^' {
			p.ConsumeByte('^')
			parent, err = applyParentFunc(p, CommitNthParent)
		} else if p.PeekByte() == '~' {
			p.ConsumeByte('~')
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

func (p *revParser) parseObjectType() (objectType ObjectType) {
	otype := p.ConsumeStrings(objectTypes)
	return ObjectType(otype)
}

func applyParentFunc(p *revParser, f parentFunc) (*Commit, error) {
	n := p.number()

	c, err := CommitFromObject(p.repo, p.o)
	if err != nil {
		return nil, err
	}
	return f(p.repo, c, n)
}

func applyDereference(p *revParser, otype ObjectType) error {
	switch otype {
	case ObjectCommit:
		c, err := CommitFromObject(p.repo, p.o)
		if err != nil {
			return err
		}
		p.o = c
	case ObjectTree:
		c, err := CommitFromObject(p.repo, p.o)
		if err != nil {
			return err
		}
		p.o, err = p.repo.ObjectFromOid(c.tree)
		if err != nil {
			return err
		}
	case ObjectTag:
		if p.o.Header().Type() != ObjectTag {
			return errors.New("cannot dereference non-tag to tag")
		}
	case ObjectBlob:
		if p.o.Header().Type() != ObjectBlob {
			return errors.New("cannot dereference non-blob to blob")
		}
	}
	return nil
}

func (p *revParser) findObject(simpleRev string) (err error) {
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
