// File: cmd/gogpt/main.go

package main

import (
	"os"

	"github.com/daemonp/gogpt/cmd/gogpt/flags"
	"github.com/daemonp/gogpt/cmd/gogpt/logger"
	"github.com/daemonp/gogpt/pkg/exporter"
	"github.com/rs/zerolog/log"
)

var (
	osExit         = os.Exit
	osGetwd        = os.Getwd
	exporterNew    = exporter.New
	parseFlagsFunc = flags.ParseFlags
)

func main() {
	flags := parseFlagsFunc()
	logger.SetupLogger(flags.Verbose)

	dir, err := osGetwd()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current working directory")
		osExit(1)
		return
	}

	log.Debug().Str("dir", dir).Msg("Current working directory")

	exp, err := exporterNew(dir, flags)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create exporter")
		osExit(1)
		return
	}

	if exp == nil {
		log.Error().Msg("Exporter is nil")
		osExit(1)
		return
	}

	if err := exp.Export(); err != nil {
		log.Error().Err(err).Msg("Failed to export repository contents")
		osExit(1)
		return
	}

	log.Info().Msg("Export completed successfully")
}
