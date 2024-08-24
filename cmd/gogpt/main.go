package main

import (
	"os"

	"github.com/daemonp/gogpt/pkg/exporter"
	"github.com/rs/zerolog/log"
)

func main() {
	flags := parseFlags()
	setupLogger(flags.Verbose)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get current working directory")
	}

	exp, err := exporter.New(dir, flags)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create exporter")
	}

	if err := exp.Export(); err != nil {
		log.Fatal().Err(err).Msg("Failed to export repository contents")
	}

	log.Info().Msg("Export completed successfully")
}
