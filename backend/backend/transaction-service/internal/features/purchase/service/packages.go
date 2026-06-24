package service

import (
	"context"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase"
)

// ListPackages возвращает доступные пакеты SUDA для покупки.
// Сейчас данные захардкожены в коде; ctx добавлен на случай миграции
// в БД позже.
func (s *Service) ListPackages(_ context.Context) []purchase.Package {
	return purchase.Packages
}
