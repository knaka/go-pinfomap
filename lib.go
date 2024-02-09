package pinfomap

import (
	"go/types"
	"golang.org/x/tools/go/packages"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type Field struct {
	Name   string
	Type   string
	Tag    string
	Params map[string]any
	Data   any
}

func (f *Field) CapName() string {
	return strings.ToUpper(f.Name[0:1]) + f.Name[1:]
}

// Lazily initialized variables
var v1 struct {
	reId   *regexp.Regexp
	reEach *regexp.Regexp
	once   sync.Once
}

func ForceCamel(sIn string) (s string) {
	v1.once.Do(func() {
		// Does not support look ahead/behind
		v1.reId = regexp.MustCompile(`(^|[a-z0-9])ID($|[A-Z])`)
		v1.reEach = regexp.MustCompile(`([A-Z][a-z0-9]*)`)
	})
	s = sIn
	s = v1.reId.ReplaceAllString(s, "${1}Id${2}")
	return s
}

// Camel2Snake converts a camel case string to snake case.
func Camel2Snake(sIn string) (s string) {
	v1.once.Do(func() {
		// Does not support look ahead/behind
		v1.reId = regexp.MustCompile(`(^|[a-z0-9])ID($|[A-Z])`)
		v1.reEach = regexp.MustCompile(`([A-Z][a-z0-9]*)`)
	})
	s = sIn
	s = v1.reId.ReplaceAllString(s, "${1}Id${2}")
	s = v1.reEach.ReplaceAllStringFunc(s, func(s string) string {
		return "_" + strings.ToLower(s)
	})
	s = strings.TrimPrefix(s, "_")
	return s
}

// SnakeName converts a field name to snake case.
func (f *Field) SnakeName() string {
	return Camel2Snake(f.Name)
}

func (f *Field) CamelName() string {
	return ForceCamel(f.Name)
}

func (f *Field) LowerCamelName() string {
	s := ForceCamel(f.Name)
	s = strings.ToLower(s[0:1]) + s[1:]
	return s
}

type Method struct {
	Name            string
	Type            string
	PointerReceiver bool
}

type Package struct {
	Name string
	Path string
	Pkg  *packages.Package
}

type Struct struct {
	StructName  string
	Package     *Package
	PackagePath string
	Fields      []*Field
	Methods     []*Method
	Imports     []string
	Data        any
}

func (s *Struct) SnakeStructName() string {
	return Camel2Snake(s.StructName)
}

func (s *Struct) PackageName() string {
	return s.Package.Name
}

func (p *Package) GetStringTypes() []string {
	var ret []string
	for _, name := range p.Pkg.Types.Scope().Names() {
		//ret = append(ret, name)
		x := p.Pkg.Types.Scope().Lookup(name)
		t := x.Type()
		//u := t.String()
		//log.Println(t, u)
		// Underlying type が types.Basic で、その BasicKind が `String` なものを列挙したい？ あるいは、もっと汎用的なメソッドの方が使いやすいか
		tu, ok := t.Underlying().(*types.Basic)
		if !ok {
			continue
		}
		// `type ... string` の型以外を排除
		if tu.Kind() != types.String {
			continue
		}
		ret = append(ret, name)
	}
	return ret
}

func (s *Struct) PrivateFields() []*Field {
	privateFields := make([]*Field, 0)
	for _, field := range s.Fields {
		if strings.ToLower(field.Name[0:1]) == field.Name[0:1] {
			privateFields = append(privateFields, field)
		}
	}
	return privateFields
}

type ctorParams struct {
	// Tag names to parse
	Tags []string
	// Additional data to pass to the template
	Data any
}

func (p *Package) GetTypes() []string {
	var ret []string
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedTypes | packages.NeedSyntax,
		Tests: false,
	}
	packages_, err := packages.Load(cfg, p.Path)
	if err != nil {
		panic(err)
	}
	pkg := packages_[0]
	for _, name := range pkg.Types.Scope().Names() {
		t := pkg.Types.Scope().Lookup(name).Type()
		// todo: filter
		ret = append(ret, t.String())
	}
	return ret
}

type InspectorOptSetter func(*ctorParams)

func WithTags(tags ...string) InspectorOptSetter {
	return func(params *ctorParams) {
		params.Tags = tags
	}
}

func WithData(data any) InspectorOptSetter {
	return func(params *ctorParams) {
		params.Data = data
	}
}

