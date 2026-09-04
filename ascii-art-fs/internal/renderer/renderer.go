package renderer

import (
	"fmt"
	"strings"
)

// RenderText converts input text to ASCII art using the provided font map.
func RenderText(input string, font map[rune][]string) (string, error) {
	if input == "" {
		return "", nil
	}

	input = strings.ReplaceAll(input, "\r\n", "\\n")
	input = strings.ReplaceAll(input, "\n", "\\n")

	parts := strings.Split(input, "\\n")

	allEmpty := true
	for _, p := range parts {
		if p != "" {
			allEmpty = false
			break
		}
	}

	var result strings.Builder

	if allEmpty {
		for i := 0; i < len(parts)-1; i++ {
			result.WriteString("\n")
		}
		return result.String(), nil
	}

	for _, part := range parts {
		if part == "" {
			result.WriteString("\n")
			continue
		}
		for i := 0; i < 8; i++ {
			for _, ch := range part {
				lines, ok := font[ch]
				if !ok {
					return "", fmt.Errorf("character %q not supported", ch)
				}
				result.WriteString(lines[i])
			}
			result.WriteString("\n")
		}
	}

	return result.String(), nil
}
