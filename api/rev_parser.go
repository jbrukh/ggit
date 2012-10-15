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
	"fmt"
	"strconv"
)

func CommitFromRevision(repo Repository, spec string) (*Commit, error) {
	r := &revParser{
		repo: repo,
		spec: spec,
	}
	return r.revParse()
}

type revParser struct {
	repo Repository
	inx  int
	spec string // the whole spec
	ref  string // the ref on the left
	c    *Commit
}

func (r *revParser) revParse() (*Commit, error) {
	l := len(r.spec)
	if l < 1 {
		return nil, errors.New("spec is empty")
	}
	r.inx = 0
	//var err error

	// find the left hand revision
	for !isModifier(r.spec[r.inx]) {
		r.inx++
	}

	// write down the ref
	r.ref = r.spec[:r.inx]

	// for r.inx < l {
	// 	c := r.spec[r.inx]
	// 	switch c {
	// 	case '^':
	// 	case '~':
	// 		r.ref = r.spec[:r.inx]
	// 		r.c, err = CommitFromRef(r.repo, r.ref)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		var n int
	// 		n, err = r.parseNumber()
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		return CommitNthAncestor(r.repo, r.c, n)

	// 	default:
	// 		// TODO check if ref is parsed?
	// 		r.inx++
	// 	}
	// }

	if r.ref != "" {
		oid, _ := OidFromString(r.ref)
		var err error
		r.c, err = CommitFromOid(r.repo, oid)
		if err != nil {
			return nil, err
		}
	}

	if r.c != nil {
		return r.c, nil
	}
	return nil, fmt.Errorf("Unknown revision: %s", r.spec)
}

func (r *revParser) parseNumber() (int, error) {
	c := r.spec[r.inx]
	if c != '^' && c != '~' {
		return 0, errors.New("not expecting a number")
	}
	i := r.inx + 1
	for i < len(r.spec) {
		if !isDigit(r.spec[i]) {
			break
		}
		i++
	}
	n := r.spec[r.inx+1 : i]
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
