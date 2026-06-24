package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/userpins"
)

type Repository interface {
	Create(ctx context.Context, pin *userpins.UserPin) error
	Delete(ctx context.Context, userID uuid.UUID, pinID int64) error
	DeleteByTarget(ctx context.Context, userID, targetID uuid.UUID, pinType string) error
	GetByUser(ctx context.Context, userID uuid.UUID, pinType string) ([]userpins.UserPin, error)
	GetPinnedChatIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	GetPinnedChats(ctx context.Context, userID uuid.UUID) ([]userpins.PinnedChat, error)
	Reorder(ctx context.Context, userID uuid.UUID, items []userpins.ReorderItem) error
	GetMaxSortOrder(ctx context.Context, userID uuid.UUID, pinType string) (int, error)
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repo{db: db}
}