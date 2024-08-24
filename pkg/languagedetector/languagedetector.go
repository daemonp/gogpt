package languagedetector

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/daemonp/gogpt/pkg/fileutils"
)

func DetectLanguages(dir string) string {
	languages := make(map[string]bool)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := fileutils.GetFileExtension(path)
		for lang, extensions := range fileutils.LanguageExtensions {
			for _, e := range extensions {
				if strings.EqualFold(e, "."+ext) || strings.EqualFold(e, ext) {
					languages[lang] = true
					break
				}
			}
		}

		return nil
	})

	var detectedLangs []string
	for lang := range languages {
		detectedLangs = append(detectedLangs, lang)
	}

	return strings.Join(detectedLangs, ",")
}