func NewStructInfo(tgtStruct any, optSetters ...InspectorOptSetter) (struct_ *Struct, err error) {
	var packagePath string
	var structName string
	if structPath, ok := tgtStruct.(string); ok {
		// github.com/knaka/go-pinfomap.Struct
		fields := strings.Split(structPath, ".")
		switch len(fields) {
		case 0:
			log.Fatalf("Invalid struct path: %v", structPath)
		case 1:
			// "foo" -> struct "foo" in the current (current directory) package
			packagePath = "."
			structName = fields[0]
		default:
			// "github.com/knaka/foo/bar.baz" -> struct "baz" in the package "github.com/knaka/foo/bar"
			packagePath = strings.Join(fields[0:len(fields)-1], ".")
			structName = fields[len(fields)-1]
		}
	} else {
		type_ := reflect.TypeOf(tgtStruct)
		if type_.Kind() != reflect.Struct {
			panic("Not a struct")
		}
		packagePath = type_.PkgPath()
		structName = type_.Name()
	}
	structInfo := &ctorParams{}
	for _, optSetter := range optSetters {
		optSetter(structInfo)
	}
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedTypes | packages.NeedSyntax,
		Tests: false,
	}
	packages_, err := packages.Load(cfg, packagePath)
	if err != nil {
		return
	}
	pkg := packages_[0]
	struct_ = &Struct{
		StructName: structName,
		Package: &Package{
			Name: pkg.Name,
			Path: pkg.PkgPath,
			Pkg:  pkg,
		},
		// Correct?
		PackagePath: pkg.PkgPath,
		Data:        structInfo.Data,
	}
	for _, name := range pkg.Types.Scope().Names() {
		t := pkg.Types.Scope().Lookup(name).Type()
		st, ok := t.Underlying().(*types.Struct)
		if !ok {
			continue
		}
		if name != structName {
			continue
		}
		for i := 0; i < st.NumFields(); i++ {
			field := st.Field(i)
			fieldName := field.Name()
			type_ := field.Type()
			structTag := st.Tag(i)
			fieldX := &Field{
				Name: fieldName,
				Type: types.TypeString(type_, func(p *types.Package) string {
					struct_.Imports = append(struct_.Imports, p.Path())
					return p.Name()
				}),
				Tag:    structTag,
				Params: map[string]any{},
			}
			struct_.Fields = append(struct_.Fields, fieldX)
			for _, tag := range structInfo.Tags {
				tagStr, ok := reflect.StructTag(structTag).Lookup(tag)
				if !ok {
					continue
				}
				for _, ParamsStr := range strings.Split(tagStr, ",") {
					ParamsParts := strings.Split(ParamsStr, "=")
					if len(ParamsParts) == 1 {
						fieldX.Params[ParamsParts[0]] = true
					} else if len(ParamsParts) == 2 {
						value := ParamsParts[1]
						var value_ any
						if value == "true" {
							value_ = true
						} else if value == "false" {
							value_ = false
						} else if n, err := strconv.Atoi(value); err == nil {
							value_ = n
						} else if n, err := strconv.ParseFloat(value, 64); err == nil {
							value_ = n
						} else {
							value_ = value
						}
						fieldX.Params[ParamsParts[0]] = value_
					} else {
						panic("Invalid tag info")
					}
				}
			}
		}
		methodSet2 := types.NewMethodSet(types.NewPointer(t))
		for i := 0; i < methodSet2.Len(); i++ {
			method := methodSet2.At(i)
			//log.Println(method.Obj().Name())
			sig := method.Type().(*types.Signature)
			//x := sig.Variadic() // 可変長 (`...`)
			//log.Println(x)
			//recv := sig.Recv()
			// Receiver type
			//log.Println("type:", recv.Type())
			//log.Println(recv.Pos())
			//structInfo := sig.Params()
			//log.Println(structInfo.Len())
			//for j := 0; j < structInfo.Len(); j++ {
			//	param := structInfo.At(j)
			//	//log.Println("name:", param.Name(), param.Type())
			//}
			results := sig.Results()
			//log.Println(results.Len())
			resultTypes := ""
			delim := ""
			for k := 0; k < results.Len(); k++ {
				result := results.At(k)
				//log.Println("name:", result.Name(), result.Type())
				resultTypes += delim + result.Type().String()
				delim = ", "
			}
			struct_.Methods = append(struct_.Methods, &Method{
				Name:            method.Obj().Name(),
				Type:            resultTypes,
				PointerReceiver: true,
			})
		}
		methodSet := types.NewMethodSet(t)
		for i := 0; i < methodSet.Len(); i++ {
			method := methodSet.At(i)
			//log.Println(method.Obj().Name())
			name := method.Obj().Name()
			for x, method := range struct_.Methods {
				if method.Name == name {
					struct_.Methods[x].PointerReceiver = false
				}
			}
		}
	}
	return
}
