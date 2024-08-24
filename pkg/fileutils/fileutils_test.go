// File: pkg/fileutils/fileutils_test.go

package fileutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"Simple extension", "file.txt", "txt"},
		{"No extension", "file", "file"},
		{"Multiple dots", "file.tar.gz", "gz"},
		{"Hidden file", ".gitignore", "gitignore"},
		{"Path with directory", "/path/to/file.go", "go"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFileExtension(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsLanguageFile(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		ext      string
		expected bool
	}{
		{"Go file", "go", "go", true},
		{"JavaScript file", "js", "js", true},
		{"Non-existent language", "nonexistent", "txt", false},
		{"Wrong extension", "go", "js", false},
		{"Case insensitive", "Go", "GO", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsLanguageFile(tt.lang, tt.ext)
			assert.Equal(t, tt.expected, result)
		})
	}
}
