package internal

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"strings"

	"rsc.io/pdf"
)

func ExtractTextFromFile(f multipart.File, ext string) (string, error) {
	switch ext {
	case ".txt", ".md", ".csv":
		content, err := io.ReadAll(f)
		if err != nil {
			return "", err
		}
		return string(content), nil

	case ".pdf":
		return extractTextFromPDF(f)

	default:
		return "", errors.New("Unsupported file extension")
	}
}

func extractTextFromPDF(f multipart.File) (string, error) {
	buf, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(buf)
	doc, err := pdf.NewReader(reader, int64(len(buf)))
	if err != nil {
		return "", err
	}

	var text strings.Builder
	numPages := doc.NumPage()
	for i := 1; i <= numPages; i++ {
		page := doc.Page(i)
		if page.V.IsNull() {
			continue
		}
		content := page.Content()
		for _, txt := range content.Text {
			text.WriteString(txt.S)
			text.WriteString(" ")
		}
	}

	return text.String(), nil
}
