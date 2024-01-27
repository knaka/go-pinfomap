package pinfomap

import (
	"bytes"
	goimports "github.com/incu6us/goimports-reviser/v3/reviser"
	"go/types"
	"golang.org/x/tools/go/packages"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

var callerFilename string

func Init() {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("runtime.Caller() failed")
	}
	callerFilename = filename
}

func getGeneratorSourceFile() string {
	return callerFilename
}

func getGeneratorSourceDir() string {
	return filepath.Dir(getGeneratorSourceFile())
}

func getOutputDir() string {
	generatorSourceDir := getGeneratorSourceDir()
	if strings.HasPrefix(filepath.Base(generatorSourceDir), "gen_") {
		return filepath.Dir(generatorSourceDir)
	}
	generatorSourceFile := getGeneratorSourceFile()
	if strings.HasPrefix(filepath.Base(generatorSourceFile), "gen_") {
		return filepath.Dir(generatorSourceFile)
	}
	panic("Cannot infer output dir")
}

func getGeneratorName() string {
	name := ""
	generatorSourceDir := getGeneratorSourceDir()
	if strings.HasPrefix(filepath.Base(generatorSourceDir), "gen_") {
		name = filepath.Base(generatorSourceDir)
	} else {
		generatorSourceFile := getGeneratorSourceFile()
		if strings.HasPrefix(filepath.Base(generatorSourceFile), "gen_") {
			name = filepath.Base(generatorSourceFile)
		} else {
			panic("Cannot infer output basename")
		}
	}
	return name
}

var reExt = regexp.MustCompile(`_([[:alnum:]]+)$`)

func getOutputBasename() string {
	name := getGeneratorName()
	name = strings.TrimPrefix(name, "gen_")
	name = strings.TrimSuffix(name, ".go")
	name = reExt.ReplaceAllString(name, ".${1}")
	return name
}

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

var vars1 struct {
	reId   *regexp.Regexp
	reEach *regexp.Regexp
	once   sync.Once
}

func Camel2Snake(sIn string) (sOut string) {
	vars1.once.Do(func() {
		// Does not support look ahead/behind
		vars1.reId = regexp.MustCompile(`(^|[a-z0-9])ID($|[A-Z])`)
		vars1.reEach = regexp.MustCompile(`([A-Z][a-z0-9]*)`)
	})
	sOut = sIn
	sOut = vars1.reId.ReplaceAllString(sOut, "${1}Id${2}")
	sOut = vars1.reEach.ReplaceAllStringFunc(sOut, func(s string) string {
		return "_" + strings.ToLower(s)
	})
	sOut = strings.TrimPrefix(sOut, "_")
	return sOut
}

// SnakeCaseName converts a field name to snake case.
func (f *Field) SnakeCaseName() string {
	return Camel2Snake(f.Name)
}

type Method struct {
	Name            string
	Type            string
	PointerReceiver bool
}

type Struct struct {
	StructName    string
	PackageName   string
	PackagePath   string
	GeneratorName string
	Fields        []*Field
	Methods       []*Method
	Imports       []string
	Data          any
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

type GetStructInfoParams struct {
	// Tag names to parse
	Tags []string
	Data any
}

func initCallerFileName() {
	if callerFilename != "" {
		return
	}
	_, filename, _, ok := runtime.Caller(2)
	if !ok {
		panic("runtime.Caller() failed")
	}
	callerFilename = filename
}

func GetStructInfo(structObject any, params *GetStructInfoParams) (struct_ *Struct, err error) {
	initCallerFileName()
	type_ := reflect.TypeOf(structObject)
	if type_.Kind() != reflect.Struct {
		panic("Not a struct")
	}
	packagePath := type_.PkgPath()
	structName := type_.Name()
	return GetStructInfoByName(packagePath, structName, params)
}

func GetStructInfoByName(packagePath string, structName string, params *GetStructInfoParams) (struct_ *Struct, err error) {
	initCallerFileName()
	if params == nil {
		params = &GetStructInfoParams{}
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
		StructName:  structName,
		PackageName: pkg.Name,
		// Correct?
		PackagePath:   pkg.PkgPath,
		GeneratorName: getGeneratorName(),
		Data:          params.Data,
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
			//isPointer := (func() bool {
			//	_, ok := type_.(*types.Pointer)
			//	return ok
			//})()
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
			for _, tag := range params.Tags {
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
			//params := sig.Params()
			//log.Println(params.Len())
			//for j := 0; j < params.Len(); j++ {
			//	param := params.At(j)
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

type GenerateParams struct {
	DisableGoimports bool
}

func Generate(tmpl string, data any, params *GenerateParams) (err error) {
	if params == nil {
		params = &GenerateParams{}
	}
	outPath := filepath.Join(getOutputDir(), getOutputBasename())
	if err != nil {
		return
	}
	tmplParsed := template.Must(template.New("aaa").Parse(tmpl))
	buf := new(bytes.Buffer)
	err = tmplParsed.Execute(buf, data)
	if err != nil {
		return
	}
	err = os.WriteFile(outPath, buf.Bytes(), 0644)
	if err != nil {
		return
	}
	if !params.DisableGoimports {
		sourceFile := goimports.NewSourceFile( /* data.PackagePath */ "?", outPath)
		_ = goimports.WithRemovingUnusedImports(sourceFile)
		fixed, _, differs, err := sourceFile.Fix()
		if err != nil {
			return err
		}
		if differs {
			err = os.WriteFile(outPath, fixed, 0644)
			if err != nil {
				return err
			}
		}
	}
	return
}
