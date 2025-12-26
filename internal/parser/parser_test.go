package parser

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRegistry_ParseFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParser := NewMockParser(ctrl)
	reg := NewRegistry()
	reg.Register(mockParser, ".mock")

	tmpDir := t.TempDir()
	content := []byte("dummy content")

	tmpFile := filepath.Join(tmpDir, "file.mock")
	require.NoError(t, os.WriteFile(tmpFile, content, 0644))

	tmpFileUnsupported := filepath.Join(tmpDir, "file.mock.txt")
	require.NoError(t, os.WriteFile(tmpFileUnsupported, content, 0644))

	tests := []struct {
		name         string
		filePath     string
		setupMock    func()
		expectResult map[string]interface{}
		expectErr    bool
	}{
		{
			name:     "success",
			filePath: tmpFile,
			setupMock: func() {
				mockParser.EXPECT().Parse(content).Return(map[string]interface{}{"key": "value"}, nil)
			},
			expectResult: map[string]interface{}{"key": "value"},
		},
		{
			name:      "unsupported extension",
			filePath:  tmpFileUnsupported,
			setupMock: func() {},
			expectErr: true,
		},
		{
			name:     "parser error",
			filePath: tmpFile,
			setupMock: func() {
				mockParser.EXPECT().Parse(content).Return(nil, errors.New("parse fail"))
			},
			expectErr: true,
		},
		{
			name:      "file not found",
			filePath:  filepath.Join(tmpDir, "nofile.mock"),
			setupMock: func() {},
			expectErr: true,
		},
		{
			name:      "invalid path for Abs",
			filePath:  string([]byte{0}), // недопустимый путь
			setupMock: func() {},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			res, err := reg.ParseFile(tt.filePath)
			if tt.expectErr {
				require.Error(t, err)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectResult, res)
			}
		})
	}
}
