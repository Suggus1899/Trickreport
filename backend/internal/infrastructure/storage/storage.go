// Package storage provides an abstraction over object/blob backends used for
// ticket attachments and other large binary payloads.
//
// The Storage interface is intentionally minimal (Upload/Download/Delete) so
// that the application layer depends on a stable contract regardless of whether
// files are persisted to a local disk, an S3-compatible API, or any other
// backend.
//
// TODO(fase-12): wire a Storage implementation into the dependency graph in
// internal/interfaces/http/wire.go and into the AttachmentService so that file
// bytes are streamed through Storage instead of being held in memory. Do NOT
// edit wire.go or attachment_service.go as part of this stub — wire it in a
// follow-up change.
package storage

import (
	"context"
	"io"
)

// Storage is the port for object/blob backends.
//
// Keys are opaque, backend-relative identifiers. Implementations MUST be safe
// for concurrent use. Callers are responsible for closing any ReadCloser
// returned by Download.
type Storage interface {
	// Upload stores the contents of reader under the given key with the
	// provided content type. The reader is consumed fully; callers do not need
	// to close it.
	Upload(ctx context.Context, key string, reader io.Reader, contentType string) error

	// Download returns a ReadCloser for the object stored under key. The caller
	// MUST close the returned reader when done.
	Download(ctx context.Context, key string) (io.ReadCloser, error)

	// Delete removes the object stored under key. A missing object is not an
	// error (idempotent).
	Delete(ctx context.Context, key string) error
}
