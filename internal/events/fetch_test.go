package events

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchEvents_Table(t *testing.T) {
	MaxRetries = 1 // Only Fail once during testing...
	tests := []struct {
		name       string
		serverResp string
		statusCode int
		wantErr    bool
		wantCount  int
	}{
		{
			name:       "happy path - valid response",
			serverResp: `[{"id":1,"name":"Test Ride","description":"Desc","eventStart":"2024-09-15T10:00:00Z","durationInSeconds":3600,"imageUrl":"https://example.com/img.png","mapId":1,"sport":"CYCLING","eventType":"GROUP_RIDE","tags":["test"]}]`,
			statusCode: 200,
			wantErr:    false,
			wantCount:  1,
		},
		{
			name:       "sad path - server error",
			serverResp: `Internal Server Error`,
			statusCode: 500,
			wantErr:    true,
			wantCount:  0,
		},
		{
			name:       "sad path - invalid JSON",
			serverResp: `not json`,
			statusCode: 200,
			wantErr:    true,
			wantCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.serverResp))
			}))
			defer ts.Close()

			zwiftAPIOrig := zwiftAPI
			zwiftAPI = ts.URL
			defer func() { zwiftAPI = zwiftAPIOrig }()

			events, err := FetchEvents(1, "")
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if len(events) != tt.wantCount {
				t.Errorf("expected %d events, got %d", tt.wantCount, len(events))
			}
		})
	}
}
