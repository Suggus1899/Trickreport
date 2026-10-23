package ticket

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"encoding/xml"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"net/http"
	"unicode/utf16"
	"unicode/utf8"
)

func validateAttachmentContent(declared string, data []byte) (string, error) {
	contentType, _, err := mime.ParseMediaType(declared)
	if err != nil || !allowedContentTypes[contentType] || len(data) == 0 {
		return "", ErrInvalidContent
	}
	detected, _, _ := mime.ParseMediaType(http.DetectContentType(data))
	valid := false
	switch contentType {
	case "image/png", "image/jpeg", "image/gif":
		cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
		valid = err == nil && cfg.Width > 0 && cfg.Height > 0 && int64(cfg.Width)*int64(cfg.Height) <= 25000000 && detected == contentType
		if valid {
			_, _, err = image.Decode(bytes.NewReader(data))
			valid = err == nil
		}
	case "image/webp":
		valid = validWebP(data)
	case "image/bmp":
		valid = validBMP(data)
	case "application/pdf":
		valid = detected == contentType && bytes.HasPrefix(data, []byte("%PDF-")) && bytes.HasSuffix(bytes.TrimSpace(data), []byte("%%EOF"))
	case "text/plain", "text/csv":
		valid = detected == "text/plain" && utf8.Valid(data) && !bytes.ContainsRune(data, 0)
		if valid && contentType == "text/csv" {
			reader := csv.NewReader(bytes.NewReader(data))
			reader.FieldsPerRecord = -1
			_, err := reader.ReadAll()
			valid = err == nil
		}
	case "application/msword", "application/vnd.ms-excel":
		valid = validCompoundDocument(data, contentType)
	default:
		valid = validZipDocument(data, contentType)
	}
	if !valid {
		return "", ErrInvalidContent
	}
	return contentType, nil
}

func validZipDocument(data []byte, contentType string) bool {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return false
	}
	part := ""
	root := ""
	switch contentType {
	case "application/zip":
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		part, root = "word/document.xml", "document"
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		part, root = "xl/workbook.xml", "workbook"
	case "application/vnd.openxmlformats-officedocument.presentationml.presentation":
		part, root = "ppt/presentation.xml", "presentation"
	default:
		return false
	}
	found := map[string]bool{}
	var total uint64
	for _, file := range reader.File {
		if file.Flags&1 != 0 {
			return false
		}
		total += file.UncompressedSize64
		// Bound decompression while checking container integrity and CRCs.
		if total > 100<<20 {
			return false
		}
		r, err := file.Open()
		if err != nil {
			return false
		}
		if file.Name == part || file.Name == "[Content_Types].xml" || file.Name == "_rels/.rels" {
			decoder := xml.NewDecoder(io.LimitReader(r, 100<<20))
			first := ""
			for {
				token, err := decoder.Token()
				if err == io.EOF {
					break
				}
				if err != nil {
					_ = r.Close()
					return false
				}
				if start, ok := token.(xml.StartElement); ok && first == "" {
					first = start.Name.Local
				}
			}
			want := root
			if file.Name == "[Content_Types].xml" {
				want = "Types"
			}
			if file.Name == "_rels/.rels" {
				want = "Relationships"
			}
			if first != want {
				_ = r.Close()
				return false
			}
			found[file.Name] = true
		} else {
			if _, err := io.Copy(io.Discard, io.LimitReader(r, 100<<20)); err != nil {
				_ = r.Close()
				return false
			}
		}
		if err := r.Close(); err != nil {
			return false
		}
	}
	return part == "" || (found[part] && found["[Content_Types].xml"] && found["_rels/.rels"])
}

func validWebP(data []byte) bool {
	if len(data) < 20 || string(data[:4]) != "RIFF" || string(data[8:12]) != "WEBP" || uint64(binary.LittleEndian.Uint32(data[4:8]))+8 != uint64(len(data)) {
		return false
	}
	found := false
	for offset := 12; offset < len(data); {
		if len(data)-offset < 8 {
			return false
		}
		size := int(binary.LittleEndian.Uint32(data[offset+4 : offset+8]))
		end := offset + 8 + size
		if end > len(data) {
			return false
		}
		chunk := string(data[offset : offset+4])
		payload := data[offset+8 : end]
		switch chunk {
		case "VP8 ":
			if len(payload) < 10 || !bytes.Equal(payload[3:6], []byte{0x9d, 0x01, 0x2a}) {
				return false
			}
			found = true
		case "VP8L":
			if len(payload) < 5 || payload[0] != 0x2f {
				return false
			}
			found = true
		case "VP8X":
			if len(payload) != 10 {
				return false
			}
		}
		offset = end + (size & 1)
		if offset > len(data) {
			return false
		}
	}
	return found
}

