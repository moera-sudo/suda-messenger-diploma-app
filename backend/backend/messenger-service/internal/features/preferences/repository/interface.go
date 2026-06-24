package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/preferences"
)

type PreferencesRepo interface {
	Get(ctx context.Context, userID uuid.UUID) (*preferences.UserPreferences, error)
	Create(ctx context.Context, p *preferences.UserPreferences) error
	Update(ctx context.Context, userID uuid.UUID, req *preferences.UpdatePreferencesReq) error
}

type LastSeenHiddenRepo interface {
	Add(ctx context.Context, ownerID, hiddenID uuid.UUID) error
	Remove(ctx context.Context, ownerID, hiddenID uuid.UUID) error
	IsHidden(ctx context.Context, ownerID, hiddenID uuid.UUID) (bool, error)
	GetAll(ctx context.Context, ownerID uuid.UUID) ([]uuid.UUID, error)
}

func NewPreferencesRepo(db *pgxpool.Pool) PreferencesRepo {
	return &preferencesRepo{db: db}
}

func NewLastSeenHiddenRepo(db *pgxpool.Pool) LastSeenHiddenRepo {
	return &lastSeenHiddenRepo{db: db}
}