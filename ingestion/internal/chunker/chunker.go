package chunker

import "unicode/utf8"

func Chunk(text string, size int) []string {
	var chunks []string

	for len(text) > 0 {
		if utf8.RuneCountInString(text) <= size {
			chunks = append(chunks, text)
			break
		}

		chunk := string([]rune(text)[:size])
		chunks = append(chunks, chunk)
		text = string([]rune(text)[size:])
	}
	return chunks
}
