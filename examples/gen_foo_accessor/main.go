package main

import (
	"github.com/knaka/pinfomap"
	examples "github.com/knaka/pinfomap/examples"
	"reflect"
)

func main() {
	structInfo, err := pinfomap.GetStructInfo(
		reflect.TypeOf(examples.Foo{}).PkgPath(),
		reflect.TypeOf(examples.Foo{}).Name(),
		&pinfomap.GetStructInfoParams{
			Tags: []string{"accessor"},
			Data: map[string]any{
				"AdditionalComment": "This is a comment.",
			},
		})
	if err != nil {
		panic(err)
	}
	err = pinfomap.Generate(pinfomap.AccessorsTemplate, structInfo)
	if err != nil {
		panic(err)
	}
}
