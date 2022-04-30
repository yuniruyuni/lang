package ir

import (
	"fmt"
	"strings"
	"text/template"
)

type IR string

func (ir IR) Expand(vals ...interface{}) IR {
	return IR(fmt.Sprintf(string(ir), vals...))
}

type IRs []IR

// Concat concatinate all given IRs into single IR.
func Concat(irs ...IR) IR {
	return Join("\n", irs...)
}

// Join concatinate all given IRs into single IR with separator.
func Join(sep string, irs ...IR) IR {
	strs := make([]string, 0, len(irs))
	for _, ir := range irs {
		strs = append(strs, string(ir))
	}
	return IR(strings.Join(strs, sep))
}

type Template struct {
	tmpl *template.Template
}

func T(tmpl string) *Template {
	return &Template{
		tmpl: template.Must(template.New("").Parse(tmpl)),
	}
}

func (t *Template) Expand(invs interface{}) IR {
	b := new(strings.Builder)
	if err := t.tmpl.Execute(b, invs); err != nil {
		panic("template expansion doesn't work properly.")
	}
	return IR(b.String())
}

type Vars map[string]interface{}
