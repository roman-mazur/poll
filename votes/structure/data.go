package structure

import "time"

type Vote struct {
	TalkName  string    `json:"talk_name"`
	Timestamp time.Time `json:"timestamp"`
	// To distinguish unique users.
	VoterId string `json:"voter_id"`
	// From 1 to 10 (From "I don't know what you are talking about" to "totally support").
	Value uint8 `json:"value" cue:">=1 & <=10"`
}

type Label struct {
	TalkName  string    `json:"talk_name"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}
