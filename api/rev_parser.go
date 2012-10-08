package api

import (
//"errors"
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
	return 0, nil // TODO
}

func isDigit(c byte) bool {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}
