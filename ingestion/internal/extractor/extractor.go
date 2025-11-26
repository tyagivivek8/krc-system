package extractor

import (
	"bytes"
	"errors"
	"strings"

	"github.com/ledongthuc/pdf"
)

func ExtractText(data []byte, filename string) (string, error) {
	if strings.HasSuffix(strings.ToLower(filename), ".txt") {
		return string(data), nil
	}

	if strings.HasSuffix(strings.ToLower(filename), ".pdf") {
		return extractPDF(data)
	}

	return "", errors.New("Unsupported file type")
}

func extractPDF(data []byte) (string, error) {
	reader := bytes.NewReader(data)

	pdfReader, err := pdf.NewReader(reader, int64(len(data)))
	if err != nil {
		return "", err
	}

	var text strings.Builder

	numPages := pdfReader.NumPage()
	for pageIndex := 1; pageIndex <= numPages; pageIndex++ {
		page := pdfReader.Page(pageIndex)
		if page.V.IsNull() {
			continue
		}

		content, err := page.GetPlainText(nil)
		if err != nil {
			return "", err
		}

		text.WriteString(content)
		text.WriteString("\n")
	}

	return text.String(), nil
}
