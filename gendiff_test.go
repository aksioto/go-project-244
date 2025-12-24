package code

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenDiff(t *testing.T) {
	tests := []struct {
		name      string
		file1     string
		file2     string
		expected  string
		expectErr bool
	}{
		{
			name:  "Flat diff",
			file1: "testdata/file1.json",
			file2: "testdata/file2.json",
			expected: `{
  - follow: false
    host: hexlet.io
  - proxy: 123.234.53.22
  - timeout: 50
  + timeout: 20
  + verbose: true
}`,
		},
		{
			name:  "Same files",
			file1: "testdata/same1.json",
			file2: "testdata/same2.json",
			expected: `{
    host: hexlet.io
    timeout: 50
}`,
		},
		{
			name:  "Empty files",
			file1: "testdata/empty1.json",
			file2: "testdata/empty2.json",
			expected: `{
}`,
		},
		{
			name:  "Empty vs filled",
			file1: "testdata/empty_vs_filled1.json",
			file2: "testdata/empty_vs_filled2.json",
			expected: `{
  + host: hexlet.io
  + timeout: 50
}`,
		},
		{
			name:  "Only deleted",
			file1: "testdata/only_deleted.json",
			file2: "testdata/only_added.json",
			expected: `{
    host: hexlet.io
  - proxy: 123.234.53.22
  - timeout: 50
}`,
		},
		{
			name:  "Only added",
			file1: "testdata/only_added.json",
			file2: "testdata/only_deleted.json",
			expected: `{
    host: hexlet.io
  + proxy: 123.234.53.22
  + timeout: 50
}`,
		},
		{
			name:  "Different types",
			file1: "testdata/different_types.json",
			file2: "testdata/different_types2.json",
			expected: `{
  - boolean: true
  + boolean: false
  - float: 3.14
  + float: 2.71
    null: <nil>
  - number: 42
  + number: 100
  - string: value
  + string: different
}`,
		},
		{
			name:  "Completely different",
			file1: "testdata/completely_different1.json",
			file2: "testdata/completely_different2.json",
			expected: `{
  - key1: value1
  - key2: 123
  + key3: value3
  + key4: 456
}`,
		},
		{
			name:      "Nonexistent file",
			file1:     "testdata/nonexistent.json",
			file2:     "testdata/file1.json",
			expectErr: true,
		},
		{
			name:      "Both nonexistent files",
			file1:     "testdata/nonexistent1.json",
			file2:     "testdata/nonexistent2.json",
			expectErr: true,
		},
		{
			name:      "Unsupported format",
			file1:     "testdata/file1.txt",
			file2:     "testdata/file2.json",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			file1 := filepath.Join(tt.file1)
			file2 := filepath.Join(tt.file2)

			result, err := GenDiff(file1, file2)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}
