package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/moera-sudo/backend/backend/media-service/config"
	"github.com/moera-sudo/backend/backend/media-service/internal/domain"
	"github.com/moera-sudo/backend/backend/media-service/internal/domain/models"
	"github.com/moera-sudo/backend/backend/media-service/internal/mocks"
	"github.com/moera-sudo/backend/backend/media-service/internal/service"
)

func testConfig() *config.Config {
	return &config.Config{
		S3Bucket:              "media",
		PresignUploadExpiry:   15 * time.Minute,
		PresignViewExpiry:     4 * time.Hour,
		PresignDownloadExpiry: 1 * time.Hour,
	}
}

func setup() (*service.MediaService, *mocks.MockRepository, *mocks.MockStorage, *mocks.MockAccessChecker) {
	repo := mocks.NewMockRepository()
	storage := mocks.NewMockStorage()
	access := mocks.NewMockAccessChecker()
	svc := service.NewMediaService(repo, storage, access, testConfig())
	return svc, repo, storage, access
}

var (
	userID1 = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	userID2 = uuid.MustParse("22222222-2222-2222-2222-222222222222")
)

// ═══════════════════════════════════════════════════════════
//  InitUpload
// ═══════════════════════════════════════════════════════════

func TestInitUpload_Success(t *testing.T) {
	svc, repo, storage, _ := setup()

	req := models.InitUploadRequest{
		OwnerUserID:  userID1,
		Kind:         domain.KindAttachment,
		ContentType:  "image/jpeg",
		OriginalName: "photo.jpg",
		IsPrivate:    false,
	}

	resp, err := svc.InitUpload(context.Background(), req, "")

	require.NoError(t, err)
	assert.NotEmpty(t, resp.MediaID)
	assert.NotEmpty(t, resp.UploadURL)
	assert.Contains(t, resp.ObjectKey, "attachment/")
	assert.Contains(t, resp.ObjectKey, ".jpg")
	assert.Equal(t, 900, resp.ExpiresIn)
	assert.Equal(t, 1, repo.CreateCalls)
	assert.Equal(t, 1, storage.EnsureBucketCalls)
	assert.Equal(t, 1, storage.PresignPutURLCalls)
}

func TestInitUpload_AvatarBucket(t *testing.T) {
	svc, _, storage, _ := setup()

	var capturedBucket string
	storage.EnsureBucketFn = func(_ context.Context, bucket string) error {
		capturedBucket = bucket
		return nil
	}

	req := models.InitUploadRequest{
		OwnerUserID: userID1,
		Kind:        domain.KindAvatar,
		ContentType: "image/png",
	}

	_, err := svc.InitUpload(context.Background(), req, "")

	require.NoError(t, err)
	assert.Equal(t, "avatars", capturedBucket)
}

func TestInitUpload_InvalidKind(t *testing.T) {
	svc, _, _, _ := setup()

	req := models.InitUploadRequest{
		OwnerUserID: userID1,
		Kind:        "NONEXISTENT",
		ContentType: "image/jpeg",
	}

	resp, err := svc.InitUpload(context.Background(), req, "")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unknown media kind")
}

func TestInitUpload_WrongContentTypeForAvatar(t *testing.T) {
	svc, _, _, _ := setup()

	req := models.InitUploadRequest{
		OwnerUserID: userID1,
		Kind:        domain.KindAvatar,
		ContentType: "application/pdf",
	}

	resp, err := svc.InitUpload(context.Background(), req, "")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "not allowed")
}

func TestInitUpload_AttachmentAcceptsAnyType(t *testing.T) {
	svc, _, _, _ := setup()

	types := []string{"application/pdf", "video/mp4", "audio/wav", "text/plain"}
	for _, ct := range types {
		req := models.InitUploadRequest{
			OwnerUserID: userID1,
			Kind:        domain.KindAttachment,
			ContentType: ct,
		}
		resp, err := svc.InitUpload(context.Background(), req, "")
		assert.NoError(t, err, "should accept %s for ATTACHMENT", ct)
		assert.NotNil(t, resp)
	}
}

