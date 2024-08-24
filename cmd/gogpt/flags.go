package main

import (
	"flag"

	"github.com/daemonp/gogpt/pkg/exporter"
)

func parseFlags() *exporter.Flags {
	flags := &exporter.Flags{}

	flag.StringVar(&flags.OutputFile, "f", "", "Output file path (default: stdout)")
	flag.BoolVar(&flags.IgnoreGitIgnore, "i", false, "Ignore files specified in .gitignore")
	flag.StringVar(&flags.Languages, "l", "", "Comma-separated list of languages to include (e.g., 'go,js,md')")
	flag.IntVar(&flags.MaxTokens, "max-tokens", 1000, "Maximum number of tokens per file (default: 1000)")
	flag.BoolVar(&flags.Verbose, "v", false, "Enable verbose logging")
	flag.Parse()

	return flags
}
