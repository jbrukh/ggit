package main

import (
	"io"
	"strings"
	"text/template"
)

const unknownCommandFormat = `ggit: '%s' is not a ggit command. See 'ggit --help'.`

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

var usageTemplate = `Usage:

    ggit [commandname] [arg1] [arg2] [...]

Available commands:
{{range .}}{{if .Runnable}}
    {{.Name}}: {{.Short}}{{end}}{{end}}
`

var helpTemplate = `{{if .Runnable}}usage: go {{.UsageLine}}

{{end}}{{.Long | trim}}
`
