package api

import (
//"errors"
)

type revParser struct {
	repo Repository
	inx  int
	spec string
	ref  string
	o    Object
}

// func (r *revParser) revParse() (Object, error) {
// 	l := len(spec)
// 	if l < 1 {
// 		return nil, errors.New("spec is empty")
// 	}
// 	r.inx = 0

// 	for r.inx < l {
// 		c := spec[r.inx]
// 		switch c {
// 		case '^':

// 		case '~':
// 		default:
// 			// TODO check if ref is parsed?
// 			r.inx++
// 		}
// 	}

// }
