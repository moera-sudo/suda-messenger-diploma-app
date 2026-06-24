package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/media-service/config"
	"github.com/moera-sudo/backend/backend/media-service/internal/domain"
	"github.com/moera-sudo/backend/backend/media-service/internal/domain/models"
	"github.com/moera-sudo/backend/backend/media-service/internal/platform/clients/messenger"
	s3storage "github.com/moera-sudo/backend/backend/media-service/internal/platform/storage/s3"
	"github.com/moera-sudo/backend/backend/media-service/internal/repository"
)

type MediaService struct {
	repo          repository.Repository
	storage       s3storage.FileStorage
	accessChecker messenger.AccessChecker
	cfg           *config.Config
}

func NewMediaService(
	repo repository.Repository,
	storage s3storage.FileStorage,
	accessChecker messenger.AccessChecker,
	cfg *config.Config,
) *MediaService {
	return &MediaService{
		repo:          repo,
		storage:       storage,
		accessChecker: accessChecker,
		cfg:           cfg,
	}
}

// InitUpload создаёт запись PENDING в БД и возвращает presigned PUT URL
// Клиент загружает файл напрямую в S3 по этой ссылке
func (s *MediaService) InitUpload(ctx context.Context, req models.InitUploadRequest, presignHost string) (*models.InitUploadResponse, error) {
	if err := s.validateKind(req.Kind, req.ContentType); err != nil {
		return nil, err
	}

	mediaID := uuid.New().String()
	ext := extensionFromContentType(req.ContentType)
	objectKey := fmt.Sprintf("%s/%s/%s%s", strings.ToLower(req.Kind), time.Now().Format("2006/01/02"), mediaID, ext)
	bucket := s.bucketForKind(req.Kind)

	if err := s.storage.EnsureBucket(ctx, bucket); err != nil {
		return nil, fmt.Errorf("ensure bucket: %w", err)
	}

	media := &domain.Media{
		ID:          mediaID,
		OwnerUserID: req.OwnerUserID,
		Bucket:      bucket,
		ObjectKey:   objectKey,
		Kind:        req.Kind,
		ContentType: req.ContentType,
		SizeBytes:   req.SizeBytes,
		Status:      domain.StatusPending,
		IsPrivate:   req.IsPrivate,
		CreatedAt:   time.Now(),
	}
	if req.OriginalName != "" {
		media.OriginalName = &req.OriginalName
	}

	if err := s.repo.Create(ctx, media); err != nil {
		log.Error().Err(err).Str("media_id", mediaID).Str("kind", req.Kind).Msg("Failed to create media record")
		return nil, fmt.Errorf("create media record: %w", err)
	}

	uploadURL, err := s.storage.PresignPutURL(ctx, presignHost, bucket, objectKey, req.ContentType, s.cfg.PresignUploadExpiry)
	if err != nil {
		log.Error().Err(err).Str("media_id", mediaID).Msg("Failed to generate presigned upload URL")
		return nil, fmt.Errorf("generate upload url: %w", err)
	}

	log.Debug().Str("media_id", mediaID).Str("kind", req.Kind).Str("content_type", req.ContentType).Msg("Upload initialized")

	return &models.InitUploadResponse{
		MediaID:   mediaID,
		UploadURL: uploadURL,
		ObjectKey: objectKey,
		ExpiresIn: int(s.cfg.PresignUploadExpiry.Seconds()),
	}, nil
}

// ConfirmUpload проверяет что файл реально загружен в S3 и обновляет статус на READY
func (s *MediaService) ConfirmUpload(ctx context.Context, mediaID string, ownerUserID uuid.UUID) (*models.ConfirmUploadResponse, error) {
	media, err := s.repo.GetByID(ctx, mediaID)
	if err != nil {
		return nil, fmt.Errorf("media not found: %w", err)
	}

	if media.OwnerUserID != ownerUserID {
		return nil, fmt.Errorf("forbidden: not the owner")
	}

	if media.Status != domain.StatusPending {
		return nil, fmt.Errorf("media is not in PENDING status, current: %s", media.Status)
	}

	size, etag, err := s.storage.HeadObject(ctx, media.Bucket, media.ObjectKey)
	if err != nil {
		return nil, fmt.Errorf("file not found in storage: %w", err)
	}

	if err := s.repo.MarkReady(ctx, mediaID, size, etag); err != nil {
		log.Error().Err(err).Str("media_id", mediaID).Msg("Failed to mark media as ready")
		return nil, fmt.Errorf("mark ready: %w", err)
	}

	log.Info().Str("media_id", mediaID).Int64("size_bytes", size).Msg("Upload confirmed")

	return &models.ConfirmUploadResponse{
		MediaID:   mediaID,
		Status:    domain.StatusReady,
		SizeBytes: size,
	}, nil
}

