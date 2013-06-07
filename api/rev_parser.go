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
	"github.com/jbrukh/ggit/api/objects"
	"github.com/jbrukh/ggit/api/token"
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
type parentFunc func(Repository, *objects.Commit, int) (*objects.Commit, error)

func init() {
	hexRegex, _ = regexp.Compile("[0-9a-fA-F]{4,40}")
}

// ================================================================= //
// REV PARSER
// ================================================================= //

// revParser is a parser for revision specs.
type revParser struct {
	util.DataParser
	repo Repository

	inx int
	rev string

	o objects.Object
}

func newRevParser(repo Repository, rev string) *revParser {
	return &revParser{
		*util.NewDataParser(util.ReaderForString(rev)),
		repo,
		0,
		rev,
		nil,
	}
}

// Object returns the object that the rev spec
// refers to after (and during) parsing.
func (p *revParser) Object() objects.Object {
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
	e := util.SafeParse(func() {
		if p.rev == "" {
			util.PanicErr("revision spec is empty")
		}

		if p.PeekByte() == ':' {
			util.PanicErr(": syntaxes not supported") // TODO
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
			util.PanicErr("revision is empty")
		}

		err := p.findObject(rev)
		if err != nil {
			util.PanicErr(err.Error())
		}

		for !p.EOF() {
			var err error
			b := p.ReadByte()
			if b == '^' {
				if !p.EOF() && p.PeekByte() == '{' {
					p.ConsumeByte('{')
					otype := objects.ObjectType(p.ConsumeStrings(token.ObjectTypes))
					err = applyDereference(p, otype)
					if err != nil {

					}
					p.ConsumeByte('}')
				} else {
					err = applyParentFunc(p, CommitNthParent)
				}
			} else if b == '~' {
				err = applyParentFunc(p, CommitNthAncestor)
			} else {
				util.PanicErrf("unexpected modifier: '%s'", string(b))
			}

			if err != nil {
				util.PanicErr(err.Error())
			}
		}
	})
	return e
}

func applyParentFunc(p *revParser, f parentFunc) (err error) {
	n := p.number()
	var c, parent *objects.Commit
	c, err = CommitFromObject(p.repo, p.o)
	if err != nil {
		return err
	}
	parent, err = f(p.repo, c, n)
	if err != nil {
		return err
	}
	p.o = parent
	return
}

func applyDereference(p *revParser, otype objects.ObjectType) error {
	switch otype {
	case objects.ObjectCommit:
		c, err := CommitFromObject(p.repo, p.o)
		if err != nil {
			return err
		}
		p.o = c
	case objects.ObjectTree:
		c, err := CommitFromObject(p.repo, p.o)
		if err != nil {
			return err
		}
		p.o, err = p.repo.ObjectFromOid(c.Tree())
		if err != nil {
			return err
		}
	case objects.ObjectTag:
		if p.o.Header().Type() != objects.ObjectTag {
			return errors.New("cannot dereference non-tag to tag")
		}
	case objects.ObjectBlob:
		if p.o.Header().Type() != objects.ObjectBlob {
			return errors.New("cannot dereference non-blob to blob")
		}
	}
	return nil
}

func (p *revParser) findObject(spec string) (err error) {
	// oid or short oid
	var o objects.Object
	switch {
	case hexRegex.MatchString(spec):
		o, err = p.repo.ObjectFromShortOid(spec)
		if err != nil {
			return err
		}
	default:
		ref, err := PeeledRefFromSpec(p.repo, spec)
		if err != nil {
			return err
		}
		o, err = p.repo.ObjectFromOid(ref.ObjectId())
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
