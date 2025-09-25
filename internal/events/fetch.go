package events

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

var (
	zwiftAPI   = "https://us-or-rly101.zwift.com/api/public/events/upcoming"
	MaxRetries = 10
	MinSleep   = 100 * time.Millisecond
	MaxSleep   = 2 * time.Second
)

func FetchEvents(limit int, tag string) ([]Event, error) {
	req, err := retryablehttp.NewRequest("GET", zwiftAPI, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("limit", fmt.Sprintf("%d", limit))
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
