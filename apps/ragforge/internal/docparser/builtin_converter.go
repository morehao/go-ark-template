package docparser

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
)

var simpleFormats = map[string]bool{
	"md": true, "markdown": true,
	"txt": true, "text": true,
	"csv":  true,
	"json": true,
}

var imageFormats = map[string]bool{
	"jpg": true, "jpeg": true, "png": true, "gif": true,
	"bmp": true, "tiff": true, "webp": true,
}

var audioFormats = map[string]bool{
	"mp3": true, "wav": true, "m4a": true, "flac": true, "ogg": true,
}

func init() {
	for k := range imageFormats {
		simpleFormats[k] = true
	}
	for k := range audioFormats {
		simpleFormats[k] = true
	}
}

func IsSimpleFormat(fileType string) bool {
	return simpleFormats[strings.ToLower(strings.TrimPrefix(fileType, "."))]
}

type SimpleFormatReader struct{}

func (b *SimpleFormatReader) Read(_ context.Context, req *ReadRequest) (*ReadResult, error) {
	ft := strings.ToLower(strings.TrimPrefix(req.FileType, "."))
	if ft == "" {
		ft = strings.TrimPrefix(strings.ToLower(filepath.Ext(req.FileName)), ".")
	}

	switch {
	case ft == "md" || ft == "markdown":
		return &ReadResult{MarkdownContent: string(req.FileContent)}, nil
	case ft == "txt" || ft == "text":
		return &ReadResult{MarkdownContent: string(req.FileContent)}, nil
	case ft == "csv":
		md, err := csvToMarkdown(req.FileContent)
		if err != nil {
			return nil, fmt.Errorf("csv conversion failed: %w", err)
		}
		return &ReadResult{MarkdownContent: md}, nil
	case ft == "json":
		md, err := jsonToMarkdown(req.FileContent)
		if err != nil {
			return nil, fmt.Errorf("json conversion failed: %w", err)
		}
		return &ReadResult{MarkdownContent: md}, nil
	case imageFormats[ft]:
		return imageToResult(req.FileName, req.FileContent), nil
	case audioFormats[ft]:
		return audioToResult(req.FileName, req.FileContent), nil
	default:
		return nil, fmt.Errorf("unsupported simple format: %s", ft)
	}
}

func imageToResult(fileName string, data []byte) *ReadResult {
	if fileName == "" {
		fileName = "image.png"
	}
	refPath := "images/" + fileName
	safeRef := strings.ReplaceAll(refPath, " ", "%20")
	mime := http.DetectContentType(data)

	return &ReadResult{
		MarkdownContent: fmt.Sprintf("![%s](%s)", fileName, safeRef),
		ImageRefs: []ImageRef{
			{
				Filename:    fileName,
				OriginalRef: safeRef,
				MimeType:    mime,
				ImageData:   data,
				IsOriginal:  true,
			},
		},
	}
}

func IsImageFormat(fileType string) bool {
	return imageFormats[strings.ToLower(strings.TrimPrefix(fileType, "."))]
}

func IsAudioFormat(fileType string) bool {
	return audioFormats[strings.ToLower(strings.TrimPrefix(fileType, "."))]
}

func audioToResult(fileName string, data []byte) *ReadResult {
	if fileName == "" {
		fileName = "audio.mp3"
	}
	return &ReadResult{
		MarkdownContent: fmt.Sprintf("[Audio file: %s]", fileName),
		IsAudio:         true,
		AudioData:       data,
	}
}

func ensureOriginalImageRef(req *ReadRequest, mdContent string, imageRefs []ImageRef) (string, []ImageRef) {
	ft := strings.ToLower(strings.TrimPrefix(req.FileType, "."))
	if ft == "" {
		ft = strings.TrimPrefix(strings.ToLower(filepath.Ext(req.FileName)), ".")
	}
	if !imageFormats[ft] {
		return mdContent, imageRefs
	}
	if len(req.FileContent) == 0 {
		return mdContent, imageRefs
	}

	fileName := req.FileName
	if fileName == "" {
		fileName = "image." + ft
	}
	refPath := "images/" + fileName

	if strings.Contains(mdContent, refPath) {
		return mdContent, imageRefs
	}

	imgLine := fmt.Sprintf("![%s](%s)", fileName, refPath)
	if strings.TrimSpace(mdContent) == "" {
		mdContent = imgLine
	} else {
		mdContent = imgLine + "\n\n" + mdContent
	}

	mime := http.DetectContentType(req.FileContent)
	imageRefs = append(imageRefs, ImageRef{
		Filename:    fileName,
		OriginalRef: refPath,
		MimeType:    mime,
		ImageData:   req.FileContent,
		IsOriginal:  true,
	})

	return mdContent, imageRefs
}

func csvToMarkdown(data []byte) (string, error) {
	reader := csv.NewReader(strings.NewReader(string(data)))
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}
	if len(records) == 0 {
		return "", nil
	}

	var sb strings.Builder

	header := records[0]
	sb.WriteString("| ")
	sb.WriteString(strings.Join(header, " | "))
	sb.WriteString(" |\n")

	sb.WriteString("|")
	for range header {
		sb.WriteString(" --- |")
	}
	sb.WriteString("\n")

	for _, row := range records[1:] {
		sb.WriteString("| ")
		cells := make([]string, len(header))
		for i := range cells {
			if i < len(row) {
				cells[i] = row[i]
			}
		}
		sb.WriteString(strings.Join(cells, " | "))
		sb.WriteString(" |\n")
	}

	return sb.String(), nil
}
