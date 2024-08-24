// File: pkg/exporter/exporter.go

package exporter

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/daemonp/gogpt/pkg/gitignore"
	"github.com/daemonp/gogpt/pkg/types"
	"github.com/rs/zerolog/log"
)

type Exporter struct {
	rootDir       string
	flags         *types.Flags
	fileProcessor *FileProcessor
	contentFilter *ContentFilter
	treeGenerator *TreeGenerator
	writer        *Writer
}

func New(rootDir string, flags *types.Flags) (*Exporter, error) {
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	if _, err := os.Stat(absRootDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", absRootDir)
	}

	var gitIgnore *gitignore.GitIgnore
	if flags.UseGitIgnore {
		gitIgnore, err = gitignore.NewGitIgnore(absRootDir)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to parse .gitignore files, continuing without gitignore")
		}
	}

	contentFilter, err := NewContentFilter(flags.ExcludePattern)
	if err != nil {
		return nil, fmt.Errorf("failed to create content filter: %w", err)
	}

	fileProcessor := NewFileProcessor(absRootDir, flags, gitIgnore)
	treeGenerator := NewTreeGenerator()
	writer := NewWriter(os.Stdout)

	return &Exporter{
		rootDir:       absRootDir,
		flags:         flags,
		fileProcessor: fileProcessor,
		contentFilter: contentFilter,
		treeGenerator: treeGenerator,
		writer:        writer,
	}, nil
}

func (e *Exporter) Export() error {
	files, err := e.fileProcessor.ScanFiles()
	if err != nil {
		return fmt.Errorf("failed to scan files: %w", err)
	}

	treeStructure, err := e.treeGenerator.Generate(files)
	if err != nil {
		return fmt.Errorf("failed to generate tree structure: %w", err)
	}

	if e.flags.OutputFile != "" {
		file, err := os.Create(e.flags.OutputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer file.Close()
		e.writer = NewWriter(file)
	}

	e.writer.Write("# Repository Export\n\n")
	e.writer.Write("This document is a structured representation of the contents of the repository. It includes a list of files and their contents as per the following criteria:\n\n")
	e.writer.Write(fmt.Sprintf("* Files are included based on the specified languages: %s.\n", e.flags.Languages))
	e.writer.Write("* Files ignored by .gitignore are excluded.\n")
	e.writer.Write(fmt.Sprintf("* Files exceeding the token limit (%d tokens) are noted but not included.\n", e.flags.MaxTokens))
	e.writer.Write(fmt.Sprintf("* Lines matching the exclude pattern '%s' are filtered out.\n", e.flags.ExcludePattern))
	e.writer.Write("\n")

	e.writer.Write(treeStructure)

	for _, file := range files {
		if !file.Excluded {
			file.Content = e.contentFilter.Filter(file.Content)
		}
	}

	if err := e.writer.WriteFileContents(files); err != nil {
		return fmt.Errorf("failed to write file contents: %w", err)
	}

	return nil
}
