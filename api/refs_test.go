package api

import (
	"sort"
	"testing"
)

func Test_Sorting(t *testing.T) {
	refs := make([]Ref, 3)
	refs[0] = &NamedRef{nil, "zoo"}
	refs[1] = &NamedRef{nil, "yogurt"}
	refs[2] = &NamedRef{nil, "xavier"}
	sort.Sort(refByName(refs))
	assert(t, refs[0].Name() == "xavier")
	assert(t, refs[1].Name() == "yogurt")
	assert(t, refs[2].Name() == "zoo")
}

func Test_matchRefs(t *testing.T) {
	assert(t, matchRefs("refs/heads/master", "master"))
	assert(t, matchRefs("refs/heads/master", "heads/master"))
	assert(t, matchRefs("refs/heads/master", "refs/heads/master"))
	assert(t, matchRefs("m/a/s/t/e/r", "a/s/t/e/r"))
	assert(t, !matchRefs("refs/heads/master", "ster"))
	assert(t, !matchRefs("refs/heads/master", "ds/master"))
	assert(t, !matchRefs("refs/heads/master", ""))
	assert(t, !matchRefs("refs/heads/master", "/refs/heads/master"))
	assert(t, !matchRefs("master", "refs/heads/master"))
	assert(t, !matchRefs("", ""))
}

func Test_RegexpCaret(t *testing.T) {
	assert(t, regexpCaret.MatchString("master^"))
	assert(t, regexpCaret.MatchString("master^^"))
	assert(t, regexpCaret.MatchString("master^^^"))
	assert(t, regexpCaret.MatchString("jake-dev^"))
	assert(t, regexpCaret.MatchString("ab334def^"))
	assert(t, !regexpCaret.MatchString("master"))
	assert(t, !regexpCaret.MatchString("^"))
	assert(t, !regexpCaret.MatchString("^^"))
	assert(t, !regexpCaret.MatchString("^^^"))
}

func Test_RegexpTilde(t *testing.T) {
	assert(t, regexpTilde.MatchString("master~1"))
	assert(t, regexpTilde.MatchString("master~19"))
	assert(t, regexpTilde.MatchString("stupid_shit~19"))
	assert(t, !regexpTilde.MatchString("master~"))
	assert(t, !regexpTilde.MatchString("master~0"))
	assert(t, !regexpTilde.MatchString("~0"))
}

func Test_RegexpHex(t *testing.T) {
	assert(t, regexpHex.MatchString("abcdef1234567890000000000000000000000000"))
	assert(t, regexpHex.MatchString("8cca23d1b1da6d712f5171aef414f1298906be85"))
	assert(t, regexpHex.MatchString("87e5cfb478d6d33e33f931ef32e896bf1b7f590b"))
	assert(t, regexpHex.MatchString("d4e482f121739ba54c02fb54dfa9e5ee91c6afd3"))
	assert(t, !regexpHex.MatchString("4e482f121739ba54c02fb54dfa9e5ee91c6afd3"))   // shorter
	assert(t, !regexpHex.MatchString("4e482f121739ba54c02fb54dfa9e5ee91c6afd3e2")) // longer
	assert(t, !regexpHex.MatchString("d4e482f121739ba54c02fx54dfa9e5ee91c6afd3"))  // has wayward character
	assert(t, !regexpHex.MatchString("d4e482f121739ba54c02fb54dfa9e5ee91c6afg3"))  // has wayward character
	assert(t, !regexpHex.MatchString(""))                                          // wtf?
}

func Test_RegexpShortHex(t *testing.T) {
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f5171aef414f1298906be8")) // 39
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f5171aef414f1298906be"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f5171aef414f1298906b"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f5171aef414f1298906"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f5171aef414f129890"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f5171aef414f12989"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f5171aef414f1298"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f5171aef414f129"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f5171aef414f12"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f5171aef414f1"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f5171aef414f"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f5171aef414"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f5171aef41"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f5171aef4"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f5171aef"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f5171ae"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f5171a"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f5171"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f517"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f51"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f5"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712f"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d712"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d71"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d7"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6d"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da6"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1da"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1d"))
	assert(t, regexpShortHex.MatchString("8cca23d1b1"))
	assert(t, regexpShortHex.MatchString("8cca23d1b"))
	assert(t, regexpShortHex.MatchString("8cca23d1"))
	assert(t, regexpShortHex.MatchString("8cca23d"))
	assert(t, regexpShortHex.MatchString("8cca23"))
	assert(t, regexpShortHex.MatchString("8cca2"))
	assert(t, regexpShortHex.MatchString("8cca"))
	assert(t, regexpShortHex.MatchString("8cc"))
	assert(t, !regexpShortHex.MatchString("8c"))
	assert(t, !regexpShortHex.MatchString("8"))
	assert(t, !regexpShortHex.MatchString("8ccg23d1b1da6d712f5"))
	assert(t, !regexpShortHex.MatchString("8cca23dzb1da6d712f"))

}
