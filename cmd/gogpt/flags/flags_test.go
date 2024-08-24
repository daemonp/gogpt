// File: cmd/gogpt/flags/flags_test.go

package flags

import (
	"flag"
	"os"
	"testing"

	"github.com/daemonp/gogpt/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedFlags *types.Flags
	}{
		{
			name: "Default flags",
			args: []string{"cmd"},
			expectedFlags: &types.Flags{
				Verbose:      false,
				OutputFile:   "",
				UseGitIgnore: true,
				MaxTokens:    1000,
			},
		},
		{
			name: "Verbose flag",
			args: []string{"cmd", "-v"},
			expectedFlags: &types.Flags{
				Verbose:      true,
				OutputFile:   "",
				UseGitIgnore: true,
				MaxTokens:    1000,
			},
		},
		{
			name: "Output file flag",
			args: []string{"cmd", "-f", "test.txt"},
			expectedFlags: &types.Flags{
				Verbose:      false,
				OutputFile:   "test.txt",
				UseGitIgnore: true,
				MaxTokens:    1000,
			},
		},
		{
			name: "All flags",
			args: []string{"cmd", "-v", "-f", "test.txt", "-i=false", "-l=go,js", "--max-tokens=500"},
			expectedFlags: &types.Flags{
				Verbose:      true,
				OutputFile:   "test.txt",
				UseGitIgnore: false,
				Languages:    "go,js",
				MaxTokens:    500,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags before each test case
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Save original os.Args and restore it at the end of the test
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			os.Args = tt.args

			flags := ParseFlags()
			assert.Equal(t, tt.expectedFlags, flags)
		})
	}
}
