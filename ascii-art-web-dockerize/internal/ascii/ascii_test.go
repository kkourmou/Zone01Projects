package ascii

import (
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// Tests run in package directory, so we move up to root to find ascii-style/
	os.Chdir("../../")
	os.Exit(m.Run())
}

func TestGenerateAscii(t *testing.T) {
	// Let's run the project tests mapped directly from auditors_tests.md
	tests := []struct {
		name     string
		text     string
		banner   string
		expected string
	}{
		{
			name:   "Standard 123",
			text:   "123??",
			banner: "standard",
			expected: `                     ___    ___  
 _   ____    _____  |__ \  |__ \ 
/ | |___ \  |___ /     ) |    ) |
| |   __) |   |_ \    / /    / / 
| |  / __/   ___) |  |_|    |_|  
|_| |_____| |____/   (_)    (_)  
                                 
                                 
`,
		},
		{
			name:   "Shadow Chars",
			text:   `$% "=`,
			banner: "shadow",
			expected: `                        _|  _|  
  _|   _|_|    _|       _|  _|  
_|_|_| _|_|  _|                _|_|_|_|_| 
_|_|       _|                   
  _|_|   _|  _|_|              _|_|_|_|_| 
_|_|_| _|    _|_|               
  _|                            
                                
`,
		},
		{
			name:   "Thinkertoy Alnum",
			text:   "123 T/fs#R",
			banner: "thinkertoy",
			expected: `                                                                      
  0    --  o-o        o-O-o     o  o-o      | |  o--o                 
 /|   o  o    |         |      /   |       -O-O- |   |                
o |     /   oo          |     o   -O-  o-o  | |  O-Oo                 
  |    /      |         |    /     |    \  -O-O- |  \                 
o-o-o o--o o-o          o   o      o   o-o  | |  o   o                
                                                                      
                                                                      
`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GenerateAscii(tc.text, tc.banner)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			// The engine relies on precise trailing whitespaces per character.
			// We trim trailing whitespaces per line to avoid markdown hardcoding discrepancies.
			resultTrimmed := trimTrailingSpaces(result)
			expectedTrimmed := trimTrailingSpaces(tc.expected)

			if resultTrimmed != expectedTrimmed {
				t.Errorf("\n=== Expected (len %d) ===\n%q\n=== Got (len %d) ===\n%q", len(expectedTrimmed), expectedTrimmed, len(resultTrimmed), resultTrimmed)
			}
		})
	}
}

func trimTrailingSpaces(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t\r")
	}
	return strings.Join(lines, "\n")
}
