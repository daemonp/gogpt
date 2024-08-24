package exporter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/daemonp/gogpt/internal/fileutils"
	"github.com/daemonp/gogpt/internal/gitignore"
	"github.com/daemonp/gogpt/pkg/tiktoken"
	"github.com/rs/zerolog/log"
)

type Exporter struct {
	rootDir         string
	outputFile      string
	ignoreGitIgnore bool
	languages       []string
	gitIgnore       *gitignore.GitIgnore
	output          io.Writer
	maxTokens       int
	includedFiles   []string
}

func New(rootDir, outputFile string, ignoreGitIgnore bool, languages string, maxTokens int) (*Exporter, error) {
	var output io.Writer = os.Stdout
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create output file: %w", err)
		}
		output = file
	}

	var gitIgnore *gitignore.GitIgnore
	var err error
	gitIgnore, err = gitignore.NewGitIgnore(rootDir)
	if err != nil {
		return nil, fmt.Errorf("failed to parse .gitignore: %w", err)
	}

	return &Exporter{
		rootDir:         rootDir,
		outputFile:      outputFile,
		ignoreGitIgnore: ignoreGitIgnore,
		languages:       strings.Split(languages, ","),
		gitIgnore:       gitIgnore,
		output:          output,
		maxTokens:       maxTokens,
		includedFiles:   []string{},
	}, nil
}

func (e *Exporter) scanFiles() error {
	var wg sync.WaitGroup
	var mu sync.Mutex

	err := filepath.Walk(e.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if e.shouldIgnoreFile(path) {
				return filepath.SkipDir
			}
			return nil
		}

		if e.shouldIgnoreFile(path) {
			return nil
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			content, err := os.ReadFile(path)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to read file: %s", path)
				return
			}

			tokenCount := tiktoken.CountTokens(string(content))
			if tokenCount > e.maxTokens {
				log.Warn().Str("file", path).Int("tokens", tokenCount).Msg("File excluded due to size")
				mu.Lock()
				fmt.Fprintf(e.output, "// File: %s (excluded due to size: %d tokens)\n\n", path, tokenCount)
				mu.Unlock()
				return
			}

			mu.Lock()
			e.includedFiles = append(e.includedFiles, path)
			mu.Unlock()
		}()

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	wg.Wait()

	// Ensure README.md and .gitignore are included if they exist
	e.includeSpecialFiles()

	return nil
}

func (e *Exporter) shouldIgnoreFile(path string) bool {
	relPath, err := filepath.Rel(e.rootDir, path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to get relative path")
		return false
	}

	if e.gitIgnore.ShouldIgnore(relPath) {
		return true
	}

	return !e.shouldIncludeFile(path)
}

func (e *Exporter) shouldIncludeFile(path string) bool {
	ext := fileutils.GetFileExtension(path)
	for _, lang := range e.languages {
		if fileutils.IsLanguageFile(lang, ext) {
			return true
		}
	}
	return false
}

func (e *Exporter) writeFileContents(path string, content []byte) {
	fmt.Fprintf(e.output, "// File: %s\n", path)
	fmt.Fprintf(e.output, "```%s\n", fileutils.GetFileExtension(path))
	fmt.Fprintf(e.output, "%s\n", content)
	fmt.Fprintf(e.output, "```\n\n")
}
