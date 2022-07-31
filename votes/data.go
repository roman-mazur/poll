package votes

import (
	"time"
)

type Vote struct {
	TalkName  string    `json:"talk_name"`
	Timestamp time.Time `json:"timestamp"`
	// To distinguish unique users.
	VoterId string `json:"voter_id"`
	// From 1 to 10 (From "I don't know what you are talking about" to "totally support").
	Value uint8 `json:"value" cue:">=1 & <=10"`
}

func (v Vote) talkName() string {
	return v.TalkName
}

func (v Vote) touch() {
	v.Timestamp = time.Now()
}

type Label struct {
	TalkName  string    `json:"talk_name"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

func (l Label) talkName() string {
	return l.TalkName
}

func (l Label) touch() {
	l.Timestamp = time.Now()
}
