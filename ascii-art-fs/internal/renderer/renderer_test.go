package renderer

import (
	"strings"
	"testing"
)

func makeFont(chars string, lines [8]string) map[rune][]string {
	font := map[rune][]string{}
	for _, ch := range chars {
		font[ch] = lines[:]
	}
	return font
}

func TestRenderText_Empty(t *testing.T) {
	result, err := RenderText("", nil)
	if err != nil {
		t.Fatal(err)
	}
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestRenderText_Newline(t *testing.T) {
	font := makeFont("a", [8]string{"1", "2", "3", "4", "5", "6", "7", "8"})
	result, err := RenderText("a\\nb", font)
	if err == nil && !strings.Contains(result, "\n") {
		t.Error("expected newline in result")
	}
}

func TestRenderText_UnsupportedChar(t *testing.T) {
	font := map[rune][]string{}
	_, err := RenderText("a", font)
	if err == nil {
		t.Error("expected error for unsupported character")
	}
}

func TestRenderText_OnlyNewlines(t *testing.T) {
	result, err := RenderText("\\n\\n", nil)
	if err != nil {
		t.Fatal(err)
	}
	if result != "\n\n" {
		t.Errorf("expected two newlines, got %q", result)
	}
}
