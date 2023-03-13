package decl

import (
	"go/types"
)

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
