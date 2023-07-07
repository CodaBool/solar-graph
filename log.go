package main

import (
	"os"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func buildLogger() {
	o := zerolog.ConsoleWriter{Out: os.Stdout, PartsExclude: []string{zerolog.TimestampFieldName}}
	log = zerolog.New(o).With().Logger()
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
}

func check(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
