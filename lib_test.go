package pinfomap

import (
	"app/db/sqlcgen"
	"github.com/knaka/go-pinfomap/examples"
	. "github.com/knaka/go-utils"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestGetStructInfo(t *testing.T) {
	structInfo := Ensure(NewStructInfo(examples.Foo{}, &ctorParams{}))
	assert.Equal(t, structInfo.StructName, "Foo")
	assert.Equal(t, structInfo.Package.Name, "examples")
	assert.Equal(t, structInfo.Package.Path, "github.com/knaka/go-pinfomap/examples")
}

func TestGetStructInfo2(t *testing.T) {
	structInfo := Ensure(NewStructInfo(sqlcgen.User{}, &ctorParams{}))
	assert.Equal(t, structInfo.StructName, "User")
	assert.Equal(t, structInfo.Package.Name, "sqlcgen")
	assert.Equal(t, structInfo.Package.Path, "app/db/sqlcgen")
	for _, name := range structInfo.Package.GetStringTypes() {
		log.Println("1469FD4", name)
	}
}
