package main

import (
	"github.com/knaka/pkginfomapper"
)

func main() {
	structInfo, err := pkginfomapper.GetStructInfo(".", "Foo", map[string]any{
		"AdditionalComment": "This is a comment.",
	})
	if err != nil {
		panic(err)
	}
	err = pkginfomapper.GenerateForStruct(pkginfomapper.AccessorsTemplate, structInfo)
	if err != nil {
		panic(err)
	}
}
