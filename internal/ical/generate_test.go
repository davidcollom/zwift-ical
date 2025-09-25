package ical

import (
	"testing"
	"time"

	"github.com/davidcollom/zwift-ical/internal/events"
)

func TestEventsToICal_Table(t *testing.T) {
	tests := []struct {
		name   string
		events []events.Event
		expect string
	}{
		{
			name: "happy path - one event",
			events: []events.Event{{
				ID: 1, Name: "Test Ride", Description: "Desc",
				EventStart:        time.Now().Format(time.RFC3339),
				DurationInSeconds: 3600, MapId: 1,
			}},
			expect: "BEGIN:VCALENDAR",
		},
		{
			name:   "sad path - empty events",
			events: []events.Event{},
			expect: "BEGIN:VCALENDAR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ical := EventsToICal(tt.events)
			if len(ical) == 0 || ical[:15] != tt.expect {
				t.Errorf("iCal output invalid: %s", ical)
			}
		})
	}
}
