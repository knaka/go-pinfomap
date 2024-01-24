package main

import (
	goimports "github.com/incu6us/goimports-reviser/v3/reviser"
	"text/template"
)

type Foo struct {
	bar  string `pkginfomapper:"setter,getter=true"`
	Baz  int
	qux  template.Template   `pkginfomapper:"getter,setter=false"`
	quux goimports.SourceDir `pkginfomapper:"setter"`
}

func (rcv *Foo) Greet() string {
	return "Hello"
}

func (rcv Foo) Greet2() string {
	return "World"
}
