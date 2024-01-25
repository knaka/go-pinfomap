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