func TestInitUpload_VoiceMessageOnlyAudio(t *testing.T) {
	svc, _, _, _ := setup()

	req := models.InitUploadRequest{
		OwnerUserID: userID1,
		Kind:        domain.KindVoiceMessage,
		ContentType: "audio/ogg",
	}
	resp, err := svc.InitUpload(context.Background(), req, "")
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	req.ContentType = "image/jpeg"
	resp, err = svc.InitUpload(context.Background(), req, "")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestInitUpload_S3Failure(t *testing.T) {
	svc, _, storage, _ := setup()

	storage.EnsureBucketFn = func(_ context.Context, _ string) error {
		return fmt.Errorf("S3 unreachable")
	}

	req := models.InitUploadRequest{
		OwnerUserID: userID1,
		Kind:        domain.KindAttachment,
		ContentType: "image/jpeg",
	}

	resp, err := svc.InitUpload(context.Background(), req, "")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "ensure bucket")
}

func TestInitUpload_DBFailure(t *testing.T) {
	svc, repo, _, _ := setup()

	repo.CreateFn = func(_ context.Context, _ *domain.Media) error {
		return fmt.Errorf("connection refused")
	}

	req := models.InitUploadRequest{
		OwnerUserID: userID1,
		Kind:        domain.KindAttachment,
		ContentType: "image/jpeg",
	}

	resp, err := svc.InitUpload(context.Background(), req, "")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "create media record")
}

// ═══════════════════════════════════════════════════════════
//  ConfirmUpload
// ═══════════════════════════════════════════════════════════

func TestConfirmUpload_Success(t *testing.T) {
	svc, repo, storage, _ := setup()
	mediaID := uuid.New().String()

	repo.GetByIDFn = func(_ context.Context, id string) (*domain.Media, error) {
		return &domain.Media{
			ID: id, OwnerUserID: userID1,
			Bucket: "media", ObjectKey: "attachment/test.jpg",
			Status: domain.StatusPending,
		}, nil
	}

	storage.HeadObjectFn = func(_ context.Context, _, _ string) (int64, string, error) {
		return 30000, "etag-xyz", nil
	}

	resp, err := svc.ConfirmUpload(context.Background(), mediaID, userID1)

	require.NoError(t, err)
	assert.Equal(t, mediaID, resp.MediaID)
	assert.Equal(t, domain.StatusReady, resp.Status)
	assert.Equal(t, int64(30000), resp.SizeBytes)
	assert.Equal(t, 1, repo.MarkReadyCalls)
	assert.Equal(t, 1, storage.HeadObjectCalls)
}

func TestConfirmUpload_NotOwner(t *testing.T) {
	svc, repo, _, _ := setup()

	repo.GetByIDFn = func(_ context.Context, id string) (*domain.Media, error) {
		return &domain.Media{ID: id, OwnerUserID: userID1, Status: domain.StatusPending}, nil
	}

	resp, err := svc.ConfirmUpload(context.Background(), "id", userID2)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestConfirmUpload_AlreadyReady(t *testing.T) {
	svc, repo, _, _ := setup()

	repo.GetByIDFn = func(_ context.Context, id string) (*domain.Media, error) {
		return &domain.Media{ID: id, OwnerUserID: userID1, Status: domain.StatusReady}, nil
	}

	resp, err := svc.ConfirmUpload(context.Background(), "id", userID1)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "not in PENDING status")
}

func TestConfirmUpload_FileNotInS3(t *testing.T) {
	svc, repo, storage, _ := setup()

	repo.GetByIDFn = func(_ context.Context, id string) (*domain.Media, error) {
		return &domain.Media{
			ID: id, OwnerUserID: userID1,
			Bucket: "media", ObjectKey: "missing.jpg",
			Status: domain.StatusPending,
		}, nil
	}

	storage.HeadObjectFn = func(_ context.Context, _, _ string) (int64, string, error) {
		return 0, "", fmt.Errorf("object not found")
	}

	resp, err := svc.ConfirmUpload(context.Background(), "id", userID1)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "file not found in storage")
}

