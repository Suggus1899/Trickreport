package storage

import "testing"

// TestValidateContentType_SVGRejected is a regression test: SVG must never
// be an accepted upload type because it can embed scripts (stored XSS).
// SVG content sniffs as text/plain via http.DetectContentType, so both the
// actual-type and declared-type checks must reject it.
func TestValidateContentType_SVGRejected(t *testing.T) {
	cases := []struct {
		name     string
		actual   string
		declared string
	}{
		{"svg declared as svg", "image/svg+xml", "image/svg+xml"},
		{"svg sniffed as text", "text/plain; charset=utf-8", "image/svg+xml"},
		{"svg declared with parameters", "text/plain", "image/svg+xml; charset=utf-8"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if ValidateContentType(tc.actual, tc.declared) {
				t.Errorf("ValidateContentType(%q, %q) = true, want false (SVG must be rejected)", tc.actual, tc.declared)
			}
		})
	}
	if IsAllowedContentType("image/svg+xml") {
		t.Error("IsAllowedContentType(\"image/svg+xml\") = true, want false")
	}
}

func TestValidateContentType_AllowedTypesStillAccepted(t *testing.T) {
	cases := []struct {
		name     string
		actual   string
		declared string
	}{
		{"png", "image/png", "image/png"},
		{"jpeg with parameters", "image/jpeg", "image/jpeg; charset=binary"},
		{"pdf", "application/pdf", "application/pdf"},
		{"plain text", "text/plain; charset=utf-8", "text/plain"},
		{"zip bundle", "application/zip", "application/zip"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if !ValidateContentType(tc.actual, tc.declared) {
				t.Errorf("ValidateContentType(%q, %q) = false, want true", tc.actual, tc.declared)
			}
		})
	}
}

func TestValidateContentType_UnknownTypesRejected(t *testing.T) {
	cases := []struct {
		name     string
		actual   string
		declared string
	}{
		{"html", "text/html", "text/html"},
		{"executable", "application/x-msdownload", "application/x-msdownload"},
		{"octet stream", "application/octet-stream", "application/octet-stream"},
		{"allowed actual but forbidden declared", "image/png", "text/html"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if ValidateContentType(tc.actual, tc.declared) {
				t.Errorf("ValidateContentType(%q, %q) = true, want false", tc.actual, tc.declared)
			}
		})
	}
}
