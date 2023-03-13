package decl

import (
	v "github.com/RussellLuo/validating/v3"
	"github.com/RussellLuo/vext"
)

var _ = []any{
	// type=comparable args=0
	v.Nonzero[string],

	// type=comparable args=0
	v.Zero[string],

	// name=len type=string args=2
	v.LenString,

	// name=len type=slice args=2
	v.LenSlice[[]string],

	// name=runecnt type=string|bytes args=2
	v.RuneCount,

	// type=comparable args=1
	v.Eq[string],

	// type=comparable args=1
	v.Ne[string],

	// type=ordered args=1
	v.Gt[string],

	// type=ordered args=1
	v.Gte[string],

	// type=ordered args=1
	v.Lt[string],

	// type=ordered args=1
	v.Lte[string],

	// name=xrange type=ordered args=2
	v.Range[string],

	// type=ordered args=1+
	v.In[string],

	// type=ordered args=1+
	v.Nin[string],

	// type=string|bytes args=1
	v.Match,

	// type=string args=0
	vext.Email,

	// type=string args=0
	vext.IP,

	// type=string args=1
	vext.Time,
}
