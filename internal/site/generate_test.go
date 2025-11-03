package site

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteICal_Table(t *testing.T) {
	tests := []struct {
		name    string
		content string
		path    string
		wantErr bool
	}{
		{
			name:    "happy path - valid file",
			content: "BEGIN:VCALENDAR\nEND:VCALENDAR",
			path:    filepath.Join(os.TempDir(), "test.ics"),
			wantErr: false,
		},
		{
			name:    "sad path - invalid file",
			content: "BEGIN:VCALENDAR\nEND:VCALENDAR",
			path:    string([]byte{0}),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(tt.path)
			err := WriteICal(tt.content, tt.path)
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestRenderIndexLinks(t *testing.T) {
	tmp := filepath.Join(os.TempDir(), "index.html")
	defer os.Remove(tmp)
	paths := []string{"cycling/rides/events.ics", "running/races/events.ics"}
	err := RenderIndexLinks(paths, tmp)
	if err != nil {
		t.Fatalf("RenderIndexLinks failed: %v", err)
	}
	data, err := os.ReadFile(tmp)
	if err != nil || len(data) == 0 {
		t.Errorf("Index file not created or empty")
	}
}

func TestGenerateRedirects(t *testing.T) {
	tmp := filepath.Join(os.TempDir(), "_redirects")
	defer os.Remove(tmp)
	paths := []string{"cycling/rides/events.ics", "running/races/events.ics"}
	err := GenerateRedirects(paths, tmp)
	if err != nil {
		t.Fatalf("GenerateRedirects failed: %v", err)
	}
	data, err := os.ReadFile(tmp)
	if err != nil || len(data) == 0 {
		t.Errorf("Redirects file not created or empty")
	}
}
