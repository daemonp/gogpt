// File: cmd/gogpt/flags/flags.go
package flags

import (
	"flag"

	"github.com/daemonp/gogpt/pkg/types"
)

func ParseFlags() *types.Flags {
	flags := &types.Flags{}

	flag.StringVar(&flags.OutputFile, "f", "", "Output file path (default: stdout)")
	flag.BoolVar(&flags.UseGitIgnore, "i", true, "Use .gitignore (default: true)")
	flag.StringVar(&flags.Languages, "l", "", "Comma-separated list of languages to include (e.g., 'go,js,md')")
	flag.IntVar(&flags.MaxTokens, "max-tokens", 1000, "Maximum number of tokens per file (default: 1000)")
	flag.BoolVar(&flags.Verbose, "v", false, "Enable verbose logging")
	flag.StringVar(&flags.ExcludePattern, "exclude", "", "Regex pattern to exclude lines (e.g., '^\\s*//')")
	flag.Parse()

	return flags
}
