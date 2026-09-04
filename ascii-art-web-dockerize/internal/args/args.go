package args

import (
	"errors"
	"strings"
)

var ErrUsage = errors.New("Usage: go run . [OPTION] [STRING] [BANNER]\n\nEX: go run . --output=<fileName.txt> something standard")

var validBanners = map[string]bool{"standard": true, "shadow": true, "thinkertoy": true, "stef": true}

// ParseArgs extracts output file, text, and banner from command line arguments.
// Supported forms:
//
//	go run . "text"
//	go run . "text" banner
//	go run . --output=file.txt "text"
//	go run . --output=file.txt "text" banner
func ParseArgs(a []string) (output, text, banner string, err error) {
	if len(a) == 0 {
		return "", "", "", ErrUsage
	}

	i := 0

	// Parse optional --output= flag
	if strings.HasPrefix(a[i], "--output=") {
		output = strings.TrimPrefix(a[i], "--output=")
		if output == "" {
			return "", "", "", ErrUsage
		}
		i++
	} else if strings.HasPrefix(a[i], "--") {
		return "", "", "", ErrUsage
	}

	remaining := a[i:]

	switch len(remaining) {
	case 1:
		text = remaining[0]
		banner = "standard"
	case 2:
		if !validBanners[remaining[1]] {
			return "", "", "", ErrUsage
		}
		text = remaining[0]
		banner = remaining[1]
	default:
		return "", "", "", ErrUsage
	}

	return output, text, banner, nil
}
