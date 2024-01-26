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