func TestConfirmUpload_MediaNotFound(t *testing.T) {
	svc, repo, _, _ := setup()

	repo.GetByIDFn = func(_ context.Context, _ string) (*domain.Media, error) {
		return nil, fmt.Errorf("media not found: fake-id")
	}

	resp, err := svc.ConfirmUpload(context.Background(), "fake-id", userID1)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

// ═══════════════════════════════════════════════════════════
//  GetViewURL
// ═══════════════════════════════════════════════════════════

func TestGetViewURL_PublicMedia(t *testing.T) {
	svc, repo, _, _ := setup()

	repo.GetByIDFn = func(_ context.Context, id string) (*domain.Media, error) {
		return &domain.Media{
			ID: id, Bucket: "media", ObjectKey: "photo.jpg",
			Status: domain.StatusReady, IsPrivate: false,
		}, nil
	}

	resp, err := svc.GetViewURL(context.Background(), "id", userID1, "")

	require.NoError(t, err)
	assert.NotEmpty(t, resp.URL)
	assert.Equal(t, 14400, resp.ExpiresIn)
}

func TestGetViewURL_PrivateMedia_OwnerAllowed(t *testing.T) {
	svc, repo, _, _ := setup()

	repo.GetByIDFn = func(_ context.Context, id string) (*domain.Media, error) {
		return &domain.Media{
			ID: id, OwnerUserID: userID1, Bucket: "media",
			ObjectKey: "private.jpg", Status: domain.StatusReady, IsPrivate: true,
		}, nil
	}

	resp, err := svc.GetViewURL(context.Background(), "id", userID1, "")

	require.NoError(t, err)
	assert.NotEmpty(t, resp.URL)
}

func TestGetViewURL_PrivateMedia_MemberAllowed(t *testing.T) {
	svc, repo, _, access := setup()

	repo.GetByIDFn = func(_ context.Context, id string) (*domain.Media, error) {
		return &domain.Media{
			ID: id, OwnerUserID: userID1, Bucket: "media",
			ObjectKey: "private.jpg", Status: domain.StatusReady, IsPrivate: true,
		}, nil
	}

	repo.GetEntityLinksFn = func(_ context.Context, _ string) ([]domain.MediaEntityLink, error) {
		return []domain.MediaEntityLink{
			{MediaID: "id", EntityType: "MESSAGE", EntityID: "456"},
		}, nil
	}

	access.CheckAccessFn = func(_ context.Context, _, _, _ string) (bool, error) {
		return true, nil
	}

	resp, err := svc.GetViewURL(context.Background(), "id", userID2, "")

	require.NoError(t, err)
	assert.NotEmpty(t, resp.URL)
	assert.Equal(t, 1, access.CheckAccessCalls)
}

func TestGetViewURL_PrivateMedia_Denied(t *testing.T) {
	svc, repo, _, access := setup()

	repo.GetByIDFn = func(_ context.Context, id string) (*domain.Media, error) {
		return &domain.Media{
			ID: id, OwnerUserID: userID1, Bucket: "media",
			ObjectKey: "private.jpg", Status: domain.StatusReady, IsPrivate: true,
		}, nil
	}

	repo.GetEntityLinksFn = func(_ context.Context, _ string) ([]domain.MediaEntityLink, error) {
		return []domain.MediaEntityLink{
			{MediaID: "id", EntityType: "MESSAGE", EntityID: "456"},
		}, nil
	}

	access.CheckAccessFn = func(_ context.Context, _, _, _ string) (bool, error) {
		return false, nil
	}

	resp, err := svc.GetViewURL(context.Background(), "id", userID2, "")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestGetViewURL_DeletedMedia(t *testing.T) {
	svc, repo, _, _ := setup()

	repo.GetByIDFn = func(_ context.Context, id string) (*domain.Media, error) {
		return &domain.Media{ID: id, Status: domain.StatusDeleted}, nil
	}

	resp, err := svc.GetViewURL(context.Background(), "id", userID1, "")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "deleted")
}

func TestGetViewURL_PendingMedia(t *testing.T) {
	svc, repo, _, _ := setup()

	repo.GetByIDFn = func(_ context.Context, id string) (*domain.Media, error) {
		return &domain.Media{ID: id, Status: domain.StatusPending}, nil
	}

	resp, err := svc.GetViewURL(context.Background(), "id", userID1, "")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "not ready")
}

func TestGetViewURL_NoEntityLinks_Denied(t *testing.T) {
	svc, repo, _, _ := setup()

	repo.GetByIDFn = func(_ context.Context, id string) (*domain.Media, error) {
		return &domain.Media{
			ID: id, OwnerUserID: userID1, Status: domain.StatusReady,
			IsPrivate: true, Bucket: "media", ObjectKey: "test.jpg",
		}, nil
	}

	repo.GetEntityLinksFn = func(_ context.Context, _ string) ([]domain.MediaEntityLink, error) {
		return nil, nil
	}

	resp, err := svc.GetViewURL(context.Background(), "id", userID2, "")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "without entity links")
}

