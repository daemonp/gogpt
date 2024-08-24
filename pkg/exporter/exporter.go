package exporter

import (
	"fmt"
	"io"
	"os"

	"github.com/daemonp/gogpt/pkg/gitignore"
	"github.com/rs/zerolog/log"
)

type Exporter struct {
	rootDir   string
	flags     *Flags
	gitIgnore *gitignore.GitIgnore
	output    io.Writer
	processor *FileProcessor
	generator *TreeGenerator
	writer    *Writer
}

type Flags struct {
	OutputFile      string
	IgnoreGitIgnore bool
	Languages       string
	MaxTokens       int
	Verbose         bool
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
	if flags.IgnoreGitIgnore {
		var err error
		gitIgnore, err = gitignore.NewGitIgnore(rootDir)
		if err != nil {
			return nil, fmt.Errorf("failed to parse .gitignore: %w", err)
		}
	}

	processor := NewFileProcessor(rootDir, flags.Languages, flags.MaxTokens, gitIgnore)
	generator := NewTreeGenerator()
	writer := NewWriter(output)

	return &Exporter{
		rootDir:   rootDir,
		flags:     flags,
		gitIgnore: gitIgnore,
		output:    output,
		processor: processor,
		generator: generator,
		writer:    writer,
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

	if err := e.writer.WriteFileContents(files); err != nil {
		log.Error().Err(err).Msg("Failed to write file contents")
	}

	return nil
}

func (e *Exporter) generatePreamble() string {
	return fmt.Sprintf(`# Repository Export

This document is a structured representation of the contents of the repository. It includes a list of files and their contents as per the following criteria:

* Files are included based on the specified languages: %s.
* Files ignored by .gitignore are excluded.
* Files exceeding the token limit (%d tokens) are noted but not included.

`, e.flags.Languages, e.flags.MaxTokens)
}