// GetViewURL — presigned URL для просмотра (Content-Disposition: inline)
// Используется для показа картинок/видео/аудио прямо в чате
func (s *MediaService) GetViewURL(ctx context.Context, mediaID string, requesterUserID uuid.UUID, presignHost string) (*models.GetURLResponse, error) {
	media, err := s.getReadyMedia(ctx, mediaID)
	if err != nil {
		return nil, err
	}

	if media.IsPrivate {
		if err := s.checkAccess(ctx, media, requesterUserID); err != nil {
			return nil, err
		}
	}

	url, err := s.storage.PresignGetURL(ctx, presignHost, media.Bucket, media.ObjectKey, s.cfg.PresignViewExpiry)
	if err != nil {
		return nil, fmt.Errorf("generate view url: %w", err)
	}

	return &models.GetURLResponse{
		MediaID:   mediaID,
		URL:       url,
		ExpiresIn: int(s.cfg.PresignViewExpiry.Seconds()),
	}, nil
}

// GetDownloadURL — presigned URL для скачивания (Content-Disposition: attachment)
// Используется при нажатии кнопки "Скачать" в чате
func (s *MediaService) GetDownloadURL(ctx context.Context, mediaID string, requesterUserID uuid.UUID, presignHost string) (*models.GetURLResponse, error) {
	media, err := s.getReadyMedia(ctx, mediaID)
	if err != nil {
		return nil, err
	}

	if media.IsPrivate {
		if err := s.checkAccess(ctx, media, requesterUserID); err != nil {
			return nil, err
		}
	}

	filename := "download"
	if media.OriginalName != nil {
		filename = *media.OriginalName
	}

	url, err := s.storage.PresignGetDownloadURL(ctx, presignHost, media.Bucket, media.ObjectKey, filename, s.cfg.PresignDownloadExpiry)
	if err != nil {
		return nil, fmt.Errorf("generate download url: %w", err)
	}

	return &models.GetURLResponse{
		MediaID:   mediaID,
		URL:       url,
		ExpiresIn: int(s.cfg.PresignDownloadExpiry.Seconds()),
	}, nil
}

func (s *MediaService) GetMetadata(ctx context.Context, mediaID string, requesterUserID uuid.UUID) (*models.MediaMetadata, error) {
	media, err := s.repo.GetByID(ctx, mediaID)
	if err != nil {
		return nil, fmt.Errorf("media not found: %w", err)
	}

	if media.IsDeleted() {
		return nil, fmt.Errorf("media is deleted")
	}

	if media.IsPrivate {
		if err := s.checkAccess(ctx, media, requesterUserID); err != nil {
			return nil, err
		}
	}

	return &models.MediaMetadata{
		ID:           media.ID,
		OwnerUserID:  media.OwnerUserID,
		Kind:         media.Kind,
		ContentType:  media.ContentType,
		SizeBytes:    media.SizeBytes,
		OriginalName: media.OriginalName,
		Status:       media.Status,
		IsPrivate:    media.IsPrivate,
		CreatedAt:    media.CreatedAt,
		ReadyAt:      media.ReadyAt,
	}, nil
}

func (s *MediaService) LinkToEntity(ctx context.Context, mediaID, entityType, entityID string) error {
	link := &domain.MediaEntityLink{
		ID:         uuid.New().String(),
		MediaID:    mediaID,
		EntityType: entityType,
		EntityID:   entityID,
		CreatedAt:  time.Now(),
	}
	if err := s.repo.CreateEntityLink(ctx, link); err != nil {
		log.Error().Err(err).Str("media_id", mediaID).Str("entity_type", entityType).Str("entity_id", entityID).Msg("Failed to link media to entity")
		return err
	}

	log.Debug().Str("media_id", mediaID).Str("entity_type", entityType).Str("entity_id", entityID).Msg("Media linked to entity")
	return nil
}

