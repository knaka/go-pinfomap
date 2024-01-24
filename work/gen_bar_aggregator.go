// //go:build ignore
// // +build ignore

package main

import (
	"github.com/knaka/pkginfomapper"
)

func main() {
	structInfo, err := pkginfomapper.GetStructInfo(".", "Bar", map[string]any{
		"AdditionalComment": "This is a comment.",
	})
	if err != nil {
		panic(err)
	}
	err = pkginfomapper.GenerateForStruct(`
// Code generated by {{.GeneratorName}}; DO NOT EDIT.

package {{$.PackageName}}

import (
{{range $.Imports}}
	"{{.}}"
{{end}}
)

{{range $.Methods}}
func (recv {{if .PointerReceiver}}*{{end}}Baz) {{.Name}}() {
}
{{end}}{{/* range */}}
	`, structInfo)
	if err != nil {
		panic(err)
	}
}
