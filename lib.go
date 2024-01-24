package pkginfomapper

import (
	"bytes"
	goimports "github.com/incu6us/goimports-reviser/v3/reviser"
	"go/types"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
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
	log.Println("AFC701F", filepath.Base(generatorSourceDir))
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

func getOutputBasename() string {
	name := getGeneratorName()
	name = strings.TrimPrefix(name, "gen_")
	name = strings.TrimSuffix(name, "_go")
	name = name + ".go"
	return name
}

type Field struct {
	Name   string
	Type   string
	Tag    string
	Params map[string]string
}

func (f *Field) CapName() string {
	return strings.ToUpper(f.Name[0:1]) + f.Name[1:]
}

type Struct struct {
	StructName    string
	PackageName   string
	PackagePath   string
	GeneratorName string
	Fields        []*Field
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

func GetStructInfo(packageName string, structName string, data any) (struct_ *Struct, err error) {
	if callerFilename == "" {
		_, filename, _, ok := runtime.Caller(1)
		if !ok {
			panic("runtime.Caller() failed")
		}
		callerFilename = filename
	}
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedTypes | packages.NeedSyntax,
		Tests: false,
	}
	if strings.HasPrefix(packageName, ".") {
		packageName = filepath.Join(getOutputDir(), packageName)
	}
	packages_, err := packages.Load(cfg, packageName)
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
		Data:          data,
	}
	for _, name := range pkg.Types.Scope().Names() {
		st, ok := pkg.Types.Scope().Lookup(name).Type().Underlying().(*types.Struct)
		if !ok {
			continue
		}
		if name != structName {
			continue
		}
		x := st.String()
		log.Println(x)
		for i := 0; i < st.NumFields(); i++ {
			field := st.Field(i)
			fieldName := field.Name()
			type_ := field.Type()
			tag := st.Tag(i)
			fieldX := &Field{
				Name: fieldName,
				Type: types.TypeString(type_, func(p *types.Package) string {
					struct_.Imports = append(struct_.Imports, p.Path())
					return p.Name()
				}),
				Tag: tag,
			}
			struct_.Fields = append(struct_.Fields, fieldX)
			tagStr, ok := reflect.StructTag(tag).Lookup("pkginfomapper")
			if !ok {
				continue
			}
			Params := make(map[string]string)
			for _, ParamsStr := range strings.Split(tagStr, ",") {
				ParamsParts := strings.Split(ParamsStr, "=")
				if len(ParamsParts) == 1 {
					Params[ParamsParts[0]] = "true"
				} else if len(ParamsParts) == 2 {
					Params[ParamsParts[0]] = ParamsParts[1]
				} else {
					panic("Invalid tag info")
				}
			}
			fieldX.Params = Params
		}
	}
	return
}

func GenerateForStruct(tmpl string, info *Struct) (err error) {
	outPath := filepath.Join(getOutputDir(), getOutputBasename())
	if err != nil {
		return
	}
	tmplParsed := template.Must(template.New("aaa").Parse(tmpl))
	buf := new(bytes.Buffer)
	err = tmplParsed.Execute(buf, info)
	if err != nil {
		return
	}
	err = os.WriteFile(outPath, buf.Bytes(), 0644)
	if err != nil {
		return
	}
	sourceFile := goimports.NewSourceFile(info.PackagePath, outPath)
	_ = goimports.WithRemovingUnusedImports(sourceFile)
	fixed, _, differs, err := sourceFile.Fix()
	if err != nil {
		return
	}
	if differs {
		err = os.WriteFile(outPath, fixed, 0644)
		if err != nil {
			return
		}
	}
	return
}
