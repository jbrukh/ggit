package builtin

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// the set of supported builtin commands for ggit
var builtins = make(map[string]*Builtin)

// add a builtin command
func Add(cmd *Builtin) {
	builtins[cmd.Name] = cmd
}

func Get(name string) (*Builtin, bool) {
	b, ok := builtins[name]
	return b, ok
}

func All() []*Builtin {
	b := make([]*Builtin, len(builtins))
	for _, v := range builtins {
		b = append(b, v)
	}
	// TODO: sort
	return b
}

// Builtin describes a built-in command
type Builtin struct {
	// ExecFunc describes the function that executes the command.
	Execute func(cmd *Builtin, args []string)

	// Name is the name of the command, a string with no spaces, 
	// usually consistng of lowercase letters.
	Name string

	// UsageLine is the one-line usage message.
	UsageLine string

	// ManPage display's this command's man page.
	ManPage string

	// Flag is a set of flags specific to this command.
	FlagSet flag.FlagSet
}

func (c *Builtin) Usage() {
	// TODO: review
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.ManPage))
	os.Exit(2)
}
