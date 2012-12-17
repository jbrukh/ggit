//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package api

import (
	"fmt"
	"github.com/jbrukh/ggit/api/objects"
	"strings"
)

// ================================================================= //
// REF SORTING
// ================================================================= //

// sort interface for sorting refs
type refByName []objects.Ref

func (s refByName) Len() int           { return len(s) }
func (s refByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s refByName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }

// ================================================================= //
// REF FILTERING
// ================================================================= //

// FilterRefs applies a ref filter to a list of refs.
func FilterRefs(refs []objects.Ref, f Filter) []objects.Ref {
	var r []objects.Ref
	for _, v := range refs {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

func FilterRefPattern(pattern string) Filter {
	return func(ref interface{}) bool {
		return matchRefs(ref.(objects.Ref).Name(), pattern)
	}
}

// FilterRefPrefix returns a filter that matches
// by the given prefix.
func FilterRefPrefix(prefix string) Filter {
	return func(ref interface{}) bool {
		return strings.HasPrefix(ref.(objects.Ref).Name(), prefix)
	}
}

// matchRefs performs the matching of a partial ref with a full (or longer)
// ref. Matching occurs from the end and matches on completed parts of the
// name. So for instance, refs/heads/master and master would match, but "ter"
// would not match the former.
func matchRefs(full, partial string) bool {
	const SL = "/"
	if full == "" || partial == "" {
		return false
	}

	f, p := strings.Split(full, SL), strings.Split(partial, SL)
	i, j := len(f), len(p)
	if i == 0 || j == 0 || i < j { // partial must be shorter
		return false
	}
	for j > 0 {
		i--
		j--
		if f[i] != p[j] {
			return false
		}
	}
	return true
}

// ================================================================= //
// ERRORS - used by DiskRepository
// ================================================================= //

type (
	noSuchRef    error
	ambiguousRef error
)

func noSuchRefErrf(ref string) noSuchRef {
	return noSuchRef(fmt.Errorf("no such ref: %s", ref))
}

func IsNoSuchRef(e error) bool {
	switch t := e.(type) {
	case noSuchRef:
		e = t // must use t
		return true
	}
	return false
}

// ================================================================= //
// OPERATIONS
// ================================================================= //

var refSearchPath = []string{
	"%s",
	"refs/%s",
	"refs/tags/%s",
	"refs/heads/%s",
	"refs/remotes/%s",
	"refs/remotes/%s/HEAD",
}

// RefFromSpec delivers a ref from the repository that is in
// the form of a Ref object. The ref may be symbolic or peeled.
// 
// This method resolves "shorthand" refs (e.g. when one writes 
// "master" for "refs/head/master". It follows the rules specified 
// here:
//
//    http://www.kernel.org/pub/software/scm/git/docs/gitrevisions.html
//
// under "Specifying Revisions/<refname>". If the ref does not exist
// then you can check the returned error with api.IsNoSuchRef.
// TODO: what about ambiguous refs?
func RefFromSpec(repo Repository, spec string) (ref objects.Ref, err error) {
	for _, prefix := range refSearchPath {
		refPath := fmt.Sprintf(prefix, spec)
		if ref, err = repo.Ref(refPath); err == nil {
			return ref, nil
		} else if !IsNoSuchRef(err) {
			return nil, err // something went wrong
		}
	}
	return nil, err // no such ref
}

// PeeledRefFromSpec takes the same arguments as RefFromSpec, but
// peels the ref before sending it back.
func PeeledRefFromSpec(repo Repository, spec string) (ref objects.Ref, err error) {
	if ref, err = RefFromSpec(repo, spec); err != nil {
		return nil, err
	}
	return PeelRef(repo, ref)
}

// PeelRef resolves the final target oid of the ref and returns
// a peeled ref for this target. It examines the target of the
// given ref and if the target is symbolic, it is followed and
// resolved. This process repeats as many times as necessary to
// obtain a peeled ref.
func PeelRef(repo Repository, r objects.Ref) (objects.Ref, error) {
	var (
		err      error
		symbolic bool
		target   interface{}
	)

	// TODO: make a limit
	for {
		symbolic, target = r.Target()
		if symbolic {
			r, err = repo.Ref(target.(string))
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return r, nil
}

func expandHeadRef(short string) string {
	return "refs/heads/" + short
}

func expandTagRef(short string) string {
	return "refs/tags/" + short
}
