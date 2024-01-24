//go:build ignore
// +build ignore

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
	err = pkginfomapper.GenerateForStruct(pkginfomapper.AccessorsTemplate, structInfo)
	if err != nil {
		panic(err)
	}
}
