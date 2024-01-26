# pinfomap

## Description

Generates Go code from a template by referencing package information such as package names, imports, and struct fields.

## Examples

### Accessor Generation

`./foo.go` contains the sample struct `Foo`:

<!-- mdppcode src=examples/foo.go -->
```go
package examples

import (
	goimports "github.com/incu6us/goimports-reviser/v3/reviser"
	"text/template"
)

type Foo struct {
	bar  string `accessor:"setter,getter=true"`
	Baz  int
	qux  template.Template    `accessor:"getter,setter=false"`
	quux *goimports.SourceDir `accessor:"setter"`
}
```

`./gen_foo_accessor/main.go` serves as an accessor generator for the struct. You can place the generator comment `//go:generate go run ./gen_foo_accessor/` in the same directory as `./foo.go`:

<!-- mdppcode src=examples/gen_foo_accessor/main.go -->
```go
package main

import (
	"github.com/knaka/pinfomap"
	"github.com/knaka/pinfomap/examples"
)

func main() {
	// GetStructInfo extracts information about the struct of the object passed as an argument.
	structInfo, err := pinfomap.GetStructInfo(examples.Foo{},
		&pinfomap.GetStructInfoParams{
			// Tags to consider when extracting information. In this case, it looks for "accessor" tag.
			Tags: []string{"accessor"},
		},
	)

	// Alternatively, specify the package and the struct by name to extract its information.
	// This method might be useful when the struct is not directly accessible or if you want to reference it dynamically.

	//structInfo, err := pinfomap.GetStructInfoByName(".", "Foo", ...)

	if err != nil {
		panic(err)
	}

	// Additional data that can be passed to the template.
	structInfo.Data = map[string]any{
		"AdditionalComment": "This is a comment.",
	}

	err = pinfomap.Generate(pinfomap.AccessorTemplate, structInfo, nil)
	if err != nil {
		panic(err)
	}
}
```

Running `go run ./gen_foo_accessor/` will generate the following code as `./foo_accessor.go`:

<!-- mdppcode src=examples/foo_accessor.go -->
```go
// Code generated by gen_foo_accessor; DO NOT EDIT.

package examples

import (
	"text/template"

	"github.com/incu6us/goimports-reviser/v3/reviser"
)

func (rcv *Foo) GetBar() string {
	return rcv.bar
}

func (rcv *Foo) SetBar(value string) {
	rcv.bar = value
}

func (rcv *Foo) GetQux() template.Template {
	return rcv.qux
}

func (rcv *Foo) SetQuux(value *reviser.SourceDir) {
	rcv.quux = value
}
```

Running `go run ./gen_foo_accessor.go` is also feasible if you prefer to place the generator code in the same directory as `./foo.go`. Be sure to add `//go:build ignore` to the top of the generator's code to prevent it from being compiled and linked with other code.

While the above example is using the template `AccessorTemplate`, you can use your own template by passing it to `pinfomap.Generate`. The `AccessorTemplate` is defined as follows:

<!-- mdppcode src=./templates/accessor.tmpl -->
```gotemplate
// Code generated by {{.GeneratorName}}; DO NOT EDIT.

package {{$.PackageName}}

import (
{{range $.Imports}}
	"{{.}}"
{{end}}{{/* range */}}
)

{{range $.PrivateFields}}
{{if .Params.getter }}
func (rcv *{{$.StructName}}) Get{{.CapName}}() {{.Type}} {
	return rcv.{{.Name}}
}
{{end}}{{/* if */}}

{{if .Params.setter }}
func (rcv *{{$.StructName}}) Set{{.CapName}}(value {{.Type}}) {
	rcv.{{.Name}} = value
}
{{end}}{{/* if */}}
{{end}}{{/* range */}}
```

The resulting Go code is automatically formatted with [goimports-reviser](https://github.com/incu6us/goimports-reviser), allowing you to write the template without worrying about the strictness of import statements or unnecessary empty lines.
