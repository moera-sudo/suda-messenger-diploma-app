package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Setup инициализирует глобальный zerolog логгер.
//
// Уровень: Debug (включает Info/Warn/Error/Fatal/Panic).
// Формат:  ConsoleWriter в stdout с человекочитаемыми timestamps.
// Caller:  имя файла и строка вызова автоматически добавляется в каждое сообщение.
//
// Использование в коде:
//
//	import "github.com/rs/zerolog/log"
//	log.Info().Str("user_id", uid.String()).Msg("wallet created")
//	log.Error().Err(err).Str("addr", addr).Msg("balance fetch failed")
//	log.Warn().Int("retry", n).Msg("retrying broadcast")
//
// Важно: уровень Error должен использоваться ВЕЗДЕ где есть err != nil
// в критичных операциях (БД, блокчейн, gRPC). Info — для бизнес-событий
// (создан кошелёк, выполнен перевод). Warn — для подозрительных но не
// фатальных ситуаций (retry, deprecated path).
func Setup() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
	log.Logger = zerolog.New(output).With().Caller().Logger()
}