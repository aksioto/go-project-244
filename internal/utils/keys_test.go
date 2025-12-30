package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type sortedKeysTestCase struct {
	name     string
	input    map[string]interface{}
	expected []string
}

type mergedKeysTestCase struct {
	name     string
	m1       map[string]interface{}
	m2       map[string]interface{}
	expected []string
}

func TestSortedKeys(t *testing.T) {
	testCases := []sortedKeysTestCase{
		{
			name:     "multiple keys",
			input:    map[string]interface{}{"z": 1, "a": 2, "m": 3},
			expected: []string{"a", "m", "z"},
		},
		{
			name:     "empty map",
			input:    map[string]interface{}{},
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, SortedKeys(tc.input))
		})
	}
}

func TestMergedSortedKeys(t *testing.T) {
	testCases := []mergedKeysTestCase{
		{
			name: "with overlaps",
			m1:   map[string]interface{}{"foo": 1, "bar": 2},
			m2:   map[string]interface{}{"baz": 3, "foo": 4},
			expected: []string{
				"bar", "baz", "foo",
			},
		},
		{
			name:     "both empty",
			m1:       map[string]interface{}{},
			m2:       map[string]interface{}{},
			expected: []string{},
		},
		{
			name:     "first empty",
			m1:       map[string]interface{}{},
			m2:       map[string]interface{}{"a": 1},
			expected: []string{"a"},
		},
		{
			name:     "second empty",
			m1:       map[string]interface{}{"b": 2},
			m2:       map[string]interface{}{},
			expected: []string{"b"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, MergedSortedKeys(tc.m1, tc.m2))
		})
	}
}
