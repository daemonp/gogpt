// File: cmd/gogpt/logger/logger_test.go
package logger

import (
	"bytes"
	"encoding/json"
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
			expectedOutput: []string{"debug", "debug message", "info message"},
		},
		{
			name:           "Non-verbose logging",
			verbose:        false,
			expectedLevel:  zerolog.InfoLevel,
			expectedOutput: []string{"info", "info message"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var buf bytes.Buffer
			log.Logger = zerolog.New(&buf)

			SetupLogger(tt.verbose)

			// Check log level
			assert.Equal(t, tt.expectedLevel, zerolog.GlobalLevel())

			// Log messages
			log.Debug().Msg("debug message")
			log.Info().Msg("info message")

			// Parse and check log output
			output := buf.String()
			lines := strings.Split(strings.TrimSpace(output), "\n")
			for i, line := range lines {
				var logEntry map[string]interface{}
				err := json.Unmarshal([]byte(line), &logEntry)
				assert.NoError(t, err, "Failed to parse log entry")

				level, ok := logEntry["level"].(string)
				assert.True(t, ok, "Log entry should have a 'level' field")

				message, ok := logEntry["message"].(string)
				assert.True(t, ok, "Log entry should have a 'message' field")

				assert.Contains(t, tt.expectedOutput[i], strings.ToLower(level), "Unexpected log level")
				assert.Contains(t, message, tt.expectedOutput[i], "Unexpected log message")
			}
		})
	}
}
