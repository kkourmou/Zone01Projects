package parser

import (
	"bufio"
	"fmt"
	"os"
)

// ParseBanner reads an ASCII art banner file and returns a font map.
// The file must contain 95 characters (ASCII 32-126), each represented by 8 lines.
func ParseBanner(filepath string) (map[rune][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open banner file: %w", err)
	}
	defer file.Close()

	font := make(map[rune][]string)
	scanner := bufio.NewScanner(file)

	currentChar := rune(32)
	var currentLines []string
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		if lineCount == 1 {
			continue
		}

		currentLines = append(currentLines, line)

		if len(currentLines) == 8 {
			font[currentChar] = currentLines
			currentChar++
			currentLines = nil
			lineCount = 0
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading banner file: %w", err)
	}

	if len(font) != 95 {
		return nil, fmt.Errorf("banner file parsing error: expected 95 characters, parsed %d (file has invalid format)", len(font))
	}

	return font, nil
}
