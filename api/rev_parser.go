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
	"regexp"
	"strconv"
)

var hexRegex *regexp.Regexp

func init() {
	hexRegex, _ = regexp.Compile("[0-9a-fA-F]{4,40}")
}

func ObjectFromRevision(repo Repository, rev string) (Object, error) {
	p := &revParser{
		repo: repo,
		rev:  rev,
	}
	e := p.revParse()
	if e != nil {
		return nil, e
	}
	return p.Object(), nil
}

type revParser struct {
	repo Repository

	inx int
	rev string

	o Object
}

func (p *revParser) Object() Object {
	return p.o
}

func (p *revParser) more() bool {
	return p.inx < len(p.rev)
}

func (p *revParser) next() {
	p.inx++
}

func (p *revParser) curr() byte {
	return p.rev[p.inx]
}

func (p *revParser) symbol() string {
	return p.rev[:p.inx]
}

func (p *revParser) peek(n int) string {
	return p.rev[p.inx : p.inx+n]
}

func (p *revParser) number() (n int, err error) {
	p.next()
	start := p.inx
	if !p.more() {
		return 1, nil // 1 by default
	}
	for p.more() && isDigit(p.curr()) {
		p.next()
	}
	strNum := p.rev[start:p.inx]
	if strNum == "" {
		return 1, nil
	}
	return strconv.Atoi(strNum)
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

type parentFunc func(Repository, *Commit, int) (*Commit, error)

func applyParentFunc(p *revParser, f parentFunc) (*Commit, error) {
	n, err := p.number()
	if err != nil {
		return nil, err
	}

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
		if !isDigit(r.rev[i]) {
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

func isDigit(c byte) bool {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}

func isModifier(c byte) bool {
	switch c {
	case '^', '~', '@':
		return true
	}
	return false
}
