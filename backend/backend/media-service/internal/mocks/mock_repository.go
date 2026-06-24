package mocks

import (
	"context"

	"github.com/moera-sudo/backend/backend/media-service/internal/domain"
)

type MockRepository struct {
	CreateFn           func(ctx context.Context, media *domain.Media) error
	GetByIDFn          func(ctx context.Context, id string) (*domain.Media, error)
	GetByObjectKeyFn   func(ctx context.Context, objectKey string) (*domain.Media, error)
	MarkReadyFn        func(ctx context.Context, id string, sizeBytes int64, checksum string) error
	SoftDeleteFn       func(ctx context.Context, id string) error
	CreateEntityLinkFn func(ctx context.Context, link *domain.MediaEntityLink) error
	GetEntityLinksFn   func(ctx context.Context, mediaID string) ([]domain.MediaEntityLink, error)
	GetMediaByEntityFn func(ctx context.Context, entityType, entityID string) ([]domain.Media, error)

	CreateCalls           int
	GetByIDCalls          int
	MarkReadyCalls        int
	SoftDeleteCalls       int
	CreateEntityLinkCalls int
	GetEntityLinksCalls   int
}

func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

func (m *MockRepository) Create(ctx context.Context, media *domain.Media) error {
	m.CreateCalls++
	if m.CreateFn != nil {
		return m.CreateFn(ctx, media)
	}
	return nil
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*domain.Media, error) {
	m.GetByIDCalls++
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockRepository) GetByObjectKey(ctx context.Context, objectKey string) (*domain.Media, error) {
	if m.GetByObjectKeyFn != nil {
		return m.GetByObjectKeyFn(ctx, objectKey)
	}
	return nil, nil
}

func (m *MockRepository) MarkReady(ctx context.Context, id string, sizeBytes int64, checksum string) error {
	m.MarkReadyCalls++
	if m.MarkReadyFn != nil {
		return m.MarkReadyFn(ctx, id, sizeBytes, checksum)
	}
	return nil
}

func (m *MockRepository) SoftDelete(ctx context.Context, id string) error {
	m.SoftDeleteCalls++
	if m.SoftDeleteFn != nil {
		return m.SoftDeleteFn(ctx, id)
	}
	return nil
}

func (m *MockRepository) CreateEntityLink(ctx context.Context, link *domain.MediaEntityLink) error {
	m.CreateEntityLinkCalls++
	if m.CreateEntityLinkFn != nil {
		return m.CreateEntityLinkFn(ctx, link)
	}
	return nil
}

func (m *MockRepository) GetEntityLinks(ctx context.Context, mediaID string) ([]domain.MediaEntityLink, error) {
	m.GetEntityLinksCalls++
	if m.GetEntityLinksFn != nil {
		return m.GetEntityLinksFn(ctx, mediaID)
	}
	return nil, nil
}

func (m *MockRepository) GetMediaByEntity(ctx context.Context, entityType, entityID string) ([]domain.Media, error) {
	if m.GetMediaByEntityFn != nil {
		return m.GetMediaByEntityFn(ctx, entityType, entityID)
	}
	return nil, nil
}