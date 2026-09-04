package parser

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseBanner(t *testing.T) {
	mockContent := ""
	for i := 32; i <= 126; i++ {
		mockContent += "\n"
		for j := 0; j < 8; j++ {
			mockContent += "MockLine\n"
		}
	}

	tmpFile := filepath.Join(t.TempDir(), "mock_banner.txt")
	err := os.WriteFile(tmpFile, []byte(mockContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create mock banner file: %v", err)
	}

	font, err := ParseBanner(tmpFile)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(font) != 95 {
		t.Errorf("Expected exactly 95 characters parsed, got %d", len(font))
	}

	spaceLines := font[' ']
	if len(spaceLines) != 8 {
		t.Errorf("Expected 8 lines for space character, got %d", len(spaceLines))
	}

	expected := []string{"MockLine", "MockLine", "MockLine", "MockLine", "MockLine", "MockLine", "MockLine", "MockLine"}
	if !reflect.DeepEqual(spaceLines, expected) {
		t.Errorf("Parsed lines do not match expected correctly.")
	}

	_, err = ParseBanner("non_existent_file.txt")
	if err == nil {
		t.Errorf("Expected error for missing file, got nil")
	}

	badContent := "\nLine1\nLine2\n"
	badFile := filepath.Join(t.TempDir(), "bad_banner.txt")
	os.WriteFile(badFile, []byte(badContent), 0o644)

	_, err = ParseBanner(badFile)
	if err == nil {
		t.Errorf("Expected error for invalid file format, got nil")
	}
}
