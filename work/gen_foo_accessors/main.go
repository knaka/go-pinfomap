package main

import (
	"github.com/knaka/pinfomap"
)

func main() {
	structInfo, err := pinfomap.GetStructInfo(".", "Foo", &pinfomap.GetStructInfoParams{
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
