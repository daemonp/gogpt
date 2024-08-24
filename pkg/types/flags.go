// File: pkg/types/flags.go
package types

type Flags struct {
	OutputFile     string
	UseGitIgnore   bool
	Languages      string
	MaxTokens      int
	Verbose        bool
	ExcludePattern string
}
