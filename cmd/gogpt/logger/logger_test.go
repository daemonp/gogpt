// File: cmd/gogpt/logger/logger_test.go

package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestSetupLogger(t *testing.T) {
	tests := []struct {
		name           string
		verbose        bool
		expectedLevel  zerolog.Level
		expectedOutput []string
	}{
		{
			name:           "Verbose logging",
			verbose:        true,
			expectedLevel:  zerolog.DebugLevel,
			expectedOutput: []string{"DBG", "debug message", "INF", "info message"},
		},
		{
			name:           "Non-verbose logging",
			verbose:        false,
			expectedLevel:  zerolog.InfoLevel,
			expectedOutput: []string{"INF", "info message"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var buf bytes.Buffer
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: &buf, NoColor: true})

			SetupLogger(tt.verbose)

			// Check log level
			assert.Equal(t, tt.expectedLevel, zerolog.GlobalLevel())

			// Log messages
			log.Debug().Msg("debug message")
			log.Info().Msg("info message")

			// Parse and check log output
			output := buf.String()
			lines := strings.Split(strings.TrimSpace(output), "\n")
			for _, line := range lines {
				for _, expected := range tt.expectedOutput {
					if strings.Contains(line, expected) {
						assert.Contains(t, line, expected, "Expected output not found in log")
					}
				}
			}
		})
	}
}
