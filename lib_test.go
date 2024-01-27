package pinfomap

import (
	"app/db/sqlcgen"
	. "github.com/knaka/go-utils"
	"github.com/knaka/pinfomap/examples"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestGetStructInfo(t *testing.T) {
	structInfo := Ensure(GetStructInfo(examples.Foo{}, &GetStructInfoParams{}))
	assert.Equal(t, structInfo.StructName, "Foo")
	assert.Equal(t, structInfo.Package.Name, "examples")
	assert.Equal(t, structInfo.Package.Path, "github.com/knaka/pinfomap/examples")
	for _, name := range structInfo.Package.Names() {
		log.Println(name)
	}
}

func TestGetStructInfo2(t *testing.T) {
	structInfo := Ensure(GetStructInfo(sqlcgen.User{}, &GetStructInfoParams{}))
	assert.Equal(t, structInfo.StructName, "User")
	assert.Equal(t, structInfo.Package.Name, "sqlcgen")
	assert.Equal(t, structInfo.Package.Path, "app/db/sqlcgen")
	for _, name := range structInfo.Package.Names() {
		log.Println(name)
	}
}
