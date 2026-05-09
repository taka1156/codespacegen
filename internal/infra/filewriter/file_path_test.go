package filewriter

import (
	"path/filepath"
	"testing"
)

func TestResolveOutputPath_Success(t *testing.T) {
	tests := []struct {
		name         string
		outputDir    string
		relativePath string
		want         string
	}{
		{
			name:         "simple file",
			outputDir:    "/out",
			relativePath: "file.txt",
			want:         "/out/file.txt",
		},
		{
			name:         "nested file",
			outputDir:    "/out",
			relativePath: "sub/file.txt",
			want:         "/out/sub/file.txt",
		},
		{
			name:         "redundant dot-slash is cleaned",
			outputDir:    "/out",
			relativePath: "./file.txt",
			want:         "/out/file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveOutputPath(tt.outputDir, tt.relativePath)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != filepath.FromSlash(tt.want) {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolveOutputPath_Error(t *testing.T) {
	tests := []struct {
		name         string
		outputDir    string
		relativePath string
		wantContains string
	}{
		{
			name:         "path escapes with ..",
			outputDir:    "/out",
			relativePath: "../secret.txt",
			wantContains: "path escapes output directory",
		},
		{
			name:         "deep escape",
			outputDir:    "/out",
			relativePath: "sub/../../secret.txt",
			wantContains: "path escapes output directory",
		},
		{
			name:         "absolute path rejected",
			outputDir:    "/out",
			relativePath: "/etc/passwd",
			wantContains: "absolute path is not allowed",
		},
		{
			name:         "dot-only path rejected",
			outputDir:    "/out",
			relativePath: ".",
			wantContains: "invalid file path",
		},
		{
			name:         "empty path rejected",
			outputDir:    "/out",
			relativePath: "",
			wantContains: "invalid file path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ResolveOutputPath(tt.outputDir, tt.relativePath)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !contains(err.Error(), tt.wantContains) {
				t.Errorf("error %q does not contain %q", err.Error(), tt.wantContains)
			}
		})
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
