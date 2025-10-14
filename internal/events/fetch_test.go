package events

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			name: "happy path - valid response",
			serverResp: `[
				{"id":1,"name":"Test Ride","description":"Desc","eventStart":"2025-10-14T17:45:00.000+0000","durationInSeconds":3600,"imageUrl":"https://example.com/img.png","mapId":1,"sport":"CYCLING","eventType":"GROUP_RIDE","tags":["test"]},
				{"id":2,"name":"Test Running","description":"Running Event","eventStart":"2025-10-14T17:45:00.000+0000","durationInSeconds":3600,"imageUrl":"https://example.com/img.png","mapId":1,"sport":"RUNNING","eventType":"WORKOUT","tags":["workout"]}
				]`,
			statusCode: 200,
			wantErr:    false,
			wantCount:  2,
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
		{
			name: "timezone offset handling - different offsets",
			serverResp: `[
				{"id":3,"name":"IST Event","description":"India Standard Time","eventStart":"2025-10-14T23:15:00.000+0530","durationInSeconds":3600,"imageUrl":"https://example.com/img.png","mapId":1,"sport":"CYCLING","eventType":"GROUP_RIDE","tags":["test"]},
				{"id":4,"name":"PST Event","description":"Pacific Standard Time","eventStart":"2025-10-14T09:45:00.000-0800","durationInSeconds":3600,"imageUrl":"https://example.com/img.png","mapId":1,"sport":"RUNNING","eventType":"WORKOUT","tags":["workout"]}
				]`,
			statusCode: 200,
			wantErr:    false,
			wantCount:  2,
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

			events, err := FetchEvents(2, "")
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Len(t, events, tt.wantCount)

			if len(events) > 0 && tt.name == "happy path - valid response" {
				cycling := events[0]
				running := events[1]

				assert.Equal(t, 1, cycling.ID)
				assert.Equal(t, "Test Ride", cycling.Name)
				assert.Equal(t, "Desc", cycling.Description)
				// assert.Equal(t, "2024-09-15T10:00:00Z", cycling.EventStart.Format(time.RFC3339))
				assert.Equal(t, 3600, cycling.DurationInSeconds)
				assert.Equal(t, "https://example.com/img.png", cycling.ImageUrl)
				assert.Equal(t, 1, cycling.MapId)
				assert.Equal(t, "CYCLING", cycling.Sport)
				assert.Equal(t, "GROUP_RIDE", cycling.EventType)
				assert.Contains(t, cycling.Tags, "test")

				assert.Equal(t, 2, running.ID)
				assert.Equal(t, "Test Running", running.Name)
				assert.Equal(t, "Running Event", running.Description)
				assert.NotEmpty(t, running.EventStart)
				// assert.Equal(t, "2025-10-14T17:45:00.000+0000", running.EventStart.Format(time.RFC3339Nano))
				assert.Equal(t, 3600, running.DurationInSeconds)
				assert.Equal(t, "https://example.com/img.png", running.ImageUrl)
				assert.Equal(t, 1, running.MapId)
				assert.Equal(t, "RUNNING", running.Sport)
				assert.Equal(t, "WORKOUT", running.EventType)
				assert.Contains(t, running.Tags, "workout")
			}

			if len(events) > 0 && tt.name == "timezone offset handling - different offsets" {
				istEvent := events[0]
				pstEvent := events[1]

				assert.Equal(t, 3, istEvent.ID)
				assert.Equal(t, "IST Event", istEvent.Name)
				assert.Equal(t, "India Standard Time", istEvent.Description)
				assert.NotEmpty(t, istEvent.EventStart)
				assert.Equal(t, 3600, istEvent.DurationInSeconds)

				assert.Equal(t, 4, pstEvent.ID)
				assert.Equal(t, "PST Event", pstEvent.Name)
				assert.Equal(t, "Pacific Standard Time", pstEvent.Description)
				assert.NotEmpty(t, pstEvent.EventStart)
				assert.Equal(t, 3600, pstEvent.DurationInSeconds)
			}
		})
	}
}
