package api

import (
	"io"
)

// Formatter is the central hub for formatting ggit objects.
type Format struct {
	W io.Writer
}
