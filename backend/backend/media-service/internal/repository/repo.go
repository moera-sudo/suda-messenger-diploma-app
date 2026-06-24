package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/media-service/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, media *domain.Media) error
	GetByID(ctx context.Context, id string) (*domain.Media, error)
	GetByObjectKey(ctx context.Context, objectKey string) (*domain.Media, error)
	MarkReady(ctx context.Context, id string, sizeBytes int64, checksum string) error
	SoftDelete(ctx context.Context, id string) error
	CreateEntityLink(ctx context.Context, link *domain.MediaEntityLink) error
	GetEntityLinks(ctx context.Context, mediaID string) ([]domain.MediaEntityLink, error)
	GetMediaByEntity(ctx context.Context, entityType, entityID string) ([]domain.Media, error)
}

type mediaRepo struct {
	pool *pgxpool.Pool
}

func NewMediaRepository(pool *pgxpool.Pool) Repository {
	return &mediaRepo{pool: pool}
}

// ── Media CRUD ──────────────────────────────────────────────

func (r *mediaRepo) Create(ctx context.Context, media *domain.Media) error {
	query := `
		INSERT INTO media (id, owner_user_id, bucket, object_key, original_name, kind, content_type, size_bytes, status, is_private, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.pool.Exec(ctx, query,
		media.ID, media.OwnerUserID, media.Bucket, media.ObjectKey,
		media.OriginalName, media.Kind, media.ContentType, media.SizeBytes,
		media.Status, media.IsPrivate, media.CreatedAt,
	)
	return err
}

func (r *mediaRepo) GetByID(ctx context.Context, id string) (*domain.Media, error) {
	query := `
		SELECT id, owner_user_id, bucket, object_key, original_name, kind,
		       content_type, size_bytes, checksum_sha256, status, is_private,
		       created_at, ready_at, deleted_at
		FROM media WHERE id = $1`

	m := &domain.Media{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.OwnerUserID, &m.Bucket, &m.ObjectKey, &m.OriginalName,
		&m.Kind, &m.ContentType, &m.SizeBytes, &m.ChecksumSHA256,
		&m.Status, &m.IsPrivate, &m.CreatedAt, &m.ReadyAt, &m.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("media not found: %s", id)
		}
		return nil, err
	}
	return m, nil
}

func (r *mediaRepo) GetByObjectKey(ctx context.Context, objectKey string) (*domain.Media, error) {
	query := `
		SELECT id, owner_user_id, bucket, object_key, original_name, kind,
		       content_type, size_bytes, checksum_sha256, status, is_private,
		       created_at, ready_at, deleted_at
		FROM media WHERE object_key = $1`

	m := &domain.Media{}
	err := r.pool.QueryRow(ctx, query, objectKey).Scan(
		&m.ID, &m.OwnerUserID, &m.Bucket, &m.ObjectKey, &m.OriginalName,
		&m.Kind, &m.ContentType, &m.SizeBytes, &m.ChecksumSHA256,
		&m.Status, &m.IsPrivate, &m.CreatedAt, &m.ReadyAt, &m.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("media not found by key: %s", objectKey)
		}
		return nil, err
	}
	return m, nil
}

func (r *mediaRepo) MarkReady(ctx context.Context, id string, sizeBytes int64, checksum string) error {
	now := time.Now()
	query := `
		UPDATE media
		SET status = $1, size_bytes = $2, checksum_sha256 = $3, ready_at = $4
		WHERE id = $5 AND status = 'PENDING'`

	tag, err := r.pool.Exec(ctx, query, domain.StatusReady, sizeBytes, checksum, now, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("media %s is not in PENDING status or not found", id)
	}
	return nil
}

func (r *mediaRepo) SoftDelete(ctx context.Context, id string) error {
	now := time.Now()
	query := `UPDATE media SET status = $1, deleted_at = $2 WHERE id = $3`
	_, err := r.pool.Exec(ctx, query, domain.StatusDeleted, now, id)
	return err
}

// ── Entity Links ────────────────────────────────────────────

func (r *mediaRepo) CreateEntityLink(ctx context.Context, link *domain.MediaEntityLink) error {
	query := `
		INSERT INTO media_entity_links (id, media_id, entity_type, entity_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (media_id, entity_type, entity_id) DO NOTHING`

	_, err := r.pool.Exec(ctx, query,
		link.ID, link.MediaID, link.EntityType, link.EntityID, link.CreatedAt,
	)
	return err
}

func (r *mediaRepo) GetEntityLinks(ctx context.Context, mediaID string) ([]domain.MediaEntityLink, error) {
	query := `SELECT id, media_id, entity_type, entity_id, created_at FROM media_entity_links WHERE media_id = $1`

	rows, err := r.pool.Query(ctx, query, mediaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []domain.MediaEntityLink
	for rows.Next() {
		var l domain.MediaEntityLink
		if err := rows.Scan(&l.ID, &l.MediaID, &l.EntityType, &l.EntityID, &l.CreatedAt); err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	return links, rows.Err()
}

func (r *mediaRepo) GetMediaByEntity(ctx context.Context, entityType, entityID string) ([]domain.Media, error) {
	query := `
		SELECT m.id, m.owner_user_id, m.bucket, m.object_key, m.original_name, m.kind,
		       m.content_type, m.size_bytes, m.checksum_sha256, m.status, m.is_private,
		       m.created_at, m.ready_at, m.deleted_at
		FROM media m
		JOIN media_entity_links mel ON m.id = mel.media_id
		WHERE mel.entity_type = $1 AND mel.entity_id = $2 AND m.status != 'DELETED'`

	rows, err := r.pool.Query(ctx, query, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var medias []domain.Media
	for rows.Next() {
		var m domain.Media
		if err := rows.Scan(
			&m.ID, &m.OwnerUserID, &m.Bucket, &m.ObjectKey, &m.OriginalName,
			&m.Kind, &m.ContentType, &m.SizeBytes, &m.ChecksumSHA256,
			&m.Status, &m.IsPrivate, &m.CreatedAt, &m.ReadyAt, &m.DeletedAt,
		); err != nil {
			return nil, err
		}
		medias = append(medias, m)
	}
	return medias, rows.Err()
}