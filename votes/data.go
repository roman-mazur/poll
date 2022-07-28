package votes

import (
	"time"
)

type Vote struct {
	TalkName  string `cue:"=~\"^.{1,50}$\""`
	Timestamp time.Time
	// To distinguish unique users.
	VoterId string
	// From 0 to 10 (From "I don't know what you are talking about" to "totally support").
	Value uint8 `cue:">=0 & <=10"`
}

type Label struct {
	TalkName  string
	Name      string
	Timestamp time.Time
}
