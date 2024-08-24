// pkg/exporter/file_processor.go

package exporter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/daemonp/gogpt/pkg/fileutils"
	"github.com/daemonp/gogpt/pkg/gitignore"
	"github.com/daemonp/gogpt/pkg/tiktoken"
	"github.com/rs/zerolog/log"
)

type FileProcessor struct {
	rootDir   string
	languages []string
	maxTokens int
	gitIgnore *gitignore.GitIgnore
}

type FileInfo struct {
	Path       string
	Content    []byte
	TokenCount int
	Excluded   bool
}

func NewFileProcessor(rootDir, languages string, maxTokens int, gitIgnore *gitignore.GitIgnore) *FileProcessor {
	return &FileProcessor{
		rootDir:   rootDir,
		languages: strings.Split(languages, ","),
		maxTokens: maxTokens,
		gitIgnore: gitIgnore,
	}
}

func (fp *FileProcessor) ScanFiles() ([]FileInfo, error) {
	var files []FileInfo
	var wg sync.WaitGroup
	var mu sync.Mutex

	err := filepath.Walk(fp.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || fp.shouldIgnoreFile(path) {
			return nil
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			fileInfo, err := fp.processFile(path)
			if err != nil {
				log.Error().Err(err).Str("file", path).Msg("Failed to process file")
				return
			}

			mu.Lock()
			files = append(files, fileInfo)
			mu.Unlock()
		}()

		return nil
	})

	wg.Wait()

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return fp.includeSpecialFiles(files), nil
}

func (fp *FileProcessor) processFile(path string) (FileInfo, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return FileInfo{}, fmt.Errorf("failed to read file: %w", err)
	}

	tokenCount := tiktoken.CountTokens(string(content))
	excluded := tokenCount > fp.maxTokens

	if excluded {
		log.Warn().Str("file", path).Int("tokens", tokenCount).Msg("File excluded due to size")
		content = []byte(fmt.Sprintf("// File excluded due to size: %d tokens", tokenCount))
	}

	return FileInfo{
		Path:       path,
		Content:    content,
		TokenCount: tokenCount,
		Excluded:   excluded,
	}, nil
}

func (fp *FileProcessor) shouldIgnoreFile(path string) bool {
	if fp.gitIgnore != nil && fp.gitIgnore.ShouldIgnore(path) {
		return true
	}

	ext := fileutils.GetFileExtension(path)
	for _, lang := range fp.languages {
		if fileutils.IsLanguageFile(lang, ext) {
			return false
		}
	}

	return true
}

func (fp *FileProcessor) includeSpecialFiles(files []FileInfo) []FileInfo {
	specialFiles := []string{
		filepath.Join(fp.rootDir, "README.md"),
		filepath.Join(fp.rootDir, ".gitignore"),
	}

	for _, specialFile := range specialFiles {
		if _, err := os.Stat(specialFile); err == nil {
			fileInfo, err := fp.processFile(specialFile)
			if err != nil {
				log.Error().Err(err).Str("file", specialFile).Msg("Failed to process special file")
				continue
			}
			files = append(files, fileInfo)
		}
	}

	return files
}
