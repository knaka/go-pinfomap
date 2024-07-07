package stringer

import (
	"fmt"
	"go/ast"
)

type DtoValue struct {
	Name  string
	Value uint64
}

func NewIntTypeInfo(packageDirOrPath string, typeName string) (dtoValues []DtoValue, err error) {
	var tags []string
	g := Generator{}
	g.parsePackage([]string{packageDirOrPath}, tags)
	if g.pkg == nil {
		err = fmt.Errorf("no package found for %s", packageDirOrPath)
		return
	}
	values := make([]Value, 0, 100)
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.typeName = typeName
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			values = append(values, file.values...)
		}
	}
	if len(values) == 0 {
		err = fmt.Errorf("no values found for %s", typeName)
		return
	}
	for _, value := range values {
		dtoValues = append(dtoValues, DtoValue{Name: value.name, Value: value.value})
	}
	return
}
