package expr_test

import (
	"go/types"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/protogodev/validate/decl"
	"github.com/protogodev/validate/expr"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name           string
		inStr          string
		inParam        expr.Param
		wantExprString string
		wantErrStr     string
	}{
		/*{
			name:  "not",
			inStr: "!lt(0)",
			inParam: expr.Param{
				Name: "x",
				Type: types.Typ[types.Int],
			},
			wantExprString: "v.Not(v.Lt[int](0))",
		},
		{
			name:  "range",
			inStr: "xrange(0, 100)",
			inParam: expr.Param{
				Name: "x",
				Type: types.Typ[types.Int],
			},
			wantExprString: "v.Range[int](0, 100)",
		},
		{
			name:  "nonzero string",
			inStr: "nonzero",
			inParam: expr.Param{
				Name: "x",
				Type: types.Typ[types.String],
			},
			wantExprString: "v.Nonzero[string]()",
		},
		{
			name:  "zero string",
			inStr: "zero",
			inParam: expr.Param{
				Name: "x",
				Type: types.Typ[types.String],
			},
			wantExprString: "v.Zero[string]()",
		},
		{
			name:  "len string",
			inStr: "len(0, 20).msg(\"bad length\") && match(`^\\w+$`)",
			inParam: expr.Param{
				Name: "x",
				Type: types.Typ[types.String],
			},
			wantExprString: "v.All(v.LenString(0, 20).Msg(\"bad length\"), v.Match(regexp.MustCompile(`^\\w+$`)))",
		},*/
		{
			name:  "len slice",
			inStr: "len(0, 20).msg(\"bad length\")",
			inParam: expr.Param{
				Name: "x",
				Type: types.NewSlice(types.Typ[types.String]),
			},
			wantExprString: "v.LenSlice[[]string](0, 20).Msg(\"bad length\")",
		}, /*
			{
				name:  "match slice",
				inStr: "match(`^\\w+$`)",
				inParam: expr.Param{
					Name: "x",
					Type: types.NewSlice(types.Typ[types.String]),
				},
				wantErrStr:     "cannot use validator `match` on type *types.Slice",
				wantExprString: "",
			},
			{
				name:  "match slice",
				inStr: "match(`^\\w+$`)",
				inParam: expr.Param{
					Name: "x",
					Type: types.NewSlice(types.Typ[types.String]),
				},
				wantErrStr:     "cannot use validator `match` on type *types.Slice",
				wantExprString: "",
			},
			{
				name:  "underscore struct",
				inStr: "_",
				inParam: expr.Param{
					Name: "x",
					Type: newStruct([]*structField{
						{
							name: "Name",
							typ:  types.Typ[types.String],
						},
						{
							name: "Age",
							typ:  types.Typ[types.Int],
						},
					}),
				},
				wantExprString: "x.Schema()",
			},
			{
				name:  "complex",
				inStr: "!nonzero || runecnt(10, 20) && match(`^\\w+$`)",
				inParam: expr.Param{
					Name: "x",
					Type: types.Typ[types.String],
				},
				wantExprString: "v.Any(v.Not(v.Nonzero[string]()), v.All(v.RuneCount(10, 20), v.Match(regexp.MustCompile(`^\\w+$`))))",
			},*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builtin, err := decl.Parse(decl.BuiltinDecls)
			if err != nil {
				t.Fatalf("err: %v\n", err)
			}

			decls := make(map[string][]*decl.Validator)
			for _, d := range builtin {
				decls[d.Alias] = append(decls[d.Alias], d)
			}

			validator, err1 := expr.Parse(tt.inStr)
			//t.Logf("validator: %s", gotValidator.ExprString())

			var err2 error
			if validator != nil {
				err2 = validator.Bind(tt.inParam, decls)
			}
			cmpError(t, err1, err2, tt.wantErrStr)

			var gotExprString string
			if err1 == nil && err2 == nil {
				gotExprString = validator.ExprString()
			}

			if !cmp.Equal(gotExprString, tt.wantExprString) {
				diff := cmp.Diff(gotExprString, tt.wantExprString)
				t.Errorf("Want - Got: %s", diff)
			}
		})
	}
}

func cmpError(t *testing.T, err1, err2 error, wantErrStr string) {
	switch {
	case err1 != nil:
		if err1.Error() != wantErrStr {
			t.Errorf("Err: got (%#v), want (%#v)", err1.Error(), wantErrStr)
		}
	case err2 != nil:
		if err2.Error() != wantErrStr {
			t.Errorf("Err: got (%#v), want (%#v)", err2.Error(), wantErrStr)
		}
	default:
		if wantErrStr != "" {
			t.Errorf("Err: got (%#v), want (%#v)", "", wantErrStr)
		}
	}
}

type structField struct {
	name string
	typ  types.Type
	tag  string
}

func newStruct(fields []*structField) *types.Struct {
	var fs []*types.Var
	var tags []string
	for _, f := range fields {
		fs = append(fs, types.NewField(0, nil, f.name, f.typ, false))
		tags = append(tags, f.tag)
	}
	return types.NewStruct(fs, tags)
}
