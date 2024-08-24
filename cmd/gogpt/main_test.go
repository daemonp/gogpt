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
	"github.com/stretchr/testify/mock"
)

// MockExporter is a mock implementation of the Exporter
type MockExporter struct {
	mock.Mock
}

func (m *MockExporter) Export() error {
	args := m.Called()
	return args.Error(0)
}

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
		setupMocks     func(*testing.T)
		expectedOutput string
		expectedExit   int
	}{
		{
			name: "Successful export",
			setupMocks: func(t *testing.T) {
				mockExporter := new(MockExporter)
				mockExporter.On("Export").Return(nil)

				osGetwd = func() (string, error) {
					return "/test/dir", nil
				}

				exporterNew = func(rootDir string, flags *types.Flags) (*exporter.Exporter, error) {
					return &exporter.Exporter{}, nil
				}

				parseFlagsFunc = func() *types.Flags {
					return &types.Flags{Verbose: false}
				}
			},
			expectedOutput: "",
			expectedExit:   0,
		},
		{
			name: "Failed to get current working directory",
			setupMocks: func(t *testing.T) {
				osGetwd = func() (string, error) {
					return "", errors.New("failed to get current working directory")
				}

				parseFlagsFunc = func() *types.Flags {
					return &types.Flags{Verbose: false}
				}
			},
			expectedOutput: "Failed to get current working directory",
			expectedExit:   1,
		},
		{
			name: "Failed to create exporter",
			setupMocks: func(t *testing.T) {
				osGetwd = func() (string, error) {
					return "/test/dir", nil
				}

				exporterNew = func(rootDir string, flags *types.Flags) (*exporter.Exporter, error) {
					return nil, errors.New("failed to create exporter")
				}

				parseFlagsFunc = func() *types.Flags {
					return &types.Flags{Verbose: false}
				}
			},
			expectedOutput: "Failed to create exporter",
			expectedExit:   1,
		},
		{
			name: "Failed to export repository contents",
			setupMocks: func(t *testing.T) {
				mockExporter := new(MockExporter)
				mockExporter.On("Export").Return(errors.New("failed to export repository contents"))

				osGetwd = func() (string, error) {
					return "/test/dir", nil
				}

				exporterNew = func(rootDir string, flags *types.Flags) (*exporter.Exporter, error) {
					// We're returning a mock that satisfies the Exporter interface
					return &exporter.Exporter{}, nil
				}

				parseFlagsFunc = func() *types.Flags {
					return &types.Flags{Verbose: false}
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
			log.Logger = zerolog.New(&buf)

			tt.setupMocks(t)

			// Mock os.Exit
			var got int
			osExit = func(code int) {
				got = code
			}

			// Run main
			main()

			// Check exit code
			assert.Equal(t, tt.expectedExit, got)

			// Check log output
			output := buf.String()
			if tt.expectedOutput != "" {
				assert.Contains(t, output, tt.expectedOutput)
			}
		})
	}
}
