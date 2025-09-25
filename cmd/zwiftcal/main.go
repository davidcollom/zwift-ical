package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/davidcollom/zwift-ical/internal/events"
	"github.com/davidcollom/zwift-ical/internal/ical"
	"github.com/davidcollom/zwift-ical/internal/site"
)

var publicDir string = filepath.Join("public")

func main() {
	if err := Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// Run executes the main logic and returns error for testability.
func Run() error {
	// Clean up public/ directory before generating new files
	log.Println("Cleaning up public/ directory...")
	if err := os.RemoveAll(publicDir); err != nil {
		return err
	}
	if err := site.EnsureDir(publicDir); err != nil {
		return err
	}

	log.Println("Fetching events from Zwift API...")
	eventsList, err := events.FetchEvents(200, "")
	if err != nil {
		return err
	}
	log.Printf("Fetched %d events", len(eventsList))

	sportMap := groupEventsBySport(eventsList)

	if err := ensureSportDirs(sportMap); err != nil {
		return err
	}

	if err := generateICalFiles(sportMap); err != nil {
		return err
	}

	icalData := ical.EventsToICal(eventsList)
	if err := site.GenerateSite(eventsList, icalData, publicDir); err != nil {
		return err
	}

	icalLinks := collectICalLinks(sportMap)
	indexOutput := filepath.Join(publicDir, "index.html")
	if err := site.RenderIndexLinks(icalLinks, indexOutput); err != nil {
		return err
	}

	log.Println("All iCal files and static site generated in public/")
	return nil
}

// groupEventsBySport groups events by their sport.
func groupEventsBySport(eventsList []events.Event) map[string][]events.Event {
	sportMap := map[string][]events.Event{}
	for _, e := range eventsList {
		sport := strings.ToLower(e.Sport)
		if sport == "" {
			continue
		}
		sportMap[sport] = append(sportMap[sport], e)
	}
	return sportMap
}

// ensureSportDirs creates necessary directories for each sport.
func ensureSportDirs(sportMap map[string][]events.Event) error {
	for sport := range sportMap {
		for _, sub := range []string{"rides", "workouts", "races", "tag"} {
			dir := filepath.Join(publicDir, sport, sub)
			if err := site.EnsureDir(dir); err != nil {
				log.Printf("Error creating directory %s: %v", dir, err)
				return err
			}
		}
	}
	return nil
}

// generateICalFiles generates iCal files for each sport/type/tag.
func generateICalFiles(sportMap map[string][]events.Event) error {
	for sport, sportEvents := range sportMap {
		log.Printf("Processing sport: %s (%d events)", sport, len(sportEvents))
		types := map[string]string{
			"rides":    "GROUP_RIDE",
			"workouts": "GROUP_WORKOUT",
			"races":    "RACE",
		}
		for typ, eventType := range types {
			filtered := filterEventsByType(sportEvents, eventType)
			icalData := ical.EventsToICal(filtered)
			icalPath := filepath.Join(publicDir, sport, typ, "events.ics")
			if err := site.WriteICal(icalData, icalPath); err != nil {
				log.Printf("Error writing %s: %v", icalPath, err)
			} else {
				log.Printf("Generated %s (%d events)", icalPath, len(filtered))
			}
		}
		// Tags
		tagSet := collectTags(sportEvents)
		for tag := range tagSet {
			tagEvents := filterEventsByTag(sportEvents, tag)
			icalData := ical.EventsToICal(tagEvents)
			icalPath := filepath.Join(publicDir, sport, "tag", tag+".ics")
			if err := site.WriteICal(icalData, icalPath); err != nil {
				log.Printf("Error writing %s: %v", icalPath, err)
			} else {
				log.Printf("Generated %s (%d events)", icalPath, len(tagEvents))
			}
		}
	}
	return nil
}

// filterEventsByType filters events by event type.
func filterEventsByType(eventsList []events.Event, eventType string) []events.Event {
	filtered := []events.Event{}
	for _, e := range eventsList {
		if e.EventType == eventType {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// collectTags collects unique tags from events.
func collectTags(eventsList []events.Event) map[string]struct{} {
	tagSet := map[string]struct{}{}
	for _, e := range eventsList {
		for _, tag := range e.Tags {
			if tag != "" && !strings.Contains(tag, "=") {
				tagSet[tag] = struct{}{}
			}
		}
	}
	return tagSet
}

// filterEventsByTag filters events by a specific tag.
func filterEventsByTag(eventsList []events.Event, tag string) []events.Event {
	tagEvents := []events.Event{}
	for _, e := range eventsList {
		for _, t := range e.Tags {
			if t == tag {
				tagEvents = append(tagEvents, e)
				break
			}
		}
	}
	return tagEvents
}

// collectICalLinks collects all iCal file paths for index rendering.
func collectICalLinks(sportMap map[string][]events.Event) []string {
	var icalLinks []string
	for sport, sportEvents := range sportMap {
		sportLower := strings.ToLower(sport)
		types := map[string]string{
			"rides":    "GROUP_RIDE",
			"workouts": "GROUP_WORKOUT",
			"races":    "RACE",
		}
		for typ := range types {
			icalPath := filepath.Join(sportLower, typ, "events.ics")
			icalLinks = append(icalLinks, icalPath)
		}
		tagSet := collectTags(sportEvents)
		for tag := range tagSet {
			tagPath := filepath.Join(sportLower, "tag", tag+".ics")
			icalLinks = append(icalLinks, tagPath)
		}
	}
	return icalLinks
}
