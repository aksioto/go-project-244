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
			file1: filepath.Join("testdata", "fixture", "file1.json"),
			file2: filepath.Join("testdata", "fixture", "file2.json"),
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
			file1: filepath.Join("testdata", "fixture", "same1.json"),
			file2: filepath.Join("testdata", "fixture", "same2.json"),
			expected: `{
    host: hexlet.io
    timeout: 50
}`,
		},
		{
			name:  "Empty files",
			file1: filepath.Join("testdata", "fixture", "empty1.json"),
			file2: filepath.Join("testdata", "fixture", "empty2.json"),
			expected: `{
}`,
		},
		{
			name:  "Empty vs filled",
			file1: filepath.Join("testdata", "fixture", "empty_vs_filled1.json"),
			file2: filepath.Join("testdata", "fixture", "empty_vs_filled2.json"),
			expected: `{
  + host: hexlet.io
  + timeout: 50
}`,
		},
		{
			name:  "Only deleted",
			file1: filepath.Join("testdata", "fixture", "only_deleted.json"),
			file2: filepath.Join("testdata", "fixture", "only_added.json"),
			expected: `{
    host: hexlet.io
  - proxy: 123.234.53.22
  - timeout: 50
}`,
		},
		{
			name:  "Only added",
			file1: filepath.Join("testdata", "fixture", "only_added.json"),
			file2: filepath.Join("testdata", "fixture", "only_deleted.json"),
			expected: `{
    host: hexlet.io
  + proxy: 123.234.53.22
  + timeout: 50
}`,
		},
		{
			name:  "Different types",
			file1: filepath.Join("testdata", "fixture", "different_types.json"),
			file2: filepath.Join("testdata", "fixture", "different_types2.json"),
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
			file1: filepath.Join("testdata", "fixture", "completely_different1.json"),
			file2: filepath.Join("testdata", "fixture", "completely_different2.json"),
			expected: `{
  - key1: value1
  - key2: 123
  + key3: value3
  + key4: 456
}`,
		},
		{
			name:      "Nonexistent file",
			file1:     filepath.Join("testdata", "fixture", "nonexistent.json"),
			file2:     filepath.Join("testdata", "fixture", "file1.json"),
			expectErr: true,
		},
		{
			name:      "Both nonexistent files",
			file1:     filepath.Join("testdata", "fixture", "nonexistent1.json"),
			file2:     filepath.Join("testdata", "fixture", "nonexistent2.json"),
			expectErr: true,
		},
		{
			name:      "Unsupported format",
			file1:     filepath.Join("testdata", "fixture", "file1.txt"),
			file2:     filepath.Join("testdata", "fixture", "file2.json"),
			expectErr: true,
		},
		// YAML tests
		{
			name:  "YAML flat diff",
			file1: filepath.Join("testdata", "fixture", "file1.yml"),
			file2: filepath.Join("testdata", "fixture", "file2.yml"),
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
			name:  "YAML same files",
			file1: filepath.Join("testdata", "fixture", "same1.yml"),
			file2: filepath.Join("testdata", "fixture", "same2.yml"),
			expected: `{
    host: hexlet.io
    timeout: 50
}`,
		},
		{
			name:  "YAML empty files",
			file1: filepath.Join("testdata", "fixture", "empty1.yml"),
			file2: filepath.Join("testdata", "fixture", "empty2.yml"),
			expected: `{
}`,
		},
		{
			name:  "YAML empty vs filled",
			file1: filepath.Join("testdata", "fixture", "empty_vs_filled1.yml"),
			file2: filepath.Join("testdata", "fixture", "empty_vs_filled2.yml"),
			expected: `{
  + host: hexlet.io
  + timeout: 50
}`,
		},
		{
			name:  "YAML only deleted",
			file1: filepath.Join("testdata", "fixture", "only_deleted.yml"),
			file2: filepath.Join("testdata", "fixture", "only_added.yml"),
			expected: `{
    host: hexlet.io
  - proxy: 123.234.53.22
  - timeout: 50
}`,
		},
		{
			name:  "YAML only added",
			file1: filepath.Join("testdata", "fixture", "only_added.yml"),
			file2: filepath.Join("testdata", "fixture", "only_deleted.yml"),
			expected: `{
    host: hexlet.io
  + proxy: 123.234.53.22
  + timeout: 50
}`,
		},
		{
			name:  "YAML different types",
			file1: filepath.Join("testdata", "fixture", "different_types.yml"),
			file2: filepath.Join("testdata", "fixture", "different_types2.yml"),
			expected: `{
  - boolean: true
  + boolean: false
  - float: 3.14
  + float: 2.71
  - number: 42
  + number: 100
  - string: value
  + string: different
}`,
		},
		{
			name:  "YAML completely different",
			file1: filepath.Join("testdata", "fixture", "completely_different1.yml"),
			file2: filepath.Join("testdata", "fixture", "completely_different2.yml"),
			expected: `{
  - key1: value1
  - key2: 123
  + key3: value3
  + key4: 456
}`,
		},
		{
			name:  "YAML and JSON mixed",
			file1: filepath.Join("testdata", "fixture", "file1.yml"),
			file2: filepath.Join("testdata", "fixture", "file2.json"),
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
			name:  "JSON and YAML mixed",
			file1: filepath.Join("testdata", "fixture", "file1.json"),
			file2: filepath.Join("testdata", "fixture", "file2.yml"),
			expected: `{
  - follow: false
    host: hexlet.io
  - proxy: 123.234.53.22
  - timeout: 50
  + timeout: 20
  + verbose: true
}`,
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
