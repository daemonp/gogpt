package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/daemonp/gogpt/internal/exporter"
	"github.com/daemonp/gogpt/internal/languagedetector"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Set up command-line flags
	outputFile := flag.String("f", "", "Output file path (default: stdout)")
	ignoreGitIgnore := flag.Bool("i", false, "Ignore files specified in .gitignore")
	languages := flag.String("l", "", "Comma-separated list of languages to include (e.g., 'go,js,md')")
	maxTokensStr := flag.String("t", "1000", "Maximum number of tokens per file")
	flag.Parse()

	// Convert maxTokens to an integer
	maxTokens, err := strconv.Atoi(*maxTokensStr)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid value for max tokens")
	}

	// Configure zerolog
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: zerolog.TimeFormatUnix}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()

	// Determine if output is being piped
	if !isTerminal() {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true})
	}

	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get current working directory")
	}

	// Detect repository type if no languages are specified
	if *languages == "" {
		detectedLangs := languagedetector.DetectLanguages(dir)
		*languages = detectedLangs
		log.Info().Str("languages", *languages).Msg("Detected languages")
	}

	// Create exporter with maxTokens
	exp, err := exporter.New(dir, *outputFile, *ignoreGitIgnore, *languages, maxTokens)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create exporter")
	}

	// Export repository contents
	err = exp.Export()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to export repository contents")
	}

	log.Info().Msg("Export completed successfully")
}

func isTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
