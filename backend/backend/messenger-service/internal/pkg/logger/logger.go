package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}

	log.Logger = zerolog.New(output).With().Caller().Logger()
}