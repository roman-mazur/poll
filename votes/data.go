package votes

import (
	"time"

	"rmazur.io/poll/votes/structure"
)

// Vote represents the perception of a presentation listener at some given time.
// The Value field keeps the measure of how much the listener agrees or disagrees with the speaker.
type Vote structure.Vote

func (v Vote) talkName() string {
	return v.TalkName
}

func (v Vote) touch() {
	v.Timestamp = time.Now()
}

func (v Vote) t() time.Time {
	return v.Timestamp
}

// Label provides some context of what the presentation is about at some given time.
// It can be used as a marker pointing to a specific place in the presentation to correlate the presentation content
// with the votes.
type Label structure.Label

func (l Label) talkName() string {
	return l.TalkName
}

func (l Label) touch() {
	l.Timestamp = time.Now()
}

func (l Label) t() time.Time {
	return l.Timestamp
}
