package ggit

import (
    "bufio"
    "bytes"
    "testing"
)

func Test_parseOidLine(t *testing.T) {
    const T1 = "commit " + CRAZY + "\n"
    buf := bytes.NewBuffer([]byte(T1))
    m, oid, e := parseOidLine(bufio.NewReader(buf))
    assert(t, e == nil, "error: ", e)
    assert(t, m == "commit", "wrong marker")
    assert(t, oid.String() == CRAZY)
}
