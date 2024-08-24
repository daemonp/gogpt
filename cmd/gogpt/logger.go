package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setupLogger(verbose bool) {
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: zerolog.TimeFormatUnix}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()

	if !isTerminal() {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true})
	}
}

func isTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
