package ical

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

			require.NotEmpty(t, ical)
			assert.Contains(t, ical, tt.expect)

			// Additional checks can be added here based on the expected output
			if len(ical) > 0 {
				assert.Contains(t, ical, "VERSION:2.0")
				assert.Contains(t, ical, "PRODID:Zwift Calendar - by David Collom")
			}

		})
	}
}
