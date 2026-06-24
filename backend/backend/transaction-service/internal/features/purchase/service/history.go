package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase"
)

// GetHistory возвращает страницу истории покупок конкретного юзера.
// limit/offset уже нормализованы caller'ом (handler clamps в [1..100], [0..1M]).
func (s *Service) GetHistory(ctx context.Context, userID uuid.UUID, limit, offset int) ([]purchase.Purchase, error) {
	items, err := s.repo.ListForUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get history: %w", err)
	}
	return items, nil
}
