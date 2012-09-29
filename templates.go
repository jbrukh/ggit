package main

import (
	"github.com/jbrukh/ggit/api"
	"io"
	"text/template"
)

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.Must(template.New("ggit").Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

// formatted messages
const (
	fmtUnknownCommand = "ggit: '%s' is not a ggit command. See 'ggit --help'.\n"
)

// constant messages
const (
	msgNotARepo = "fatal: Not a git repository (or any of the parent directories): " + api.DefaultGitDir
)

var tmplUsage = `usage: ggit [--version] <command> [<args>]

Available commands:
{{range .}}
   {{.Info.Name | printf "%-11s"}} {{.Info.Description }}{{end}}
`
