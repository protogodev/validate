package expr

import (
	"fmt"
	"go/types"
)

type Converter func(name string, typ types.Type) (string, error)

var DefaultConverters = map[string]Converter{
	"_": func(name string, typ types.Type) (string, error) {
		if !IsStruct(typ) {
			return "", ConversionError(name, typ)
		}
		return "_.Schema", nil
	},
	"nonzero": func(name string, typ types.Type) (string, error) {
		if !IsComparable(typ) {
			return "", ConversionError(name, typ)
		}
		return "Nonzero[" + typ.String() + "]", nil
	},
	"zero": func(name string, typ types.Type) (string, error) {
		if !IsComparable(typ) {
			return "", ConversionError(name, typ)
		}
		return "Zero[" + typ.String() + "]", nil
	},
	"len": func(name string, typ types.Type) (string, error) {
		switch {
		case IsString(typ):
			return "LenString", nil
		case IsSlice(typ):
			return "LenSlice", nil
		default:
			return "", ConversionError(name, typ)
		}
	},
	"runecount": func(name string, typ types.Type) (string, error) {
		if !IsString(typ) && !IsBytes(typ) {
			return "", ConversionError(name, typ)
		}
		return "RuneCount", nil
	},
	"eq": func(name string, typ types.Type) (string, error) {
		if !IsComparable(typ) {
			return "", ConversionError(name, typ)
		}
		return "Eq[" + typ.String() + "]", nil
	},
	"ne": func(name string, typ types.Type) (string, error) {
		if !IsComparable(typ) {
			return "", ConversionError(name, typ)
		}
		return "Ne[" + typ.String() + "]", nil
	},
	"gt": func(name string, typ types.Type) (string, error) {
		if !IsOrdered(typ) {
			return "", ConversionError(name, typ)
		}
		return "Gt[" + typ.String() + "]", nil
	},
	"gte": func(name string, typ types.Type) (string, error) {
		if !IsOrdered(typ) {
			return "", ConversionError(name, typ)
		}
		return "Gte[" + typ.String() + "]", nil
	},
	"lt": func(name string, typ types.Type) (string, error) {
		if !IsOrdered(typ) {
			return "", ConversionError(name, typ)
		}
		return "Lt[" + typ.String() + "]", nil
	},
	"lte": func(name string, typ types.Type) (string, error) {
		if !IsOrdered(typ) {
			return "", ConversionError(name, typ)
		}
		return "Lte[" + typ.String() + "]", nil
	},
	"xrange": func(name string, typ types.Type) (string, error) {
		if !IsOrdered(typ) {
			return "", ConversionError(name, typ)
		}
		return "Range[" + typ.String() + "]", nil
	},
	"in": func(name string, typ types.Type) (string, error) {
		if !IsComparable(typ) {
			return "", ConversionError(name, typ)
		}
		return "In[" + typ.String() + "]", nil
	},
	"nin": func(name string, typ types.Type) (string, error) {
		if !IsComparable(typ) {
			return "", ConversionError(name, typ)
		}
		return "Nin[" + typ.String() + "]", nil
	},
	"match": func(name string, typ types.Type) (string, error) {
		if !IsString(typ) && !IsBytes(typ) {
			return "", ConversionError(name, typ)
		}
		return "Match", nil
	},
}

func ConversionError(name string, typ types.Type) error {
	return fmt.Errorf("cannot use validator `%s` on type %T", name, typ.Underlying())
}

func IsStruct(typ types.Type) bool {
	_, ok := typ.Underlying().(*types.Struct)
	return ok
}

func IsComparable(typ types.Type) bool {
	_, ok := typ.Underlying().(*types.Basic)
	return ok
}

func IsOrdered(typ types.Type) bool {
	switch t := typ.Underlying().(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Uintptr,
			types.Float32, types.Float64,
			types.String:
			return true
		}
	}
	return false
}

func IsString(typ types.Type) bool {
	switch t := typ.Underlying().(type) {
	case *types.Basic:
		if t.Kind() == types.String {
			return true
		}
	}
	return false
}

func IsBytes(typ types.Type) bool {
	switch t := typ.Underlying().(type) {
	case *types.Slice:
		switch et := t.Elem().(type) {
		case *types.Basic:
			if et.Kind() == types.Byte {
				return true
			}
		}
	}
	return false
}

func IsSlice(typ types.Type) bool {
	_, ok := typ.Underlying().(*types.Slice)
	return ok
}
