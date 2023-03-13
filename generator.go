package validate

import (
	_ "embed"
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	protogocmd "github.com/protogodev/protogo/cmd"
	"github.com/protogodev/protogo/generator"
	"github.com/protogodev/protogo/parser"
	"github.com/protogodev/protogo/parser/ifacetool"
	"github.com/protogodev/validate/decl"
	"github.com/protogodev/validate/expr"
)

//go:embed template.go.tmpl
var template string

func init() {
	protogocmd.MustRegister(&protogocmd.Plugin{
		Name: "validate",
		Cmd:  protogocmd.NewGen(&Generator{}),
	})
}

type Generator struct {
	OutDir    string `name:"out" default:"." help:"output directory"`
	Formatted bool   `name:"fmt" default:"true" help:"whether to make the generated code formatted"`
	Custom    string `name:"custom" help:"the declaration file of custom validators"`
}

func (g *Generator) PkgName() string {
	return parser.PkgNameFromDir(g.OutDir)
}

func (g *Generator) Generate(data *ifacetool.Data) (*generator.File, error) {
	customDecls, err := getCustomDecls(g.Custom)
	if err != nil {
		return nil, err
	}

	completeDecls, imports := buildCompleteDecls(customDecls)

	tmplData := struct {
		Imports []ifacetool.Import
		Data    *ifacetool.Data
	}{
		Imports: imports,
		Data:    data,
	}

	schemas := make(map[string]map[string]string)
	for _, method := range data.Methods {
		m := make(map[string]string)
		options := ParseDoc(method.Doc)["schema"]
		for _, opt := range options {
			m[opt.K] = opt.V
		}
		schemas[method.Name] = m
	}

	return generator.Generate(template, tmplData, generator.Options{
		Funcs: map[string]interface{}{
			"nonCtxParams": func(params []*ifacetool.Param) (out []*ifacetool.Param) {
				for _, p := range params {
					if p.TypeString != "context.Context" {
						out = append(out, p)
					}
				}
				return
			},
			"methodSchema": func(methodName string) map[string]string {
				return schemas[methodName]
			},
			"exprString": func(schema, paramName string, paramType types.Type) string {
				validator, err := expr.Parse(schema)
				if err != nil {
					panic(err)
				}

				param := expr.Param{
					Name: paramName,
					Type: paramType,
				}
				if err := validator.Bind(param, completeDecls); err != nil {
					panic(err)
				}

				return validator.ExprString()
			},
			"returnErr": func(params []*ifacetool.Param, errFormat string) string {
				var returns []string
				for i := 0; i < len(params)-1; i++ {
					returns = append(returns, emptyValue(params[i]))
				}

				returns = append(returns, fmt.Sprintf(errFormat, "err"))
				return strings.Join(returns, ", ")
			},
		},
		Formatted:      g.Formatted,
		TargetFileName: filepath.Join(g.OutDir, "validate_gen.go"),
	})
}

func getCustomDecls(filename string) (string, error) {
	if filename == "" {
		return "", nil
	}

	b, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func buildCompleteDecls(customDecls string) (map[string][]*decl.Validator, []ifacetool.Import) {
	builtin, err := decl.Parse(decl.BuiltinDecls)
	if err != nil {
		panic(err)
	}

	custom, err := decl.Parse(customDecls)
	if err != nil {
		panic(err)
	}

	decls := make(map[string][]*decl.Validator)
	imports := make(map[string]string)

	for _, d := range builtin {
		decls[d.Alias] = append(decls[d.Alias], d)
		imports[d.Qualifier] = d.Import
	}
	for _, d := range custom {
		decls[d.Alias] = append(decls[d.Alias], d)
		imports[d.Qualifier] = d.Import
	}

	var importList []ifacetool.Import
	for alias, path := range imports {
		importList = append(importList, ifacetool.Import{
			Alias: alias,
			Path:  path,
		})

	}

	return decls, importList
}

func emptyValue(param *ifacetool.Param) string {
	t := param.Type.Underlying()

	switch v := t.(type) {
	case *types.Basic:
		switch info := v.Info(); {
		case info&types.IsInteger == types.IsInteger:
			return "0"
		case info&types.IsFloat == types.IsFloat:
			return "0"
		case info&types.IsString == types.IsString:
			return `""`
		case info&types.IsBoolean == types.IsBoolean:
			return "false"
		default:
			return `""`
		}
	case *types.Map, *types.Chan, *types.Slice, *types.Array, *types.Pointer, *types.Interface:
		return "nil"
	case *types.Struct:
		return param.TypeString + "{}"
	default:
		return "nil"
	}
}
