package main

import (
	"io"
	"text/template"
)

const unknownCommandFormat = "ggit: '%s' is not a ggit command. See 'ggit --help'.\n"

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.Must(template.New("ggit").Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

var usageTemplate = `usage: ggit [--version] <command> [<args>]

Available commands:
{{range .}}
   {{.Name | printf "%-11s"}} {{.Description }}{{end}}
`