func TestGetViewURL_AccessCheckError_TriesNextLink(t *testing.T) {
	svc, repo, _, access := setup()

	repo.GetByIDFn = func(_ context.Context, id string) (*domain.Media, error) {
		return &domain.Media{
			ID: id, OwnerUserID: userID1, Status: domain.StatusReady,
			IsPrivate: true, Bucket: "media", ObjectKey: "test.jpg",
		}, nil
	}

	repo.GetEntityLinksFn = func(_ context.Context, _ string) ([]domain.MediaEntityLink, error) {
		return []domain.MediaEntityLink{
			{EntityType: "MESSAGE", EntityID: "100"},
			{EntityType: "MESSAGE", EntityID: "200"},
		}, nil
	}

	callCount := 0
	access.CheckAccessFn = func(_ context.Context, _, _, entityID string) (bool, error) {
		callCount++
		if entityID == "100" {
			return false, fmt.Errorf("gRPC timeout")
		}
		return true, nil
	}

	resp, err := svc.GetViewURL(context.Background(), "id", userID2, "")

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, callCount)
}

// ═══════════════════════════════════════════════════════════
//  GetDownloadURL
// ═══════════════════════════════════════════════════════════

func TestGetDownloadURL_Success(t *testing.T) {
	svc, repo, _, _ := setup()

	name := "document.pdf"
	repo.GetByIDFn = func(_ context.Context, id string) (*domain.Media, error) {
		return &domain.Media{
			ID: id, Bucket: "media", ObjectKey: "doc.pdf",
			OriginalName: &name, Status: domain.StatusReady, IsPrivate: false,
		}, nil
	}

	resp, err := svc.GetDownloadURL(context.Background(), "id", userID1, "")

	require.NoError(t, err)
	assert.NotEmpty(t, resp.URL)
	assert.Equal(t, 3600, resp.ExpiresIn)
}

// ═══════════════════════════════════════════════════════════
//  LinkToEntity
// ═══════════════════════════════════════════════════════════

func TestLinkToEntity_Success(t *testing.T) {
	svc, repo, _, _ := setup()

	var captured *domain.MediaEntityLink
	repo.CreateEntityLinkFn = func(_ context.Context, link *domain.MediaEntityLink) error {
		captured = link
		return nil
	}

	err := svc.LinkToEntity(context.Background(), "media-123", "MESSAGE", "456")

	require.NoError(t, err)
	assert.Equal(t, "media-123", captured.MediaID)
	assert.Equal(t, "MESSAGE", captured.EntityType)
	assert.Equal(t, "456", captured.EntityID)
	assert.NotEmpty(t, captured.ID)
}

// ═══════════════════════════════════════════════════════════
//  Delete
// ═══════════════════════════════════════════════════════════

func TestDelete_Success(t *testing.T) {
	svc, repo, storage, _ := setup()

	repo.GetByIDFn = func(_ context.Context, id string) (*domain.Media, error) {
		return &domain.Media{
			ID: id, OwnerUserID: userID1, Bucket: "media",
			ObjectKey: "photo.jpg", Status: domain.StatusReady,
		}, nil
	}

	err := svc.Delete(context.Background(), "id", userID1)

	require.NoError(t, err)
	assert.Equal(t, 1, repo.SoftDeleteCalls)
	assert.Equal(t, 1, storage.DeleteObjectCalls)
}

func TestDelete_NotOwner(t *testing.T) {
	svc, repo, storage, _ := setup()

	repo.GetByIDFn = func(_ context.Context, id string) (*domain.Media, error) {
		return &domain.Media{ID: id, OwnerUserID: userID1, Status: domain.StatusReady}, nil
	}

	err := svc.Delete(context.Background(), "id", userID2)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
	assert.Equal(t, 0, repo.SoftDeleteCalls)
	assert.Equal(t, 0, storage.DeleteObjectCalls)
}

func TestDelete_S3FailureDoesNotBlock(t *testing.T) {
	svc, repo, storage, _ := setup()

	repo.GetByIDFn = func(_ context.Context, id string) (*domain.Media, error) {
		return &domain.Media{
			ID: id, OwnerUserID: userID1, Bucket: "media",
			ObjectKey: "photo.jpg", Status: domain.StatusReady,
		}, nil
	}

	storage.DeleteObjectFn = func(_ context.Context, _, _ string) error {
		return fmt.Errorf("S3 timeout")
	}

	err := svc.Delete(context.Background(), "id", userID1)

	require.NoError(t, err)
	assert.Equal(t, 1, repo.SoftDeleteCalls)
}