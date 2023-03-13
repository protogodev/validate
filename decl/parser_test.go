package decl_test

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/protogodev/validate/decl"
)

func TestParseConverters(t *testing.T) {
	got, err := decl.Parse(decl.BuiltinDecls)
	if err != nil {
		t.Errorf("err: %v\n", err)
	}

	want := []*decl.Validator{
		{
			Import:       "github.com/RussellLuo/validating",
			Qualifier:    "v",
			Name:         "Nonzero",
			IsGeneric:    true,
			Alias:        "nonzero",
			AllowedTypes: []string{"comparable"},
			ArgNum:       decl.Range{Min: 0, Max: 0},
		},
		{
			Import:       "github.com/RussellLuo/validating",
			Qualifier:    "v",
			Name:         "Zero",
			IsGeneric:    true,
			Alias:        "zero",
			AllowedTypes: []string{"comparable"},
			ArgNum:       decl.Range{Min: 0, Max: 0},
		},
		{
			Import:       "github.com/RussellLuo/validating",
			Qualifier:    "v",
			Name:         "LenString",
			IsGeneric:    false,
			Alias:        "len",
			AllowedTypes: []string{"string"},
			ArgNum:       decl.Range{Min: 2, Max: 2},
		},
		{
			Import:       "github.com/RussellLuo/validating",
			Qualifier:    "v",
			Name:         "LenSlice",
			IsGeneric:    true,
			Alias:        "len",
			AllowedTypes: []string{"slice"},
			ArgNum:       decl.Range{Min: 2, Max: 2},
		},
		{
			Import:       "github.com/RussellLuo/validating",
			Qualifier:    "v",
			Name:         "RuneCount",
			IsGeneric:    false,
			Alias:        "runec",
			AllowedTypes: []string{"string|bytes"},
			ArgNum:       decl.Range{Min: 2, Max: 2},
		},
		{
			Import:       "github.com/RussellLuo/validating",
			Qualifier:    "v",
			Name:         "Eq",
			IsGeneric:    true,
			Alias:        "eq",
			AllowedTypes: []string{"comparable"},
			ArgNum:       decl.Range{Min: 1, Max: 1},
		},
		{
			Import:       "github.com/RussellLuo/validating",
			Qualifier:    "v",
			Name:         "Ne",
			IsGeneric:    true,
			Alias:        "ne",
			AllowedTypes: []string{"comparable"},
			ArgNum:       decl.Range{Min: 1, Max: 1},
		},
		{
			Import:       "github.com/RussellLuo/validating",
			Qualifier:    "v",
			Name:         "Gt",
			IsGeneric:    true,
			Alias:        "gt",
			AllowedTypes: []string{"ordered"},
			ArgNum:       decl.Range{Min: 1, Max: 1},
		},
		{
			Import:       "github.com/RussellLuo/validating",
			Qualifier:    "v",
			Name:         "Gte",
			IsGeneric:    true,
			Alias:        "gte",
			AllowedTypes: []string{"ordered"},
			ArgNum:       decl.Range{Min: 1, Max: 1},
		},
		{
			Import:       "github.com/RussellLuo/validating",
			Qualifier:    "v",
			Name:         "Lt",
			IsGeneric:    true,
			Alias:        "lt",
			AllowedTypes: []string{"ordered"},
			ArgNum:       decl.Range{Min: 1, Max: 1},
		},
		{
			Import:       "github.com/RussellLuo/validating",
			Qualifier:    "v",
			Name:         "Lte",
			IsGeneric:    true,
			Alias:        "lte",
			AllowedTypes: []string{"ordered"},
			ArgNum:       decl.Range{Min: 1, Max: 1},
		},
		{
			Import:       "github.com/RussellLuo/validating",
			Qualifier:    "v",
			Name:         "Range",
			IsGeneric:    true,
			Alias:        "xrange",
			AllowedTypes: []string{"ordered"},
			ArgNum:       decl.Range{Min: 2, Max: 2},
		},
		{
			Import:       "github.com/RussellLuo/validating",
			Qualifier:    "v",
			Name:         "In",
			IsGeneric:    true,
			Alias:        "in",
			AllowedTypes: []string{"ordered"},
			ArgNum:       decl.Range{Min: 1, Max: math.MaxInt},
		},
		{
			Import:       "github.com/RussellLuo/validating",
			Qualifier:    "v",
			Name:         "Nin",
			IsGeneric:    true,
			Alias:        "nin",
			AllowedTypes: []string{"ordered"},
			ArgNum:       decl.Range{Min: 1, Max: math.MaxInt},
		},
		{
			Import:       "github.com/RussellLuo/validating",
			Qualifier:    "v",
			Name:         "Match",
			IsGeneric:    false,
			Alias:        "match",
			AllowedTypes: []string{"string|bytes"},
			ArgNum:       decl.Range{Min: 1, Max: 1},
		},
		{
			Import:       "github.com/RussellLuo/vext",
			Qualifier:    "vext",
			Name:         "Email",
			IsGeneric:    false,
			Alias:        "email",
			AllowedTypes: []string{"string"},
			ArgNum:       decl.Range{Min: 0, Max: 0},
		},
		{
			Import:       "github.com/RussellLuo/vext",
			Qualifier:    "vext",
			Name:         "IP",
			IsGeneric:    false,
			Alias:        "ip",
			AllowedTypes: []string{"string"},
			ArgNum:       decl.Range{Min: 0, Max: 0},
		},
		{
			Import:       "github.com/RussellLuo/vext",
			Qualifier:    "vext",
			Name:         "Time",
			IsGeneric:    false,
			Alias:        "time",
			AllowedTypes: []string{"string"},
			ArgNum:       decl.Range{Min: 1, Max: 1},
		},
	}

	if !cmp.Equal(got, want) {
		diff := cmp.Diff(got, want)
		t.Errorf("Want - Got: %s", diff)
	}
}
