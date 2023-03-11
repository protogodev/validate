package validate_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/protogodev/validate"
)

func TestParseDoc(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want map[string][]validate.Option
	}{
		{
			name: "one",
			in: []string{
				"// @header1:",
				"//   key1: value1",
				"//   key2: value2",
			},
			want: map[string][]validate.Option{
				"header1": {
					{K: "key1", V: "value1"},
					{K: "key2", V: "value2"},
				},
			},
		},
		{
			name: "more",
			in: []string{
				"// @header1:",
				"//   key1: value1",
				"//   key2: value2",
				"//",
				"// @header2:",
				"//   key3: value3",
				"//   key4: value4",
			},
			want: map[string][]validate.Option{
				"header1": {
					{K: "key1", V: "value1"},
					{K: "key2", V: "value2"},
				},
				"header2": {
					{K: "key3", V: "value3"},
					{K: "key4", V: "value4"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validate.ParseDoc(tt.in)
			if !cmp.Equal(got, tt.want) {
				diff := cmp.Diff(got, tt.want)
				t.Errorf("Want - Got: %s", diff)
			}
		})
	}
}
