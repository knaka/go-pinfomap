package pinfomap

import (
	. "github.com/knaka/go-utils"
	"github.com/knaka/pinfomap/examples"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetStructInfo(t *testing.T) {
	structInfo := Ensure(GetStructInfo(examples.Foo{}, &GetStructInfoParams{}))
	assert.Equal(t, structInfo.StructName, "Foo")
	assert.Equal(t, structInfo.Package.Name, "examples")
	assert.Equal(t, structInfo.Package.Path, "github.com/knaka/pinfomap/examples")
}
