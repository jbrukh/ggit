package api

import (
	"io"
)

// Formatter is the central hub for formatting ggit objects.
type Formatter struct {
	W io.Writer
}