func validBMP(data []byte) bool {
	if len(data) < 54 || string(data[:2]) != "BM" || uint64(binary.LittleEndian.Uint32(data[2:6])) != uint64(len(data)) {
		return false
	}
	header := uint64(binary.LittleEndian.Uint32(data[14:18]))
	offset := uint64(binary.LittleEndian.Uint32(data[10:14]))
	width := int32(binary.LittleEndian.Uint32(data[18:22]))
	height := int32(binary.LittleEndian.Uint32(data[22:26]))
	if header < 40 || 14+header > offset || offset >= uint64(len(data)) || width <= 0 || height == 0 || height == -2147483648 || binary.LittleEndian.Uint16(data[26:28]) != 1 {
		return false
	}
	if height < 0 {
		height = -height
	}
	bits := uint64(binary.LittleEndian.Uint16(data[28:30]))
	switch bits {
	case 1, 4, 8, 16, 24, 32:
	default:
		return false
	}
	// shortcut: accept uncompressed BMP only; add a decoder if compressed BMP is required.
	if binary.LittleEndian.Uint32(data[30:34]) != 0 {
		return false
	}
	return offset+((uint64(width)*bits+31)/32)*4*uint64(height) <= uint64(len(data))
}

func validCompoundDocument(data []byte, contentType string) bool {
	if len(data) < 512 || !bytes.Equal(data[:8], []byte{0xd0, 0xcf, 0x11, 0xe0, 0xa1, 0xb1, 0x1a, 0xe1}) || binary.LittleEndian.Uint16(data[28:30]) != 0xfffe {
		return false
	}
	shift := binary.LittleEndian.Uint16(data[30:32])
	if shift != 9 && shift != 12 {
		return false
	}
	sectorSize := uint64(1) << shift
	if uint64(len(data))%sectorSize != 0 {
		return false
	}
	sectorCount := uint64(len(data))/sectorSize - 1
	// shortcut: inline DIFAT covers ordinary legacy Office uploads; add extended DIFAT support when needed.
	fatCount := uint64(binary.LittleEndian.Uint32(data[44:48]))
	if fatCount == 0 || fatCount > 109 {
		return false
	}
	fat := make([]uint32, 0, fatCount*sectorSize/4)
	for i := uint64(0); i < fatCount; i++ {
		sector := uint64(binary.LittleEndian.Uint32(data[76+i*4 : 80+i*4]))
		if sector >= sectorCount {
			return false
		}
		start := (sector + 1) * sectorSize
		for n := uint64(0); n < sectorSize; n += 4 {
			fat = append(fat, binary.LittleEndian.Uint32(data[start+n:start+n+4]))
		}
	}
	sector := binary.LittleEndian.Uint32(data[48:52])
	seen := map[uint32]bool{}
	foundRoot, foundStream := false, false
	for sector != 0xfffffffe {
		if uint64(sector) >= sectorCount || int(sector) >= len(fat) || seen[sector] {
			return false
		}
		seen[sector] = true
		start := (uint64(sector) + 1) * sectorSize
		for n := uint64(0); n < sectorSize; n += 128 {
			entry := data[start+n : start+n+128]
			length := int(binary.LittleEndian.Uint16(entry[64:66]))
			if entry[66] == 0 {
				continue
			}
			if length < 2 || length > 64 || length%2 != 0 {
				return false
			}
			chars := make([]uint16, (length-2)/2)
			for i := range chars {
				chars[i] = binary.LittleEndian.Uint16(entry[i*2 : i*2+2])
			}
			name := string(utf16.Decode(chars))
			if entry[66] == 5 && name == "Root Entry" {
				foundRoot = true
			}
			if entry[66] == 2 && ((contentType == "application/msword" && name == "WordDocument") || (contentType == "application/vnd.ms-excel" && (name == "Workbook" || name == "Book"))) {
				foundStream = true
			}
		}
		sector = fat[sector]
	}
	return foundRoot && foundStream
}
