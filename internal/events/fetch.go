package events

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

var (
	zwiftAPI   = "https://us-or-rly101.zwift.com/api/public/events/upcoming"
	MaxRetries = 10
	MinSleep   = 100 * time.Millisecond
	MaxSleep   = 2 * time.Second
	MaxStart   = 10000
)

func FetchEvents(limit int, tag string) (events []Event, err error) {
	for start := 0; ; start += limit {
		var eventSet []Event
		eventSet, err = fetchEvents(limit, tag, start)
		log.Printf("Got %d Events", len(eventSet))
		if len(eventSet) == 0 || err != nil || start >= MaxStart {
			log.Printf("Breaking on iteration %d with %d events and error: %v", start, len(eventSet), err)
			break
		}
		events = append(events, eventSet...)
	}
	unique := make(map[int]struct{})
	filtered := make([]Event, 0, len(events))
	for _, e := range events {
		if _, exists := unique[e.ID]; !exists {
			unique[e.ID] = struct{}{}
			filtered = append(filtered, e)
		}
	}
	events = filtered
	return events, err
}

func fetchEvents(limit int, tag string, start int) ([]Event, error) {
	req, err := retryablehttp.NewRequest("GET", zwiftAPI, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("limit", fmt.Sprintf("%d", limit))
	if start > 0 {
		q.Add("start", fmt.Sprintf("%d", start))
	}
	if tag != "" {
		q.Add("tags", tag)
	}

	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")

	client := retryablehttp.NewClient()
	client.RetryMax = MaxRetries
	client.RetryWaitMin = MinSleep
	client.RetryWaitMax = MaxSleep

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var events []Event
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, err
	}
	return events, nil
}
