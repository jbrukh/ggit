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
	e := safeParse(func() {
		if p.rev == "" {
			panicErr("revision spec is empty")
		}

		if p.PeekByte() == ':' {
			panicErr(": syntaxes not supported") // TODO
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
			panicErr("revision is empty")
		}

		err := p.findObject(rev)
		if err != nil {
			panicErr(err.Error())
		}

		for !p.EOF() {
			var parent *Commit
			var err error
			b := p.ReadByte()
			if b == '^' {
				if !p.EOF() && p.PeekByte() == '{' {
					p.ConsumeByte('{')
					otype := ObjectType(p.ConsumeStrings(objectTypes))
					applyDereference(p, otype)
					p.ConsumeByte('}')
				} else {
					parent, err = applyParentFunc(p, CommitNthParent)
				}
			} else if b == '~' {
				parent, err = applyParentFunc(p, CommitNthAncestor)
			} else {
				panicErrf("unexpected modifier: '%s'", string(b))
			}
			if err != nil {
				panicErr(err.Error())
			}
			p.o = parent
		}
	})
	return e
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
