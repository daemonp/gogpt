// File: pkg/exporter/exporter_test.go
package exporter

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/daemonp/gogpt/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		dir           string
		flags         *types.Flags
		expectedError bool
	}{
		{
			name:          "Valid directory",
			dir:           "/valid/dir",
			flags:         &types.Flags{},
			expectedError: false,
		},
		{
			name:          "Invalid directory",
			dir:           "",
			flags:         &types.Flags{},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp, err := New(tt.dir, tt.flags)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, exp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, exp)
			}
		})
	}
}

func TestExport(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "exporter_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create some test files
	testFiles := []string{
		"file1.txt",
		"file2.go",
		"subdir/file3.md",
	}
	for _, file := range testFiles {
		path := filepath.Join(tempDir, file)
		err := os.MkdirAll(filepath.Dir(path), 0755)
		assert.NoError(t, err)
		err = ioutil.WriteFile(path, []byte("test content"), 0644)
		assert.NoError(t, err)
	}

	tests := []struct {
		name          string
		flags         *types.Flags
		expectedError bool
	}{
		{
			name:          "Default export",
			flags:         &types.Flags{},
			expectedError: false,
		},
		{
			name: "Export with output file",
			flags: &types.Flags{
				OutputFile: filepath.Join(tempDir, "output.txt"),
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp, err := New(tempDir, tt.flags)
			assert.NoError(t, err)

			err = exp.Export()
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				if tt.flags.OutputFile != "" {
					// Check if output file was created
					_, err := os.Stat(tt.flags.OutputFile)
					assert.NoError(t, err)

					// Check if output file contains content
					content, err := ioutil.ReadFile(tt.flags.OutputFile)
					assert.NoError(t, err)
					assert.NotEmpty(t, content)
				}
			}
		})
	}
}
