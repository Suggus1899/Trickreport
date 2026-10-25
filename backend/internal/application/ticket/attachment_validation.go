package ticket

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"encoding/xml"
	"image"
	// Decoders for the image formats checked by content sniffing; imported for side effects.
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
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
	if !validContentBody(contentType, detected, data) {
		return "", ErrInvalidContent
	}
	return contentType, nil
}

// validContentBody reports whether data matches its declared content type.
// detected is the type sniffed from the bytes.
func validContentBody(contentType, detected string, data []byte) bool {
	switch contentType {
	case "image/png", "image/jpeg", "image/gif":
		return validRaster(data, detected == contentType)
	case "image/webp":
		return validWebP(data)
	case "image/bmp":
		return validBMP(data)
	case "application/pdf":
		return detected == contentType && bytes.HasPrefix(data, []byte("%PDF-")) && bytes.HasSuffix(bytes.TrimSpace(data), []byte("%%EOF"))
	case "text/plain", "text/csv":
		return validText(data, contentType, detected)
	case "application/msword", "application/vnd.ms-excel":
		return validCompoundDocument(data, contentType)
	default:
		return validZipDocument(data, contentType)
	}
}

// validRaster reports whether data decodes as an image of an allowed size.
func validRaster(data []byte, typeMatches bool) bool {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || !typeMatches || cfg.Width <= 0 || cfg.Height <= 0 || int64(cfg.Width)*int64(cfg.Height) > 25000000 {
		return false
	}
	_, _, err = image.Decode(bytes.NewReader(data))
	return err == nil
}

// validText reports whether data is NUL-free UTF-8 text. CSV content must also parse as CSV.
func validText(data []byte, contentType, detected string) bool {
	if detected != "text/plain" || !utf8.Valid(data) || bytes.ContainsRune(data, 0) {
		return false
	}
	if contentType != "text/csv" {
		return true
	}
	reader := csv.NewReader(bytes.NewReader(data))
	reader.FieldsPerRecord = -1
	_, err := reader.ReadAll()
	return err == nil
}

func validZipDocument(data []byte, contentType string) bool {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return false
	}
	part, root, ok := zipDocumentPart(contentType)
	if !ok {
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
		if !readZipEntry(file, part, root, found) {
			return false
		}
	}
	return part == "" || (found[part] && found["[Content_Types].xml"] && found["_rels/.rels"])
}

// zipDocumentPart returns the main part and its required XML root element for a
// ZIP-based document type. ok is false for types that are not ZIP containers.
func zipDocumentPart(contentType string) (part, root string, ok bool) {
	switch contentType {
	case "application/zip":
		return "", "", true
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return "word/document.xml", "document", true
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return "xl/workbook.xml", "workbook", true
	case "application/vnd.openxmlformats-officedocument.presentationml.presentation":
		return "ppt/presentation.xml", "presentation", true
	default:
		return "", "", false
	}
}

// readZipEntry reads an entry to the end, which verifies its CRC. The entries that define
// the document must be well-formed XML with the expected root element; those are recorded in found.
func readZipEntry(file *zip.File, part, root string, found map[string]bool) bool {
	r, err := file.Open()
	if err != nil {
		return false
	}
	var ok bool
	if file.Name == part || file.Name == "[Content_Types].xml" || file.Name == "_rels/.rels" {
		ok = xmlRootIs(io.LimitReader(r, 100<<20), requiredRoot(file.Name, root))
		if ok {
			found[file.Name] = true
		}
	} else {
		_, copyErr := io.Copy(io.Discard, io.LimitReader(r, 100<<20))
		ok = copyErr == nil
	}
	closeErr := r.Close()
	return ok && closeErr == nil
}

// xmlRootIs reports whether r is well-formed XML whose first element is named root.
func xmlRootIs(r io.Reader, root string) bool {
	decoder := xml.NewDecoder(r)
	first := ""
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false
		}
		if start, ok := token.(xml.StartElement); ok && first == "" {
			first = start.Name.Local
		}
	}
	return first == root
}

// requiredRoot returns the XML root element that a ZIP entry must have.
func requiredRoot(name, documentRoot string) string {
	switch name {
	case "[Content_Types].xml":
		return "Types"
	case "_rels/.rels":
		return "Relationships"
	default:
		return documentRoot
	}
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
	width := binary.LittleEndian.Uint32(data[18:22])
	height := binary.LittleEndian.Uint32(data[22:26])
	if header < 40 || 14+header > offset || offset >= uint64(len(data)) || width == 0 || width > math.MaxInt32 || height == 0 || height == 1<<31 || binary.LittleEndian.Uint16(data[26:28]) != 1 {
		return false
	}
	if height > math.MaxInt32 {
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
	fat, ok := compoundFAT(data, sectorSize, sectorCount)
	if !ok {
		return false
	}
	return compoundDirectoryHasDocument(data, fat, sectorSize, sectorCount, contentType)
}

// compoundFAT reads the file allocation table from the sectors listed in the header.
func compoundFAT(data []byte, sectorSize, sectorCount uint64) ([]uint32, bool) {
	// shortcut: inline DIFAT covers ordinary legacy Office uploads; add extended DIFAT support when needed.
	fatCount := uint64(binary.LittleEndian.Uint32(data[44:48]))
	if fatCount == 0 || fatCount > 109 {
		return nil, false
	}
	fat := make([]uint32, 0, fatCount*sectorSize/4)
	for i := uint64(0); i < fatCount; i++ {
		sector := uint64(binary.LittleEndian.Uint32(data[76+i*4 : 80+i*4]))
		if sector >= sectorCount {
			return nil, false
		}
		start := (sector + 1) * sectorSize
		for n := uint64(0); n < sectorSize; n += 4 {
			fat = append(fat, binary.LittleEndian.Uint32(data[start+n:start+n+4]))
		}
	}
	return fat, true
}

// compoundDirectoryHasDocument walks the directory sector chain. It reports whether the
// directory holds a root entry and the main stream of contentType, and fails on malformed chains.
func compoundDirectoryHasDocument(data []byte, fat []uint32, sectorSize, sectorCount uint64, contentType string) bool {
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
			if entry[66] == 0 {
				continue
			}
			name, ok := compoundEntryName(entry)
			if !ok {
				return false
			}
			if entry[66] == 5 && name == "Root Entry" {
				foundRoot = true
			}
			if entry[66] == 2 && isDocumentStream(contentType, name) {
				foundStream = true
			}
		}
		sector = fat[sector]
	}
	return foundRoot && foundStream
}

// compoundEntryName decodes the UTF-16 name of a directory entry. ok is false when the name length is invalid.
func compoundEntryName(entry []byte) (string, bool) {
	length := int(binary.LittleEndian.Uint16(entry[64:66]))
	if length < 2 || length > 64 || length%2 != 0 {
		return "", false
	}
	chars := make([]uint16, (length-2)/2)
	for i := range chars {
		chars[i] = binary.LittleEndian.Uint16(entry[i*2 : i*2+2])
	}
	return string(utf16.Decode(chars)), true
}

// isDocumentStream reports whether name is the main stream of a legacy Office document of contentType.
func isDocumentStream(contentType, name string) bool {
	switch contentType {
	case "application/msword":
		return name == "WordDocument"
	case "application/vnd.ms-excel":
		return name == "Workbook" || name == "Book"
	default:
		return false
	}
}
