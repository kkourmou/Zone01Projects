package args

import (
	"errors"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantOutput string
		wantText   string
		wantBanner string
		wantErr    bool
	}{
		{
			name:       "single string",
			args:       []string{"hello"},
			wantText:   "hello",
			wantBanner: "standard",
		},
		{
			name:       "string with standard banner",
			args:       []string{"hello", "standard"},
			wantText:   "hello",
			wantBanner: "standard",
		},
		{
			name:       "string with shadow banner",
			args:       []string{"hello", "shadow"},
			wantText:   "hello",
			wantBanner: "shadow",
		},
		{
			name:       "string with thinkertoy banner",
			args:       []string{"hello", "thinkertoy"},
			wantText:   "hello",
			wantBanner: "thinkertoy",
		},
		{
			name:       "output flag with string",
			args:       []string{"--output=test.txt", "hello"},
			wantOutput: "test.txt",
			wantText:   "hello",
			wantBanner: "standard",
		},
		{
			name:       "output flag with string and banner",
			args:       []string{"--output=test.txt", "hello", "shadow"},
			wantOutput: "test.txt",
			wantText:   "hello",
			wantBanner: "shadow",
		},
		{
			name:    "no args",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "output space separated is invalid",
			args:    []string{"--output", "test.txt", "hello"},
			wantErr: true,
		},
		{
			name:    "empty output value",
			args:    []string{"--output=", "hello"},
			wantErr: true,
		},
		{
			name:    "output flag only no string",
			args:    []string{"--output=test.txt"},
			wantErr: true,
		},
		{
			name:    "invalid banner",
			args:    []string{"hello", "unknown"},
			wantErr: true,
		},
		{
			name:    "unknown flag",
			args:    []string{"--unknown=x", "hello"},
			wantErr: true,
		},
		{
			name:    "too many args",
			args:    []string{"--output=test.txt", "hello", "standard", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, gotText, gotBanner, err := ParseArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if !errors.Is(err, ErrUsage) {
					t.Errorf("expected ErrUsage, got %v", err)
				}
				return
			}
			if gotOutput != tt.wantOutput {
				t.Errorf("output = %q, want %q", gotOutput, tt.wantOutput)
			}
			if gotText != tt.wantText {
				t.Errorf("text = %q, want %q", gotText, tt.wantText)
			}
			if gotBanner != tt.wantBanner {
				t.Errorf("banner = %q, want %q", gotBanner, tt.wantBanner)
			}
		})
	}
}
