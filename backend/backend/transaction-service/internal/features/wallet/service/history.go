package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
)

// TransactionHistoryItem ре-экспортируем для удобства handler'а —
// чтобы он не знал про repository пакет.
type TransactionHistoryItem = repository.TransactionHistoryItem

// GetHistory возвращает последние транзакции, в которых участвовал юзер.
// limit / offset валидируются на уровне handler'а (clamp до 1-100 / 0-1M).
func (s *Service) GetHistory(
	ctx context.Context, userID uuid.UUID, limit, offset int,
) ([]TransactionHistoryItem, error) {
	items, err := s.repo.GetHistoryForUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("history: %w", err)
	}
	return items, nil
}