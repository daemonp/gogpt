package exporter

import (
	"fmt"
	"io"
	"os"

	"github.com/daemonp/gogpt/pkg/gitignore"
	"github.com/daemonp/gogpt/pkg/languagedetector"
	"github.com/rs/zerolog/log"
)

type Exporter struct {
	rootDir       string
	flags         *Flags
	gitIgnore     *gitignore.GitIgnore
	output        io.Writer
	processor     *FileProcessor
	generator     *TreeGenerator
	writer        *Writer
	contentFilter *ContentFilter
}

type Flags struct {
	OutputFile     string
	UseGitIgnore   bool
	Languages      string
	MaxTokens      int
	Verbose        bool
	ExcludePattern string
}

func New(rootDir string, flags *Flags) (*Exporter, error) {
	output := os.Stdout
	if flags.OutputFile != "" {
		file, err := os.Create(flags.OutputFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create output file: %w", err)
		}
		output = file
	}

	var gitIgnore *gitignore.GitIgnore
	if flags.UseGitIgnore {
		var err error
		gitIgnore, err = gitignore.NewGitIgnore(rootDir)
		if err != nil {
			return nil, fmt.Errorf("failed to parse .gitignore: %w", err)
		}
	}

	// If no languages are specified, detect them automatically
	if flags.Languages == "" {
		detectedLangs := languagedetector.DetectLanguages(rootDir)
		flags.Languages = detectedLangs
		log.Info().Str("languages", flags.Languages).Msg("Detected languages")
	}

	contentFilter, err := NewContentFilter(flags.ExcludePattern)
	if err != nil {
		return nil, fmt.Errorf("failed to create content filter: %w", err)
	}

	processor := NewFileProcessor(rootDir, flags.Languages, flags.MaxTokens, gitIgnore, flags.UseGitIgnore)
	generator := NewTreeGenerator()
	writer := NewWriter(output)

	return &Exporter{
		rootDir:       rootDir,
		flags:         flags,
		gitIgnore:     gitIgnore,
		output:        output,
		processor:     processor,
		generator:     generator,
		writer:        writer,
		contentFilter: contentFilter,
	}, nil
}

func (e *Exporter) Export() error {
	files, err := e.processor.ScanFiles()
	if err != nil {
		return fmt.Errorf("failed to scan files: %w", err)
	}

	preamble := e.generatePreamble()
	e.writer.Write(preamble)

	treeOutput, err := e.generator.Generate(files)
	if err != nil {
		return fmt.Errorf("failed to generate tree structure: %w", err)
	}
	e.writer.Write(treeOutput)

	var totalSize int64
	var totalTokens int

	for i, file := range files {
		file.Content = e.contentFilter.Filter(file.Content)
		files[i] = file

		fileSize := int64(len(file.Content))
		totalSize += fileSize
		totalTokens += file.TokenCount

		if e.flags.Verbose {
			logFileInfo(file.Path, fileSize, file.TokenCount)
		}
	}

	if err := e.writer.WriteFileContents(files); err != nil {
		log.Error().Err(err).Msg("Failed to write file contents")
	}

	// Log summary
	log.Info().
		Float64("total_size_kb", float64(totalSize)/1024.0).
		Int("total_tokens", totalTokens).
		Msg("Export completed")

	return nil
}

func (e *Exporter) generatePreamble() string {
	gitIgnoreStatus := "included"
	if e.flags.UseGitIgnore {
		gitIgnoreStatus = "excluded"
	}

	return fmt.Sprintf(`# Repository Export

This document is a structured representation of the contents of the repository. It includes a list of files and their contents as per the following criteria:

* Files are included based on the specified languages: %s.
* Files ignored by .gitignore are %s.
* Files exceeding the token limit (%d tokens) are noted but not included.
* Lines matching the exclude pattern '%s' are filtered out.

`, e.flags.Languages,
		gitIgnoreStatus,
		e.flags.MaxTokens,
		e.flags.ExcludePattern)
}

func logFileInfo(path string, sizeInBytes int64, tokenCount int) {
	sizeInKB := float64(sizeInBytes) / 1024.0
	log.Debug().
		Str("file", path).
		Float64("size_kb", sizeInKB).
		Int("tokens", tokenCount).
		Msg("File processed")
}
