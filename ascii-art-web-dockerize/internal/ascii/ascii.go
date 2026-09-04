package ascii

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"ascii-art/internal/parser"
	"ascii-art/internal/renderer"
)

var (
	validBanners = map[string]bool{"standard": true, "shadow": true, "thinkertoy": true}
	bannerCache  = map[string]map[rune][]string{}
	cacheMu      sync.RWMutex
)

// IsValidBanner reports whether a banner name is one of the bundled fonts.
func IsValidBanner(banner string) bool {
	return validBanners[banner]
}

func GenerateAscii(text string, banner string) (string, error) {
	if text == "" {
		return "", nil
	}

	if !IsValidBanner(banner) {
		return "", fmt.Errorf("invalid banner: %q", banner)
	}

	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\\n", "\n") // Support literal '\n' characters as line breaks.

	for _, ch := range text {
		if (ch < ' ' || ch > '~') && ch != '\n' {
			return "", fmt.Errorf("unsupported character: %q", ch)
		}
	}

	font, err := loadBanner(banner)
	if err != nil {
		return "", fmt.Errorf("banner not found: %w", err)
	}

	return renderer.RenderText(text, font)
}

func loadBanner(banner string) (map[rune][]string, error) {
	cacheMu.RLock()
	font, ok := bannerCache[banner]
	cacheMu.RUnlock()
	if ok {
		return font, nil
	}

	cacheMu.Lock()
	defer cacheMu.Unlock()

	if font, ok = bannerCache[banner]; ok {
		return font, nil
	}

	font, err := parser.ParseBanner(projectFilePath(filepath.Join("ascii-style", banner+".txt")))
	if err != nil {
		return nil, err
	}
	bannerCache[banner] = font
	return font, nil
}

func projectFilePath(relativePath string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return relativePath
	}

	for {
		candidate := filepath.Join(cwd, relativePath)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}

		parent := filepath.Dir(cwd)
		if parent == cwd {
			return relativePath
		}
		cwd = parent
	}
}
