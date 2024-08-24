// File: pkg/exporter/file_processor.go

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
	rootDir        string
	languages      []string
	maxTokens      int
	gitIgnore      *gitignore.GitIgnore
	useGitIgnore   bool
	customScanFunc func() ([]FileInfo, error)
}

type FileInfo struct {
	Path       string
	Content    []byte
	TokenCount int
	Excluded   bool
}

func NewFileProcessor(rootDir, languages string, maxTokens int, gitIgnore *gitignore.GitIgnore, useGitIgnore bool) *FileProcessor {
	return &FileProcessor{
		rootDir:      rootDir,
		languages:    strings.Split(languages, ","),
		maxTokens:    maxTokens,
		gitIgnore:    gitIgnore,
		useGitIgnore: useGitIgnore,
	}
}

func (fp *FileProcessor) SetCustomScanFunc(scanFunc func() ([]FileInfo, error)) {
	fp.customScanFunc = scanFunc
}

func (fp *FileProcessor) ScanFiles() ([]FileInfo, error) {
	if fp.customScanFunc != nil {
		return fp.customScanFunc()
	}

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

			relPath, err := filepath.Rel(fp.rootDir, path)
			if err != nil {
				log.Error().Err(err).Str("file", path).Msg("Failed to get relative path")
				return
			}

			fileInfo, err := fp.processFile(relPath)
			if err != nil {
				log.Error().Err(err).Str("file", relPath).Msg("Failed to process file")
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
	content, err := os.ReadFile(filepath.Join(fp.rootDir, path))
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
	if fp.useGitIgnore && fp.gitIgnore != nil && fp.gitIgnore.ShouldIgnore(path) {
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
		".gitignore",
	}

	// Only include README.md if markdown is not in the selected languages
	if !fp.isLanguageIncluded("markdown") {
		specialFiles = append(specialFiles, "README.md")
	}

	for _, specialFile := range specialFiles {
		filePath := filepath.Join(fp.rootDir, specialFile)
		if _, err := os.Stat(filePath); err == nil {
			relPath, err := filepath.Rel(fp.rootDir, filePath)
			if err != nil {
				log.Error().Err(err).Str("file", filePath).Msg("Failed to get relative path for special file")
				continue
			}

			fileInfo, err := fp.processFile(relPath)
			if err != nil {
				log.Error().Err(err).Str("file", relPath).Msg("Failed to process special file")
				continue
			}
			files = append(files, fileInfo)
		}
	}

	return files
}

func (fp *FileProcessor) isLanguageIncluded(lang string) bool {
	for _, l := range fp.languages {
		if strings.TrimSpace(strings.ToLower(l)) == lang {
			return true
		}
	}
	return false
}
