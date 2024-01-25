# pinfomap

## Description

Generates Go code from a template by referencing package information such as package names, imports, and struct fields.

## Examples

### Accessor Generation

`./foo.go` contains the struct `Foo`:

```go
package examples

import (
	goimports "github.com/incu6us/goimports-reviser/v3/reviser"
	"text/template"
)

type Foo struct {
	bar  string `accessor:"setter,getter=true"`
	Baz  int
	qux  template.Template   `accessor:"getter,setter=false"`
	quux goimports.SourceDir `accessor:"setter"`
}

func (rcv *Foo) Greet() string {
	return "Hello"
}

func (rcv Foo) Greet2() string {
	return "World"
}
```

`./gen_foo_accessor/main.go` serves as an accessor generator for the struct. You can place the generator comment `//go:generate go run ./gen_foo_accessor/` in the same directory as `./foo.go`:

```go
package main

import (
	"github.com/knaka/pinfomap"
	"github.com/knaka/pinfomap/examples"
)

func main() {
	structInfo, err := pinfomap.GetStructInfo(
		examples.Foo{},
		&pinfomap.GetStructInfoParams{
			Tags: []string{"accessor"},
			Data: map[string]any{
				"AdditionalComment": "This is a comment.",
			},
		},
	)
	if err != nil {
		panic(err)
	}
	err = pinfomap.Generate(pinfomap.AccessorTemplate, structInfo)
	if err != nil {
		panic(err)
	}
}
```

Running `go run ./gen_foo_accessor/` will generate the following code as `./foo_accessor.go`:

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

func (rcv *Foo) SetQuux(value reviser.SourceDir) {
	rcv.quux = value
}
```

Running `go run ./gen_foo_accessor.go` is also feasible if you prefer to place the generator code in the same directory as `./foo.go`. Be sure to add `//go:build ignore` to the top of the generator's code to prevent it from being compiled and linked with other code.

While the above example is using the template `AccessorTemplate`, you can use your own template by passing it to `pinfomap.Generate`. The `AccessorTemplate` is defined as follows:

```gotemplate
// Code generated by {{.GeneratorName}}; DO NOT EDIT.

package {{$.PackageName}}

import (
{{range $.Imports}}
	"{{.}}"
{{end}}
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
