package expr_test

import (
	"go/types"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/protogodev/validate/expr"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name          string
		inStr         string
		inType        types.Type
		wantValidator expr.Validator
		wantErrStr    string
	}{
		{
			name:   "not",
			inStr:  "!lt(0)",
			inType: types.Typ[types.Int],
			wantValidator: &expr.LogicValidator{
				Qualifier: expr.DefaultQualifier,
				Name:      "Not",
				Left: &expr.LeafValidator{
					Qualifier: expr.DefaultQualifier,
					Name:      "Lt[int]",
					Args:      []string{"0"},
				},
			},
		},
		{
			name:   "range",
			inStr:  "xrange(0, 100)",
			inType: types.Typ[types.Int],
			wantValidator: &expr.LeafValidator{
				Qualifier: expr.DefaultQualifier,
				Name:      "Range[int]",
				Args:      []string{"0", "100"},
			},
		},
		{
			name:   "nonzero string",
			inStr:  "nonzero",
			inType: types.Typ[types.String],
			wantValidator: &expr.LeafValidator{
				Qualifier: expr.DefaultQualifier,
				Name:      "Nonzero[string]",
			},
		},
		{
			name:   "zero string",
			inStr:  "zero",
			inType: types.Typ[types.String],
			wantValidator: &expr.LeafValidator{
				Qualifier: expr.DefaultQualifier,
				Name:      "Zero[string]",
			},
		},
		{
			name:   "len string",
			inStr:  "len(0, 20).msg(\"bad length\") && match(`^\\w+$`)",
			inType: types.Typ[types.String],
			wantValidator: &expr.LogicValidator{
				Qualifier: expr.DefaultQualifier,
				Name:      "All",
				Left: &expr.LeafValidator{
					Qualifier: expr.DefaultQualifier,
					Name:      "LenString",
					Args:      []string{"0", "20"},
					Msg:       `"bad length"`,
				},
				Right: &expr.LeafValidator{
					Qualifier: expr.DefaultQualifier,
					Name:      "Match",
					Args:      []string{"`^\\w+$`"},
				},
			},
		},
		{
			name:   "len slice",
			inStr:  "len(0, 20).msg(\"bad length\")",
			inType: types.NewSlice(types.Typ[types.String]),
			wantValidator: &expr.LeafValidator{
				Qualifier: expr.DefaultQualifier,
				Name:      "LenSlice",
				Args:      []string{"0", "20"},
				Msg:       `"bad length"`,
			},
		},
		{
			name:       "match slice",
			inStr:      "match(`^\\w+$`)",
			inType:     types.NewSlice(types.Typ[types.String]),
			wantErrStr: "cannot use validator `match` on type *types.Slice",
			wantValidator: &expr.LeafValidator{
				Qualifier: expr.DefaultQualifier,
				Name:      "match",
				Args:      []string{"`^\\w+$`"},
			},
		},
		{
			name:       "match slice",
			inStr:      "match(`^\\w+$`)",
			inType:     types.NewSlice(types.Typ[types.String]),
			wantErrStr: "cannot use validator `match` on type *types.Slice",
			wantValidator: &expr.LeafValidator{
				Qualifier: expr.DefaultQualifier,
				Name:      "match",
				Args:      []string{"`^\\w+$`"},
			},
		},
		{
			name:  "underscore struct",
			inStr: "_",
			inType: newStruct([]*structField{
				{
					name: "Name",
					typ:  types.Typ[types.String],
				},
				{
					name: "Age",
					typ:  types.Typ[types.Int],
				},
			}),
			wantValidator: &expr.LeafValidator{
				Qualifier: expr.DefaultQualifier,
				Name:      "_.Schema",
			},
		},
		{
			name:   "complex",
			inStr:  "!nonzero || runecount(10, 20) && match(`^\\w+$`)",
			inType: types.Typ[types.String],
			wantValidator: &expr.LogicValidator{
				Qualifier: expr.DefaultQualifier,
				Name:      "Any",
				Left: &expr.LogicValidator{
					Qualifier: expr.DefaultQualifier,
					Name:      "Not",
					Left: &expr.LeafValidator{
						Qualifier: expr.DefaultQualifier,
						Name:      "Nonzero[string]",
					},
				},
				Right: &expr.LogicValidator{
					Qualifier: expr.DefaultQualifier,
					Name:      "All",
					Left: &expr.LeafValidator{
						Qualifier: expr.DefaultQualifier,
						Name:      "RuneCount",
						Args:      []string{"10", "20"},
					},
					Right: &expr.LeafValidator{
						Qualifier: expr.DefaultQualifier,
						Name:      "Match",
						Args:      []string{"`^\\w+$`"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValidator, err1 := expr.Parse(tt.inStr)
			var err2 error
			if gotValidator != nil {
				err2 = gotValidator.ConvertName(expr.DefaultConverters, tt.inType)
			}
			cmpError(t, err1, err2, tt.wantErrStr)

			//t.Logf("validator: %s", gotValidator.ExprString())

			if !cmp.Equal(gotValidator, tt.wantValidator) {
				diff := cmp.Diff(gotValidator, tt.wantValidator)
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
