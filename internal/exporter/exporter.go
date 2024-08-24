package exporter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/daemonp/gogpt/internal/fileutils"
	"github.com/daemonp/gogpt/internal/gitignore"
	"github.com/daemonp/gogpt/pkg/tiktoken"
)

type Exporter struct {
	rootDir         string
	outputFile      string
	ignoreGitIgnore bool
	languages       []string
	gitIgnore       *gitignore.GitIgnore
	output          io.Writer
	maxTokens       int
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
	if ignoreGitIgnore {
		var err error
		gitIgnore, err = gitignore.NewGitIgnore(rootDir)
		if err != nil {
			return nil, fmt.Errorf("failed to parse .gitignore: %w", err)
		}
	}

	return &Exporter{
		rootDir:         rootDir,
		outputFile:      outputFile,
		ignoreGitIgnore: ignoreGitIgnore,
		languages:       strings.Split(languages, ","),
		gitIgnore:       gitIgnore,
		output:          output,
		maxTokens:       maxTokens,
	}, nil
}

func (e *Exporter) Export() error {
	return filepath.Walk(e.rootDir, e.processFile)
}

func (e *Exporter) processFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	if e.ignoreGitIgnore && e.gitIgnore.ShouldIgnore(path) {
		return nil
	}

	if !e.shouldIncludeFile(path) {
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	tokenCount := tiktoken.CountTokens(string(content))
	if tokenCount > e.maxTokens {
		log.Warn().Str("file", path).Int("tokens", tokenCount).Msg("File excluded due to size")
		fmt.Fprintf(e.output, "// File: %s (excluded due to size: %d tokens)\n\n", path, tokenCount)
		return nil
	}

	fmt.Fprintf(e.output, "// File: %s\n", path)
	fmt.Fprintf(e.output, "```%s\n", fileutils.GetFileExtension(path))
	fmt.Fprintf(e.output, "%s\n", content)
	fmt.Fprintf(e.output, "```\n\n")

	return nil
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
