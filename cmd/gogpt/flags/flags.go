// File: cmd/gogpt/flags/flags.go
package flags

import (
	"flag"
	"strings"

	"github.com/daemonp/gogpt/pkg/types"
)

func ParseFlags() *types.Flags {
	var excludePaths string
	var maxTokens int

	flags := &types.Flags{}

	flag.StringVar(&flags.OutputFile, "f", "", "Output file path (default: stdout)")
	flag.BoolVar(&flags.UseGitIgnore, "i", true, "Use .gitignore (default: true)")
	flag.StringVar(&flags.Languages, "l", "", "Comma-separated list of languages to include (e.g., 'go,js,md')")
	flag.IntVar(&maxTokens, "max-tokens", 0, "Maximum number of tokens per file (default: no limit)")
	flag.BoolVar(&flags.Verbose, "v", false, "Enable verbose logging")
	flag.StringVar(&flags.ExcludePattern, "exclude", "", "Regex pattern to exclude lines (e.g., '^\\s*//')")
	flag.StringVar(&excludePaths, "x", "", "Comma-separated list of paths to exclude")

	flag.Parse()

	if excludePaths != "" {
		flags.ExcludePaths = strings.Split(excludePaths, ",")
	}

	if maxTokens > 0 {
		flags.MaxTokens = &maxTokens
	}

	return flags
}
