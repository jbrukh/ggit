package api

// filter represents a generic filter that desides
// whether or not the given parameter should be
// filtered out. Each individual filter needs to
// assert type of the parameter and create a
// requisite typed Apply method.
type filter func(interface{}) bool

// apply applies a filter to a slice of interface{}.
// One should define a custom apply method for each
// type they wish to filter.
func apply(s []interface{}, f filter) []interface{} {
	r := make([]interface{}, 0)
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

// filterOr combines a number of filters
// together using a logical OR.
func filterOr(fs ...filter) filter {
	return func(i interface{}) bool {
		for _, f := range fs {
			if f(i) {
				return true
			}
		}
		return false
	}
}

// filterAnd combines a number of filters
// together using a logical AND.
func filterAnd(fs ...filter) filter {
	return func(i interface{}) bool {
		for _, f := range fs {
			if !f(i) {
				return false
			}
		}
		return true
	}
}

func filterNot(f filter) filter {
	return func(i interface{}) bool {
		return !f(i)
	}
}
