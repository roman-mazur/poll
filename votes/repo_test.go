package votes

import (
	"fmt"
	"testing"
	"time"
)

func TestRepository_Aggregate(t *testing.T) {
	const talk = "talk1"
	start := time.Now()

	r := NewRepository()

	// Lines represent the calls handling sequence, not the chronological timeline.

	_ = r.Label(Label{TalkName: talk, Name: "label1", Timestamp: start})

	_ = r.Vote(Vote{TalkName: talk, VoterId: "u1", Value: 8, Timestamp: start.Add(time.Second)})
	_ = r.Vote(Vote{TalkName: talk, VoterId: "u2", Value: 8, Timestamp: start.Add(time.Second * 3 / 2)})
	_ = r.Vote(Vote{TalkName: talk, VoterId: "u1", Value: 1, Timestamp: start.Add(time.Second * 2)})
	_ = r.Vote(Vote{TalkName: talk, VoterId: "u3", Value: 7, Timestamp: start.Add(time.Second * 3)})

	_ = r.Label(Label{TalkName: talk, Name: "label2", Timestamp: start.Add(10 * time.Second)})
	_ = r.Vote(Vote{TalkName: talk, VoterId: "u1", Value: 3, Timestamp: start.Add(11 * time.Second)})
	_ = r.Vote(Vote{TalkName: talk, VoterId: "u2", Value: 9, Timestamp: start.Add(10 * time.Second)})
	_ = r.Vote(Vote{TalkName: talk, VoterId: "u3", Value: 2, Timestamp: start.Add(13 * time.Second)}) // Will be added to the next label.

	_ = r.Label(Label{TalkName: talk, Name: "label3", Timestamp: start.Add(12 * time.Second)})

	a := r.Aggregate(talk)
	t.Log(a.Data)
	if len(a.Data) != 3 {
		t.Fatal("bad length of data")
	}
	for i := 0; i < 3; i++ {
		if a.Data[i].Label != fmt.Sprintf("label%d", i+1) {
			t.Error("bad label for", i)
		}
		if i > 1 && !a.Data[i].Time.After(a.Data[i-1].Time) {
			t.Errorf("%d not after %d", i, i-1)
		}
	}
	if a.Data[0].Pos != 2 || a.Data[0].Neg != 1 {
		t.Error("bad first record")
	}
	if a.Data[1].Pos != 1 || a.Data[1].Neg != 1 {
		t.Error("bad second record")
	}
	if a.Data[2].Pos != 0 || a.Data[2].Neg != 1 {
		t.Error("bad third record")
	}
}
