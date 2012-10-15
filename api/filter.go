//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

// filter represents a generic filter that desides
// whether or not the given parameter should be
// filtered out. Each individual filter needs to
// assert type of the parameter and create a
// requisite typed Apply method.
type Filter func(interface{}) bool

// apply applies a filter to a slice of interface{}.
// One should define a custom apply method for each
// type they wish to filter.
func Apply(s []interface{}, f Filter) []interface{} {
	r := make([]interface{}, 0)
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

// FilterOr combines a number of filters
// together using a logical OR.
func FilterOr(fs ...Filter) Filter {
	return func(i interface{}) bool {
		for _, f := range fs {
			if f(i) {
				return true
			}
		}
		return false
	}
}

// FilterAnd combines a number of filters
// together using a logical AND.
func FilterAnd(fs ...Filter) Filter {
	return func(i interface{}) bool {
		for _, f := range fs {
			if !f(i) {
				return false
			}
		}
		return true
	}
}

func FilterNot(f Filter) Filter {
	return func(i interface{}) bool {
		return !f(i)
	}
}
