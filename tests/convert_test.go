package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	flam "github.com/happyhippyhippo/flam"
	config "github.com/happyhippyhippo/flam-config"
)

func Test_Convert(t *testing.T) {
	scenarios := []struct {
		name string
		val  any
		want any
	}{
		{
			name: "flam.Bag",
			val:  flam.Bag{"KEY": "value"},
			want: flam.Bag{"key": "value"},
		},
		{
			name: "slice of any",
			val:  []any{flam.Bag{"KEY": "value"}},
			want: []any{flam.Bag{"key": "value"}},
		},
		{
			name: "map[string]any",
			val:  map[string]any{"KEY": "value"},
			want: flam.Bag{"key": "value"},
		},
		{
			name: "map[any]any with string key",
			val:  map[any]any{"KEY": "value"},
			want: flam.Bag{"key": "value"},
		},
		{
			name: "map[any]any with non-string key",
			val:  map[any]any{123: "value"},
			want: flam.Bag{"123": "value"},
		},
		{
			name: "float64 convertible to int",
			val:  123.0,
			want: 123,
		},
		{
			name: "float64 not convertible to int",
			val:  123.4,
			want: 123.4,
		},
		{
			name: "other primitive types",
			val:  "a string",
			want: "a string",
		},
		{
			name: "nested structure",
			val:  flam.Bag{"L1": map[any]any{"L2": []any{1.0, "VALUE"}}},
			want: flam.Bag{"l1": flam.Bag{"l2": []any{1, "VALUE"}}},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			assert.Equal(t, scenario.want, config.Convert(scenario.val))
		})
	}
}
