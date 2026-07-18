package storage

import (
	"context"
	"errors"
	"io"
)

// S3Storage is a stub implementation of Storage backed by an S3-compatible API.
//
// This is an interface skeleton only — it does NOT add the AWS SDK dependency.
// Methods return ErrNotImplemented until the SDK is wired in. The struct fields
// document the configuration shape that the real implementation will need.
//
// TODO(fase-12): integrate the AWS SDK for Go (v2) and implement each method:
//   - Upload  -> s3.PutObject with the provided content type
//   - Download-> s3.GetObject, returning Body as the ReadCloser
//   - Delete  -> s3.DeleteObject (treat NoSuchKey as success)
// Do NOT add the dependency in this stub; wire it in a follow-up change.
type S3Storage struct {
	// Bucket is the target S3 bucket name.
	Bucket string
	// Region is the AWS region the bucket lives in.
	Region string
	// Endpoint overrides the default AWS S3 endpoint (for MinIO, R2, etc.).
	Endpoint string
	// ForcePathStyle enables path-style addressing (required by some MinIO setups).
	ForcePathStyle bool
}

// NewS3Storage returns an S3Storage stub. It performs no network calls and does
// not validate credentials.
func NewS3Storage(bucket, region, endpoint string, forcePathStyle bool) (*S3Storage, error) {
	if bucket == "" {
		return nil, errors.New("storage: s3 bucket is required")
	}
	return &S3Storage{
		Bucket:         bucket,
		Region:         region,
		Endpoint:       endpoint,
		ForcePathStyle: forcePathStyle,
	}, nil
}

// ErrNotImplemented is returned by S3Storage methods until the AWS SDK is wired.
var ErrNotImplemented = errors.New("storage: s3 backend not implemented yet")

// Upload is a stub. See the TODO on S3Storage.
func (s *S3Storage) Upload(ctx context.Context, key string, reader io.Reader, contentType string) error {
	return ErrNotImplemented
}

// Download is a stub. See the TODO on S3Storage.
func (s *S3Storage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	return nil, ErrNotImplemented
}

// Delete is a stub. See the TODO on S3Storage.
func (s *S3Storage) Delete(ctx context.Context, key string) error {
	return ErrNotImplemented
}
