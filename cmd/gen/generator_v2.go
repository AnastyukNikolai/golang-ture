package main

import (
	"fmt"
	"go/types"
	golang_ture "golang-ture"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/packages"
)

type data struct {
	PropName   string
	PropType   string
	StructName string
}

func main() {
	// Handle arguments to command
	if len(os.Args) != 2 {
		failErr(fmt.Errorf("expected exactly one argument: <source type>"))
	}
	sourceType := os.Args[1]
	sourceTypePackage, sourceTypeName := splitSourceType(sourceType)

	// Inspect package and use type checker to infer imported types
	pkg := loadPackage(sourceTypePackage)

	// Lookup the given source type name in the package declarations
	obj := pkg.Types.Scope().Lookup(sourceTypeName)
	if obj == nil {
		failErr(fmt.Errorf("%s not found in declared types of %s",
			sourceTypeName, pkg))
	}

	// We check if it is a declared type
	if _, ok := obj.(*types.TypeName); !ok {
		failErr(fmt.Errorf("%v is not a named type", obj))
	}
	// We expect the underlying type to be a struct
	structType, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		failErr(fmt.Errorf("type %v is not a struct", obj))
	}

	if err := golang_ture.LoadGetSetGenTemplate(); err != nil {
		logrus.Fatalf("error detected while load templates: %s", err.Error())
	}

	// Generate code using jennifer
	err := generateV2(sourceTypeName, structType)
	if err != nil {
		failErr(err)
	}
}

func generateV2(sourceTypeName string, structType *types.Struct) error {
	goFile := os.Getenv("GOFILE")
	ext := filepath.Ext(goFile)
	baseFilename := goFile[0 : len(goFile)-len(ext)]
	targetFilename := baseFilename + "_" + strings.ToLower(sourceTypeName) + "_gen_v2.go"
	file, err := os.Create(targetFilename)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString("package " + os.Getenv("GOPACKAGE") + "\n\n")

	if err != nil {
		return err
	}

	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		d := data{StructName: sourceTypeName, PropName: field.Name(), PropType: field.Type().String()}
		tmpl, ok := golang_ture.Templates["get_set_gen.tmpl"]
		if !ok {
			panic("template get_set_gen.tmpl not found")
		}
		err = tmpl.Execute(file, d)
		if err != nil {
			return err
		}
	}

	return nil
}

//var getSetTemplate = `
//func (m *{{.StructName}}) Get{{.PropName}}() {{.PropType}} {
//	return m.{{.PropName}}
//}
//
//func (m *{{.StructName}}) Set{{.PropName}}(val {{.PropType}}) {
//	m.{{.PropName}} = val
//}
//`

func loadPackage(path string) *packages.Package {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedImports}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		failErr(fmt.Errorf("loading packages for inspection: %v", err))
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	return pkgs[0]
}

func splitSourceType(sourceType string) (string, string) {
	idx := strings.LastIndexByte(sourceType, '.')
	if idx == -1 {
		failErr(fmt.Errorf(`expected qualified type as "pkg/path.MyType"`))
	}
	sourceTypePackage := sourceType[0:idx]
	sourceTypeName := sourceType[idx+1:]
	return sourceTypePackage, sourceTypeName
}

func failErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
