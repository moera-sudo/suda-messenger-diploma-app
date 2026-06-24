package s3

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// s3Region is the SigV4 signing region. MinIO is single-region; setting it lets
// the SDK presign offline instead of calling GetBucketLocation over the network.
const s3Region = "us-east-1"

type FileStorage interface {
	EnsureBucket(ctx context.Context, bucket string) error
	// presignHost (optional) overrides the host baked into the presigned URL, so
	// the link points at the same origin the client used to reach the API. Empty →
	// the configured S3_PUBLIC_ENDPOINT.
	PresignPutURL(ctx context.Context, presignHost, bucket, objectKey, contentType string, expiry time.Duration) (string, error)
	PresignGetURL(ctx context.Context, presignHost, bucket, objectKey string, expiry time.Duration) (string, error)
	PresignGetDownloadURL(ctx context.Context, presignHost, bucket, objectKey, filename string, expiry time.Duration) (string, error)
	HeadObject(ctx context.Context, bucket, objectKey string) (size int64, etag string, err error)
	DeleteObject(ctx context.Context, bucket, objectKey string) error
	GetFile(ctx context.Context, bucket, objectKey string) (io.ReadCloser, string, error)
}

type s3Storage struct {
	client        *minio.Client
	presignClient *minio.Client
	endpoint      string
	useSSL        bool
	accessKey     string
	secretKey     string

	// presignClients caches per-host clients so a presigned URL can be signed for
	// the exact origin the request came in on (X-Forwarded-Host).
	mu             sync.Mutex
	presignClients map[string]*minio.Client
}


func NewS3Storage(endpoint, accessKey, secretKey string, useSSL bool, publicEndpoint string) (FileStorage, error) {
	// Основной клиент — для внутренних операций (HeadObject, DeleteObject, EnsureBucket)
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
		Region: s3Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	// Клиент для presigned URL — подписывает с публичным хостом
	presignEndpoint := publicEndpoint
	if presignEndpoint == "" {
		presignEndpoint = endpoint
	}

	// Region must be set so presigning is fully offline: without it minio-go does a
	// GetBucketLocation (".../?location=") network call against the public endpoint,
	// which is unreachable from inside the container and hangs until timeout.
	presignClient, err := minio.New(presignEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
		Region: s3Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create presign minio client: %w", err)
	}

	return &s3Storage{
		client:         client,
		presignClient:  presignClient,
		endpoint:       endpoint,
		useSSL:         useSSL,
		accessKey:      accessKey,
		secretKey:      secretKey,
		presignClients: make(map[string]*minio.Client),
	}, nil
}

// presignClientFor returns a minio client whose presigned URLs use `host`
// (e.g. "10.0.2.2:8080"). Empty host → the default S3_PUBLIC_ENDPOINT client.
// Clients are cached per host. Falls back to the default client on error.
func (s *s3Storage) presignClientFor(host string) *minio.Client {
	if host == "" {
		return s.presignClient
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if c, ok := s.presignClients[host]; ok {
		return c
	}
	c, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(s.accessKey, s.secretKey, ""),
		Secure: s.useSSL,
		Region: s3Region,
	})
	if err != nil {
		return s.presignClient
	}
	s.presignClients[host] = c
	return c
}

func (s *s3Storage) EnsureBucket(ctx context.Context, bucket string) error {
	exists, err := s.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("check bucket exists: %w", err)
	}
	if exists {
		return nil
	}

	if err := s.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("create bucket: %w", err)
	}

	// Все бакеты приватные — доступ только через presigned URL
	return nil
}

// PresignPutURL — клиент загружает файл напрямую в S3 по этой ссылке
func (s *s3Storage) PresignPutURL(ctx context.Context, presignHost, bucket, objectKey, contentType string, expiry time.Duration) (string, error) {
	presignedURL, err := s.presignClientFor(presignHost).PresignedPutObject(ctx, bucket, objectKey, expiry)
	if err != nil {
		return "", fmt.Errorf("presign put: %w", err)
	}
	return presignedURL.String(), nil
}

// PresignGetURL — для просмотра в чате/браузере (Content-Disposition: inline)
func (s *s3Storage) PresignGetURL(ctx context.Context, presignHost, bucket, objectKey string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "inline")

	presignedURL, err := s.presignClientFor(presignHost).PresignedGetObject(ctx, bucket, objectKey, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("presign get: %w", err)
	}
	return presignedURL.String(), nil
}

// PresignGetDownloadURL — для скачивания файла (Content-Disposition: attachment)
func (s *s3Storage) PresignGetDownloadURL(ctx context.Context, presignHost, bucket, objectKey, filename string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)
	disposition := fmt.Sprintf("attachment; filename=\"%s\"", filename)
	reqParams.Set("response-content-disposition", disposition)

	presignedURL, err := s.presignClientFor(presignHost).PresignedGetObject(ctx, bucket, objectKey, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("presign get download: %w", err)
	}
	return presignedURL.String(), nil
}

// HeadObject — проверяет что файл реально существует в S3 (после upload)
func (s *s3Storage) HeadObject(ctx context.Context, bucket, objectKey string) (int64, string, error) {
	stat, err := s.client.StatObject(ctx, bucket, objectKey, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.StatusCode == http.StatusNotFound {
			return 0, "", fmt.Errorf("object not found: %s/%s", bucket, objectKey)
		}
		return 0, "", fmt.Errorf("stat object: %w", err)
	}
	return stat.Size, stat.ETag, nil
}

func (s *s3Storage) DeleteObject(ctx context.Context, bucket, objectKey string) error {
	return s.client.RemoveObject(ctx, bucket, objectKey, minio.RemoveObjectOptions{})
}

// GetFile — fallback для proxy-режима (если presigned URL не подходит)
func (s *s3Storage) GetFile(ctx context.Context, bucket, objectKey string) (io.ReadCloser, string, error) {
	object, err := s.client.GetObject(ctx, bucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", err
	}

	stat, err := object.Stat()
	if err != nil {
		return nil, "", err
	}

	return object, stat.ContentType, nil
}

// func (s *s3Storage) replaceHost(presignedURL string) string {
// 	if s.publicEndpoint == "" || s.publicEndpoint == s.endpoint {
// 		return presignedURL
// 	}
// 	return strings.Replace(presignedURL, s.endpoint, s.publicEndpoint, 1)
// }