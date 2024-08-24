// File: pkg/exporter/exporter_test.go

package exporter

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/daemonp/gogpt/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func intPtr(i int) *int {
	return &i
}

func TestNew(t *testing.T) {
	// Create a temporary directory for valid tests
	validDir, err := ioutil.TempDir("", "exporter_test_valid")
	require.NoError(t, err)
	defer os.RemoveAll(validDir)

	tests := []struct {
		name          string
		dir           string
		flags         *types.Flags
		expectedError bool
	}{
		{
			name:          "Valid directory",
			dir:           validDir,
			flags:         &types.Flags{},
			expectedError: false,
		},
		{
			name:          "Invalid directory",
			dir:           "/nonexistent/directory",
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
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create some test files
	testFiles := []struct {
		path    string
		content string
	}{
		{"file1.txt", "This is a text file."},
		{"file2.go", "package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}"},
		{"subdir/file3.md", "# Markdown File\n\nThis is a markdown file."},
	}
	for _, file := range testFiles {
		path := filepath.Join(tempDir, file.path)
		err := os.MkdirAll(filepath.Dir(path), 0755)
		require.NoError(t, err)
		err = ioutil.WriteFile(path, []byte(file.content), 0644)
		require.NoError(t, err)
	}

	tests := []struct {
		name          string
		flags         *types.Flags
		expectedError bool
		checkOutput   func(t *testing.T, output string)
	}{
		{
			name: "Default export",
			flags: &types.Flags{
				Languages: "go,markdown",
				MaxTokens: intPtr(1000),
			},
			expectedError: false,
			checkOutput: func(t *testing.T, output string) {
				assert.Contains(t, output, "file2.go")
				assert.Contains(t, output, "subdir/file3.md")
				assert.NotContains(t, output, "file1.txt")
			},
		},
		{
			name: "Export with small token limit",
			flags: &types.Flags{
				Languages: "go,markdown",
				MaxTokens: intPtr(1),
			},
			expectedError: false,
			checkOutput: func(t *testing.T, output string) {
				assert.Contains(t, output, "File excluded due to size")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := filepath.Join(tempDir, "output.txt")
			tt.flags.OutputFile = output

			exp, err := New(tempDir, tt.flags)
			require.NoError(t, err)

			err = exp.Export()
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Check if output file was created
				content, err := ioutil.ReadFile(output)
				assert.NoError(t, err)
				assert.NotEmpty(t, content)

				// Run custom output checks
				tt.checkOutput(t, string(content))
			}
		})
	}
}
