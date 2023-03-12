package validate

import (
	"fmt"
	"go/types"
	"path/filepath"
	"strings"

	protogocmd "github.com/protogodev/protogo/cmd"
	"github.com/protogodev/protogo/generator"
	"github.com/protogodev/protogo/parser"
	"github.com/protogodev/protogo/parser/ifacetool"
	"github.com/protogodev/validate/expr"
)

func init() {
	protogocmd.MustRegister(&protogocmd.Plugin{
		Name: "validate",
		Cmd:  protogocmd.NewGen(&Generator{}),
	})
}

type Generator struct {
	OutDir    string `name:"out" default:"." help:"output directory"`
	Formatted bool   `name:"fmt" default:"true" help:"whether to make the generated code formatted"`
}

func (g *Generator) PkgName() string {
	return parser.PkgNameFromDir(g.OutDir)
}

func (g *Generator) Generate(data *ifacetool.Data) (*generator.File, error) {
	tmplData := struct {
		Data *ifacetool.Data
	}{
		Data: data,
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

	return generator.Generate(Template, tmplData, generator.Options{
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
				if err := validator.ConvertName(expr.DefaultConverters, paramType); err != nil {
					panic(err)
				}

				validator.SetQualifier("v")
				s := validator.ExprString()
				// Replace the special validator `_` with the parameter's name.
				s = strings.Replace(s, "v._", paramName, 1)
				return s
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
