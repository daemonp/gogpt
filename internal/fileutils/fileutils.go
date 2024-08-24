package fileutils

import (
	"path/filepath"
	"strings"
)

var LanguageExtensions = map[string][]string{
	"go":         {".go"},
	"js":         {".js", ".jsx"},
	"ts":         {".ts", ".tsx"},
	"ruby":       {".rb", ".erb"},
	"python":     {".py"},
	"java":       {".java"},
	"c":          {".c", ".h"},
	"cpp":        {".cpp", ".hpp", ".cc", ".hh"},
	"csharp":     {".cs"},
	"php":        {".php"},
	"swift":      {".swift"},
	"rust":       {".rs"},
	"kotlin":     {".kt", ".kts"},
	"scala":      {".scala"},
	"html":       {".html", ".htm"},
	"css":        {".css"},
	"markdown":   {".md", ".markdown"},
	"yaml":       {".yaml", ".yml"},
	"json":       {".json"},
	"xml":        {".xml"},
	"sql":        {".sql"},
	"shell":      {".sh", ".bash"},
	"powershell": {".ps1"},
	"docker":     {"Dockerfile"},
	"make":       {"Makefile"},
	"config":     {".cfg", ".conf", ".ini"},
}

func GetFileExtension(path string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return filepath.Base(path)
	}
	return strings.TrimPrefix(ext, ".") // Remove the leading dot
}

func IsLanguageFile(lang, ext string) bool {
	extensions, ok := LanguageExtensions[strings.ToLower(lang)]
	if !ok {
		return false
	}
	for _, e := range extensions {
		if strings.EqualFold(e, "."+ext) || strings.EqualFold(e, ext) {
			return true
		}
	}
	return false
}
