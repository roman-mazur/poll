package votes

import (
	"time"
)

// Vote represents the perception of a presentation listener at some given time.
// The Value field keeps the measure of how much the listener agrees or disagrees with the speaker.
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

// Label provides some context of what the presentation is about at some given time.
// It can be used as a marker pointing to a specific place in the presentation to correlate the presentation content
// with the votes.
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
