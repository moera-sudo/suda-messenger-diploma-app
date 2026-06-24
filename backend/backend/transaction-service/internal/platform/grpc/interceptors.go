package grpc

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// loggingInterceptor пишет в zerolog каждый unary-вызов: метод, длительность,
// статус-код. Ошибки логируются как Error, успехи как Debug.
func loggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()

	// recover защищает от panic'ов в handler'ах. Без него panic
	// положит весь gRPC сервер.
	resp, err := func() (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().
					Interface("panic", r).
					Str("method", info.FullMethod).
					Str("stack", string(debug.Stack())).
					Msg("gRPC handler panicked")
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}()

	duration := time.Since(start)
	if err != nil {
		log.Error().
			Err(err).
			Str("method", info.FullMethod).
			Dur("duration", duration).
			Msg("gRPC call failed")
	} else {
		log.Debug().
			Str("method", info.FullMethod).
			Dur("duration", duration).
			Msg("gRPC call ok")
	}

	return resp, err
}