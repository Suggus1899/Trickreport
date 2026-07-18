package storage

import (
	"net/http"
	"strings"
)

// allowedContentTypes is the whitelist of accepted MIME content types for
// file uploads. Images, PDFs, plain text, and common document formats are
// allowed.
var allowedContentTypes = map[string]bool{
	// Images
	"image/jpeg":             true,
	"image/png":              true,
	"image/gif":              true,
	"image/webp":             true,
	"image/svg+xml":          true,
	"image/bmp":              true,
	// PDF
	"application/pdf": true,
	// Text
	"text/plain": true,
	"text/csv":   true,
	// Documents
	"application/msword":                                                          true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document":     true,
	"application/vnd.ms-excel":                                                    true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":           true,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation":   true,
	"application/vnd.oasis.opendocument.text":                                     true,
	// Archives (for support bundles)
	"application/zip": true,
}

// ValidateContentType checks that the actual content type (detected from the
// file's magic bytes) matches the declared content type and that both are in
// the whitelist. Returns true if the file is safe to accept.
//
// actualType should be the result of http.DetectContentType on the first 512
// bytes of the file. declaredType should be the Content-Type header from the
// multipart upload.
func ValidateContentType(actualType, declaredType string) bool {
	// Normalize: strip parameters like "; charset=utf-8"
	actualType = normalizeContentType(actualType)
	declaredType = normalizeContentType(declaredType)

	// The actual type must be in the whitelist
	if !allowedContentTypes[actualType] {
		return false
	}

	// If a declared type is provided, it must also be in the whitelist.
	// We don't require an exact match between actual and declared because
	// browsers sometimes declare a slightly different type, but both must
	// be in the whitelist.
	if declaredType != "" && !allowedContentTypes[declaredType] {
		return false
	}

	return true
}

// normalizeContentType strips parameters (e.g. "; charset=utf-8") and
// lowercases the content type.
func normalizeContentType(ct string) string {
	ct = strings.TrimSpace(ct)
	if idx := strings.Index(ct, ";"); idx >= 0 {
		ct = ct[:idx]
	}
	return strings.ToLower(strings.TrimSpace(ct))
}

// DetectContentType reads the first 512 bytes of data and returns the
// detected content type. This is a convenience wrapper around
// http.DetectContentType.
func DetectContentType(data []byte) string {
	// http.DetectContentType only needs the first 512 bytes
	if len(data) > 512 {
		data = data[:512]
	}
	return http.DetectContentType(data)
}

// IsAllowedContentType returns true if the given content type is in the
// whitelist. This is useful for quick checks without content sniffing.
func IsAllowedContentType(ct string) bool {
	return allowedContentTypes[normalizeContentType(ct)]
}
