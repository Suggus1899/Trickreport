package ticket

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"github.com/google/uuid"
	"image"
	"image/png"
	"testing"
)

func pngFixture(t *testing.T) []byte {
	t.Helper()
	var data bytes.Buffer
	if err := png.Encode(&data, image.NewRGBA(image.Rect(0, 0, 1, 1))); err != nil {
		t.Fatal(err)
	}
	return data.Bytes()
}

func TestAttachment_ContentMatchesDeclaration(t *testing.T) {
	pngBytes := pngFixture(t)
	var zipBytes bytes.Buffer
	zw := zip.NewWriter(&zipBytes)
	part, _ := zw.Create("hello.txt")
	_, _ = part.Write([]byte("hello"))
	_ = zw.Close()
	creator := uuid.New()
	tk := ownedTicket(creator)
	for _, test := range []struct {
		name, mime string
		data       []byte
		valid      bool
	}{
		{"png", "image/png", pngBytes, true}, {"MIME parameters", "image/png; charset=binary", pngBytes, true},
		{"pdf", "application/pdf", []byte("%PDF-1.4\n1 0 obj<</Type/Catalog>>endobj\n%%EOF\n"), true},
		{"forged image", "image/png", []byte("text content"), false}, {"empty", "text/plain", nil, false},
		{"forged PDF", "application/pdf", []byte("plain text"), false}, {"truncated PDF", "application/pdf", []byte("%PDF-1.4"), false},
		{"zip", "application/zip", zipBytes.Bytes(), true}, {"not a docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", zipBytes.Bytes(), false},
		{"forged OLE", "application/msword", []byte("not a binary office document"), false},
		{"svg with script", "image/svg+xml", []byte(`<svg xmlns="http://www.w3.org/2000/svg"><script>alert(1)</script></svg>`), false},
		{"svg declared as png", "image/png", []byte(`<svg xmlns="http://www.w3.org/2000/svg"></svg>`), false},
	} {
		t.Run(test.name, func(t *testing.T) {
			repo := &mockAttachmentRepo{}
			svc := newAttachmentSvc(repo, &mockRepo{ticket: tk})
			_, err := svc.Upload(context.Background(), UploadInput{TenantID: tk.TenantID, TicketID: tk.ID, UserID: creator, Role: "end_user", Filename: "file", ContentType: test.mime, FileData: test.data})
			if test.valid && err != nil {
				t.Fatal(err)
			}
			if !test.valid && !errors.Is(err, ErrInvalidContent) {
				t.Fatalf("error=%v want invalid content", err)
			}
			if !test.valid && repo.createCalls != 0 {
				t.Error("invalid bytes persisted")
			}
			if test.name == "MIME parameters" && repo.gotInput.ContentType != "image/png" {
				t.Errorf("stored MIME=%s", repo.gotInput.ContentType)
			}
		})
	}
}
