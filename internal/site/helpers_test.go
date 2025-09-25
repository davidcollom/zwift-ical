package site

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureDir_Table(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		wantErr bool
	}{
		{
			name:    "happy path - create dir",
			dir:     filepath.Join(os.TempDir(), "testdir"),
			wantErr: false,
		},
		{
			name:    "sad path - invalid dir",
			dir:     string([]byte{0}), // invalid path
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(tt.dir)
			err := EnsureDir(tt.dir)
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
