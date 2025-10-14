package ical

import (
	"fmt"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/davidcollom/zwift-ical/internal/events"
)

func EventsToICal(events []events.Event) string {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	cal.SetProductId("Zwift Calendar - by David Collom")
	for _, e := range events {
		ev := cal.AddEvent(fmt.Sprintf("%d", e.ID))
		ev.SetSummary(fmt.Sprintf("[%s] %s", worldName(e), e.Name))
		ev.SetLocation(worldName(e))
		ev.SetDescription(e.Description)
		// Set start and end times
		// start, _ := time.Parse(time.ISO8601, e.EventStart)
		ev.SetStartAt(e.EventStart.UTC())
		end := e.EventStart.Add(time.Duration(e.DurationInSeconds) * time.Second)
		ev.SetEndAt(end.UTC())
		// Add URL/Image if available
		if e.ImageUrl != "" {
			ev.SetURL(e.ImageUrl)
		}
		ev.SetClass(ics.ClassificationPublic)
		ev.SetLastModifiedAt(time.Now())
	}
	return cal.Serialize()
}
