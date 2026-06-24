// Package gating — token-gating правила доступа к каналам.
//
// Owner канала задаёт правило (минимальный баланс SUDA, опц. владение NFT
// из коллекции). messenger-service при подписке на канал дёргает gRPC
// CheckTokenGating (логика — в wallet/service/gating.go) и блокирует вход,
// если правило не выполнено. HTTP API этой фичи — CRUD самих правил.
package gating

import "errors"

// Доменные ошибки. Handler маппит их в HTTP-коды.
var (
	ErrInvalidInput = errors.New("invalid input")
	ErrForbidden    = errors.New("forbidden")
)
