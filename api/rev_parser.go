package api

import (
	"errors"
	"strconv"
)

type revParser struct {
	repo Repository
	inx  int
	spec string // the whole spec
	ref  string // the ref on the left
	o    Object
}

// func (r *revParser) revParse() (Object, error) {
// 	l := len(spec)
// 	if l < 1 {
// 		return nil, errors.New("spec is empty")
// 	}
// 	r.inx = 0
// 	var err error
// 	for r.inx < l {
// 		c := spec[r.inx]
// 		switch c {
// 		case '^':
// 		case '~':
// 			r.ref = r.spec[:r.inx]
// 			r.o, err = ObjectFromRef(r.ref)
// 			if err != nil {
// 				return nil, err
// 			}

// 			var n int
// 			n, err = r.parseNumber()
// 			if err != nil {
// 				return nil, err
// 			}

// 		default:
// 			// TODO check if ref is parsed?
// 			r.inx++
// 		}
// 	}
// }

func (r *revParser) parseNumber() (int, error) {
	c := r.spec[r.inx]
	if c != '^' && c != '~' {
		return 0, errors.New("not expecting a number")
	}

	i := r.inx + 1
	for i < len(r.spec) {
		println("i: ", i)
		if !isDigit(r.spec[i]) {
			break
		}
		i++
	}

	n := r.spec[r.inx+1 : i]
	println("text: ", n)
	if n == "" {
		return 1, nil
	}

	return strconv.Atoi(n)
}

func isDigit(c byte) bool {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}
