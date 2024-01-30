package main

import (
	"github.com/knaka/go-pinfomap"
	"github.com/knaka/go-pinfomap/examples"
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

	err = pinfomap.GenerateGo(pinfomap.AccessorTemplate, structInfo, nil)
	if err != nil {
		panic(err)
	}
}
