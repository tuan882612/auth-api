package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"authapi/internal/server"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	deps, err := server.NewDependencies()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create dependencies")
	}

	svr := server.New(deps)
	if err := svr.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
