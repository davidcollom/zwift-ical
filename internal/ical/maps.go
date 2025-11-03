package ical

import "github.com/davidcollom/zwift-ical/internal/events"

var mapIDs = []string{
	"Watopia",
	"Richmond",
	"London",
	"New York",
	"Innsbuck",
	"Bologna TT",
	"Yorkshire",
	"Crit City",
	"Makuri Islands",
	"France",
	"Paris",
	"Gravel Mountain",
	"Scotland",
}

func worldName(e events.Event) string {
	mapID := e.MapId
	if mapID > 0 && mapID <= len(mapIDs) {
		return mapIDs[mapID-1]
	}
	return "Unknown"
}
