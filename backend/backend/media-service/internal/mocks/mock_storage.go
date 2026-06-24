package mocks

import (
	"context"
	"io"
	"time"
)

type MockStorage struct {
	EnsureBucketFn          func(ctx context.Context, bucket string) error
	PresignPutURLFn         func(ctx context.Context, bucket, objectKey, contentType string, expiry time.Duration) (string, error)
	PresignGetURLFn         func(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error)
	PresignGetDownloadURLFn func(ctx context.Context, bucket, objectKey, filename string, expiry time.Duration) (string, error)
	HeadObjectFn            func(ctx context.Context, bucket, objectKey string) (int64, string, error)
	DeleteObjectFn          func(ctx context.Context, bucket, objectKey string) error
	GetFileFn               func(ctx context.Context, bucket, objectKey string) (io.ReadCloser, string, error)

	EnsureBucketCalls  int
	PresignPutURLCalls int
	HeadObjectCalls    int
	DeleteObjectCalls  int
}

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

func (m *MockStorage) EnsureBucket(ctx context.Context, bucket string) error {
	m.EnsureBucketCalls++
	if m.EnsureBucketFn != nil {
		return m.EnsureBucketFn(ctx, bucket)
	}
	return nil
}

func (m *MockStorage) PresignPutURL(ctx context.Context, presignHost, bucket, objectKey, contentType string, expiry time.Duration) (string, error) {
	m.PresignPutURLCalls++
	if m.PresignPutURLFn != nil {
		return m.PresignPutURLFn(ctx, bucket, objectKey, contentType, expiry)
	}
	return "https://s3.example.com/presigned-put-url", nil
}

func (m *MockStorage) PresignGetURL(ctx context.Context, presignHost, bucket, objectKey string, expiry time.Duration) (string, error) {
	if m.PresignGetURLFn != nil {
		return m.PresignGetURLFn(ctx, bucket, objectKey, expiry)
	}
	return "https://s3.example.com/presigned-get-url?disposition=inline", nil
}

func (m *MockStorage) PresignGetDownloadURL(ctx context.Context, presignHost, bucket, objectKey, filename string, expiry time.Duration) (string, error) {
	if m.PresignGetDownloadURLFn != nil {
		return m.PresignGetDownloadURLFn(ctx, bucket, objectKey, filename, expiry)
	}
	return "https://s3.example.com/presigned-download-url?disposition=attachment", nil
}

func (m *MockStorage) HeadObject(ctx context.Context, bucket, objectKey string) (int64, string, error) {
	m.HeadObjectCalls++
	if m.HeadObjectFn != nil {
		return m.HeadObjectFn(ctx, bucket, objectKey)
	}
	return 10240, "etag-abc123", nil
}

func (m *MockStorage) DeleteObject(ctx context.Context, bucket, objectKey string) error {
	m.DeleteObjectCalls++
	if m.DeleteObjectFn != nil {
		return m.DeleteObjectFn(ctx, bucket, objectKey)
	}
	return nil
}

func (m *MockStorage) GetFile(ctx context.Context, bucket, objectKey string) (io.ReadCloser, string, error) {
	if m.GetFileFn != nil {
		return m.GetFileFn(ctx, bucket, objectKey)
	}
	return nil, "", nil
}