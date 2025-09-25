package events


type Event struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	EventStart        string   `json:"eventStart"`
	DurationInSeconds int      `json:"durationInSeconds"`
	ImageUrl          string   `json:"imageUrl"`
	MapId             int      `json:"mapId"`
	Sport             string   `json:"sport"`
	EventType         string   `json:"eventType"`
	Tags              []string `json:"tags"`
	// Ignore unknown fields during JSON unmarshalling
}
