package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	chatRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat/repository"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/userpins"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/userpins/repository"
)

type Service struct {
	repo    repository.Repository
	members chatRepo.MemberRepo
}

func NewService(repo repository.Repository, members chatRepo.MemberRepo) *Service {
	return &Service{repo: repo, members: members}
}

func (s *Service) CreatePin(ctx context.Context, userID uuid.UUID, req *userpins.CreatePinReq) (*userpins.UserPin, error) {
	// Валидация комбинации
	if !userpins.IsValidCombination(req.PinType, req.TargetType) {
		return nil, fmt.Errorf("invalid combination: %s with %s", req.PinType, req.TargetType)
	}

	// Если закрепляем чат — проверяем что юзер в нём состоит
	if req.TargetType == userpins.TargetTypeChat {
		isMember, err := s.members.IsMember(ctx, req.TargetID, userID)
		if err != nil {
			return nil, err
		}
		if !isMember {
			return nil, fmt.Errorf("forbidden: not a member of this chat")
		}
	}

	// Если APP — пока не валидируем (Application Service будет позже)
	// TODO: when AppService is ready, check that app_id exists

	// Получаем следующий sort_order
	maxOrder, err := s.repo.GetMaxSortOrder(ctx, userID, req.PinType)
	if err != nil {
		return nil, err
	}

	pin := &userpins.UserPin{
		UserID:     userID,
		PinType:    req.PinType,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		SortOrder:  maxOrder + 1,
	}

	if err := s.repo.Create(ctx, pin); err != nil {
		return nil, fmt.Errorf("create pin: %w", err)
	}

	return pin, nil
}

func (s *Service) DeletePin(ctx context.Context, userID uuid.UUID, pinID int64) error {
	return s.repo.Delete(ctx, userID, pinID)
}

func (s *Service) DeletePinByTarget(ctx context.Context, userID, targetID uuid.UUID, pinType string) error {
	return s.repo.DeleteByTarget(ctx, userID, targetID, pinType)
}

func (s *Service) GetPins(ctx context.Context, userID uuid.UUID, pinType string) ([]userpins.UserPin, error) {
	if pinType == "" {
		return nil, fmt.Errorf("pin_type query param is required")
	}
	return s.repo.GetByUser(ctx, userID, pinType)
}

func (s *Service) GetPinnedChatIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	return s.repo.GetPinnedChatIDs(ctx, userID)
}

func (s *Service) Reorder(ctx context.Context, userID uuid.UUID, items []userpins.ReorderItem) error {
	if len(items) == 0 {
		return nil
	}
	return s.repo.Reorder(ctx, userID, items)
}