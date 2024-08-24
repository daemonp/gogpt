// File: cmd/gogpt/main_test.go

package main

import (
	"bytes"
	"errors"
	"testing"

	"github.com/daemonp/gogpt/pkg/exporter"
	"github.com/daemonp/gogpt/pkg/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	// Save original functions and restore them at the end of the test
	originalOsExit := osExit
	originalOsGetwd := osGetwd
	originalExporterNew := exporterNew
	originalParseFlagsFunc := parseFlagsFunc
	defer func() {
		osExit = originalOsExit
		osGetwd = originalOsGetwd
		exporterNew = originalExporterNew
		parseFlagsFunc = originalParseFlagsFunc
	}()

	tests := []struct {
		name           string
		setupMocks     func()
		expectedOutput string
		expectedExit   int
	}{
		{
			name: "Successful export",
			setupMocks: func() {
				osGetwd = func() (string, error) {
					return "/test/dir", nil
				}
				exporterNew = func(rootDir string, flags *types.Flags) (*exporter.Exporter, error) {
					return &exporter.Exporter{}, nil
				}
				parseFlagsFunc = func() *types.Flags {
					return &types.Flags{Verbose: true}
				}
			},
			expectedOutput: "Export completed successfully",
			expectedExit:   0,
		},
		{
			name: "Failed to get current working directory",
			setupMocks: func() {
				osGetwd = func() (string, error) {
					return "", errors.New("failed to get current working directory")
				}
				parseFlagsFunc = func() *types.Flags {
					return &types.Flags{Verbose: true}
				}
			},
			expectedOutput: "Failed to get current working directory",
			expectedExit:   1,
		},
		{
			name: "Failed to create exporter",
			setupMocks: func() {
				osGetwd = func() (string, error) {
					return "/test/dir", nil
				}
				exporterNew = func(rootDir string, flags *types.Flags) (*exporter.Exporter, error) {
					return nil, errors.New("failed to create exporter")
				}
				parseFlagsFunc = func() *types.Flags {
					return &types.Flags{Verbose: true}
				}
			},
			expectedOutput: "Failed to create exporter",
			expectedExit:   1,
		},
		{
			name: "Failed to export repository contents",
			setupMocks: func() {
				osGetwd = func() (string, error) {
					return "/test/dir", nil
				}
				exporterNew = func(rootDir string, flags *types.Flags) (*exporter.Exporter, error) {
					return &exporter.Exporter{}, nil
				}
				parseFlagsFunc = func() *types.Flags {
					return &types.Flags{Verbose: true}
				}
			},
			expectedOutput: "Failed to export repository contents",
			expectedExit:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var buf bytes.Buffer
			log.Logger = zerolog.New(&buf).With().Timestamp().Logger()

			tt.setupMocks()

			// Mock os.Exit
			var exitCode int
			osExit = func(code int) {
				exitCode = code
				panic("os.Exit called")
			}

			// Run main
			func() {
				defer func() {
					if r := recover(); r != nil {
						if r != "os.Exit called" {
							t.Fatalf("Unexpected panic: %v", r)
						}
					}
				}()
				main()
			}()

			// Check exit code
			assert.Equal(t, tt.expectedExit, exitCode)

			// Check log output
			output := buf.String()
			assert.Contains(t, output, tt.expectedOutput)
		})
	}
}