func (s *MediaService) Delete(ctx context.Context, mediaID string, requesterUserID uuid.UUID) error {
	media, err := s.repo.GetByID(ctx, mediaID)
	if err != nil {
		return fmt.Errorf("media not found: %w", err)
	}

	if media.OwnerUserID != requesterUserID {
		log.Warn().Str("media_id", mediaID).Str("requester_id", requesterUserID.String()).Msg("Delete denied: not the owner")
		return fmt.Errorf("forbidden: not the owner")
	}

	if err := s.repo.SoftDelete(ctx, mediaID); err != nil {
		log.Error().Err(err).Str("media_id", mediaID).Msg("Failed to soft-delete media")
		return fmt.Errorf("soft delete: %w", err)
	}

	// Удаляем из S3 (можно сделать асинхронно позже)
	if err := s.storage.DeleteObject(ctx, media.Bucket, media.ObjectKey); err != nil {
		log.Warn().Err(err).
			Str("bucket", media.Bucket).
			Str("key", media.ObjectKey).
			Msg("Failed to delete S3 object (DB already marked as deleted)")
	}

	log.Info().Str("media_id", mediaID).Str("user_id", requesterUserID.String()).Msg("Media deleted")

	return nil
}

func (s *MediaService) getReadyMedia(ctx context.Context, mediaID string) (*domain.Media, error) {
	media, err := s.repo.GetByID(ctx, mediaID)
	if err != nil {
		return nil, fmt.Errorf("media not found: %w", err)
	}
	if media.IsDeleted() {
		return nil, fmt.Errorf("media is deleted")
	}
	if !media.IsReady() {
		return nil, fmt.Errorf("media is not ready, status: %s", media.Status)
	}
	return media, nil
}

func (s *MediaService) checkAccess(ctx context.Context, media *domain.Media, userID uuid.UUID) error {
	if media.OwnerUserID == userID {
		return nil
	}

	links, err := s.repo.GetEntityLinks(ctx, media.ID)
	if err != nil {
		return fmt.Errorf("get entity links: %w", err)
	}

	if len(links) == 0 {
		return fmt.Errorf("forbidden: private media without entity links")
	}

	for _, link := range links {
		hasAccess, err := s.accessChecker.CheckAccess(ctx, userID.String(), link.EntityType, link.EntityID)
		//                                                  ^^^^^^^^^^^^^^^^ конвертим в string для gRPC
		if err != nil {
			log.Warn().Err(err).
				Str("entity_type", link.EntityType).
				Str("entity_id", link.EntityID).
				Msg("Access check failed, trying next link")
			continue
		}
		if hasAccess {
			return nil
		}
	}

	return fmt.Errorf("forbidden: user does not have access to this media")
}

func (s *MediaService) validateKind(kind, contentType string) error {
	allowed, exists := domain.AllowedContentTypes[kind]
	if !exists {
		return fmt.Errorf("unknown media kind: %s", kind)
	}

	if len(allowed) == 0 {
		return nil
	}

	for _, prefix := range allowed {
		if strings.HasPrefix(contentType, prefix) {
			return nil
		}
	}

	return fmt.Errorf("content type %s is not allowed for kind %s", contentType, kind)
}

func (s *MediaService) bucketForKind(kind string) string {
	switch kind {
	case domain.KindAvatar, domain.KindChatAvatar:
		return "avatars"
	default:
		return s.cfg.S3Bucket
	}
}

func extensionFromContentType(ct string) string {
	extensions := map[string]string{
		"image/jpeg":      ".jpg",
		"image/png":       ".png",
		"image/gif":       ".gif",
		"image/webp":      ".webp",
		"audio/ogg":       ".ogg",
		"audio/mpeg":      ".mp3",
		"audio/wav":       ".wav",
		"video/mp4":       ".mp4",
		"video/webm":      ".webm",
		"application/pdf": ".pdf",
	}

	if ext, ok := extensions[ct]; ok {
		return ext
	}

	parts := strings.Split(ct, "/")
	if len(parts) == 2 {
		return "." + parts[1]
	}
	return ""
}
