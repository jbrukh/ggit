package ggit

import (
    "testing"
)

func assert(t *testing.T, b bool, items ...interface{}) {
    if !b {
        t.Error("error: ", items)
    }
}

func assertPanic(t *testing.T, f func()) {
    defer func() {
        if r := recover(); r != nil {
            return
        }
        // should never get here
    }()
    f()
    t.Error("was expecting a panic")
}
