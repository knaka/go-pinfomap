package main

import (
	"github.com/knaka/go-pinfomap"
	"github.com/knaka/go-pinfomap/examples"
	"github.com/knaka/go-pinfomap/generator"

	"log"
)

func main() {
	// Extracts information about a struct by passing an instance of the struct.
	// Here, we're using 'examples.Foo{}' with an option to include fields tagged as "accessor".
	structInfo, err := pinfomap.NewStructInfo(examples.Foo{}, pinfomap.WithTags("accessor"))
	if err != nil {
		log.Fatalf("failed to get struct info: %v", err)
	}

	// You can also specify a struct by its name to extract information.
	// Use just the name "Foo" for structs in the current package or provide the full path
	// like "github.com/knaka/foo/bar.baz" for the other structs.
	// This approach is useful for accessing structs dynamically or when they're not directly accessible.
	_, err = pinfomap.NewStructInfo("Foo", pinfomap.WithTags("accessor"))
	if err != nil {
		log.Fatalf("failed to get struct info: %v", err)
	}

	// Adding additional data to the 'structInfo' which can be utilized in the template.
	structInfo.Data = map[string]any{
		"AdditionalComment": "This is a comment.",
	}

	// Generates Go code based on the 'AccessorTemplate' and the provided 'structInfo'.
	err = generator.GenerateGo(pinfomap.AccessorTemplate, structInfo)
	if err != nil {
		log.Fatalf("failed to generate go: %v", err)
	}
}
