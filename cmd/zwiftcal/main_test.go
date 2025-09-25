package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/davidcollom/zwift-ical/internal/events"
)

func TestGroupEventsBySport(t *testing.T) {
	evs := []events.Event{
		{Sport: "Cycling"},
		{Sport: "Running"},
		{Sport: "Cycling"},
		{Sport: ""},
	}
	got := groupEventsBySport(evs)
	want := map[string][]events.Event{
		"cycling": {evs[0], evs[2]},
		"running": {evs[1]},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("groupEventsBySport() = %v, want %v", got, want)
	}
}

func TestEnsureSportDirs(t *testing.T) {
	tmp := filepath.Join(os.TempDir(), "zwiftcal_test")
	defer os.RemoveAll(tmp)
	sportMap := map[string][]events.Event{
		"cycling": {},
		"running": {},
	}
	publicDir = tmp
	err := ensureSportDirs(sportMap)
	if err != nil {
		t.Fatalf("ensureSportDirs() error = %v", err)
	}
	for sport := range sportMap {
		for _, sub := range []string{"rides", "workouts", "races", "tag"} {
			dir := filepath.Join(tmp, sport, sub)
			if _, err := os.Stat(dir); err != nil {
				t.Errorf("Directory %s not created", dir)
			}
		}
	}
}

func TestFilterEventsByType(t *testing.T) {
	evs := []events.Event{
		{EventType: "GROUP_RIDE"},
		{EventType: "RACE"},
		{EventType: "GROUP_WORKOUT"},
		{EventType: "GROUP_RIDE"},
	}
	got := filterEventsByType(evs, "GROUP_RIDE")
	if len(got) != 2 {
		t.Errorf("filterEventsByType() got %d, want 2", len(got))
	}
}

func TestCollectTags(t *testing.T) {
	evs := []events.Event{
		{Tags: []string{"A", "B", "foo=bar"}},
		{Tags: []string{"B", "C"}},
		{Tags: []string{""}},
	}
	got := collectTags(evs)
	want := map[string]struct{}{"A": {}, "B": {}, "C": {}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("collectTags() = %v, want %v", got, want)
	}
}

func TestFilterEventsByTag(t *testing.T) {
	evs := []events.Event{
		{Tags: []string{"A", "B"}},
		{Tags: []string{"B", "C"}},
		{Tags: []string{"D"}},
	}
	got := filterEventsByTag(evs, "B")
	if len(got) != 2 {
		t.Errorf("filterEventsByTag() got %d, want 2", len(got))
	}
}

func TestCollectICalLinks(t *testing.T) {
	evs := []events.Event{
		{Sport: "Cycling", EventType: "GROUP_RIDE", Tags: []string{"A"}},
		{Sport: "Cycling", EventType: "RACE", Tags: []string{"B"}},
		{Sport: "Running", EventType: "GROUP_WORKOUT", Tags: []string{"C"}},
	}
	sportMap := groupEventsBySport(evs)
	got := collectICalLinks(sportMap)
	// Should contain paths for rides, workouts, races, and tags for each sport
	wantContains := []string{
		filepath.Join("cycling", "rides", "events.ics"),
		filepath.Join("cycling", "races", "events.ics"),
		filepath.Join("cycling", "workouts", "events.ics"),
		filepath.Join("cycling", "tag", "A.ics"),
		filepath.Join("cycling", "tag", "B.ics"),
		filepath.Join("running", "workouts", "events.ics"),
		filepath.Join("running", "races", "events.ics"),
		filepath.Join("running", "rides", "events.ics"),
		filepath.Join("running", "tag", "C.ics"),
	}
	for _, want := range wantContains {
		found := false
		for _, gotPath := range got {
			if strings.HasSuffix(gotPath, want) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("collectICalLinks() missing %s", want)
		}
	}
}
