package generator

import (
	"bytes"
	goimports "github.com/incu6us/goimports-reviser/v3/reviser"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"text/template"
)

var callerFilename string

func initLib() {
	if callerFilename != "" {
		return
	}
	_, filename, _, ok := runtime.Caller(2)
	if !ok {
		panic("runtime.Caller() failed")
	}
	callerFilename = filename
}

var reExt = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`_([[:alnum:]]+)$`)
})

func getOutputBasename() string {
	name := getGeneratorName()
	name = strings.TrimPrefix(name, "gen_")
	name = strings.TrimSuffix(name, ".go")
	name = reExt().ReplaceAllString(name, ".${1}")
	return name
}

func getGeneratorSourceFile() string {
	return callerFilename
}

func getGeneratorSourceDir() string {
	return filepath.Dir(getGeneratorSourceFile())
}

func OutputDir() string {
	initLib()
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

type GenerateParams struct {
	ShouldRunGoImports bool
	Filename           string
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
			name = "Cannot infer output basename"
		}
	}
	return name
}

//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func GenerateGo(tmpl string, data any, optSetters ...GeneratorOptSetter) (err error) {
	initLib()
	return Generate(tmpl, data, append(optSetters, FormatsGo(true))...)
}

type GeneratorOptSetter func(*GenerateParams)

func FormatsGo(b bool) GeneratorOptSetter {
	return func(params *GenerateParams) {
		params.ShouldRunGoImports = b
	}
}

func WithFilename(filename string) GeneratorOptSetter {
	return func(params *GenerateParams) {
		params.Filename = filename
	}
}

func Generate(tmpl string, data any, optSetters ...GeneratorOptSetter) (err error) {
	initLib()
	params := &GenerateParams{}
	for _, setter := range optSetters {
		if setter == nil {
			continue
		}
		setter(params)
	}
	outPath := filepath.Join(OutputDir(), getOutputBasename())
	tmplParsed := template.Must(template.New("aaa").Funcs(map[string]any{
		"GeneratorName": getGeneratorName,
	}).Parse(tmpl))
	buf := new(bytes.Buffer)
	err = tmplParsed.Execute(buf, data)
	if err != nil {
		return
	}
	outDir, err := filepath.Abs(filepath.Dir(outPath))
	if err != nil {
		return
	}
	err = os.MkdirAll(outDir, 0755)
	if err != nil {
		return
	}
	err = os.WriteFile(outPath, buf.Bytes(), 0644)
	if err != nil {
		return
	}
	if params.ShouldRunGoImports {
		//cmd := exec.Command("go", "fmt", outPath)
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		//err = cmd.Run()
		//if err != nil {
		//	// Rename the file not to be used as a source file.
		//	_ = os.Rename(outPath, outPath+".err")
		//	return err
		//}
		sourceFile := goimports.NewSourceFile( /* data.PackagePath */ "?", outPath)
		_ = goimports.WithRemovingUnusedImports(sourceFile)
		_ = goimports.WithCodeFormatting(sourceFile)
		fixed, _, differs, err := sourceFile.Fix()
		if err != nil {
			// Rename the file not to be used as a source file.
			_ = os.Rename(outPath, outPath+".err")
			return err
		}
		if differs {
			err = os.WriteFile(outPath, fixed, 0644)
			if err != nil {
				return err
			}
		}
		_ = os.Remove(outPath + ".err")
	}
	return
}

func GetOutputPath(elem ...string) string {
	initLib()
	return filepath.Join(OutputDir(), filepath.Join(elem...))
}
