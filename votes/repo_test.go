package votes

import (
	"fmt"
	"testing"
	"time"
)

func TestRepository_Aggregate(t *testing.T) {
	const talk = "talk1"
	start := time.Now()
	now := start

	r := NewRepository()
	r.changeClock(func() time.Time { return now })

	// Lines represent the calls handling sequence, not the chronological timeline.

	_ = r.Label(Label{TalkName: talk, Name: "label1"})

	now = start.Add(time.Second)
	_ = r.Vote(Vote{TalkName: talk, VoterId: "u1", Value: 8})
	now = start.Add(time.Second / 2)
	_ = r.Vote(Vote{TalkName: talk, VoterId: "u2", Value: 8})
	now = start.Add(time.Second * 2)
	_ = r.Vote(Vote{TalkName: talk, VoterId: "u1", Value: 1})
	now = start.Add(time.Second * 3)
	_ = r.Vote(Vote{TalkName: talk, VoterId: "u3", Value: 7})

	now = start.Add(time.Second * 10)
	_ = r.Label(Label{TalkName: talk, Name: "label2"})
	now = start.Add(time.Second * 11)
	_ = r.Vote(Vote{TalkName: talk, VoterId: "u1", Value: 3})
	now = start.Add(time.Second * 10)
	_ = r.Vote(Vote{TalkName: talk, VoterId: "u2", Value: 9})
	now = start.Add(time.Second * 13)
	_ = r.Vote(Vote{TalkName: talk, VoterId: "u3", Value: 2}) // Will be added to the next label.

	now = start.Add(time.Second * 12)
	_ = r.Label(Label{TalkName: talk, Name: "label3"})

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

func TestRepository_ReproduceIssues(t *testing.T) {
	t.Run("no votes", func(t *testing.T) {
		const talk = "talk1"
		r := NewRepository()
		_ = r.Label(Label{TalkName: talk, Name: "label1"})

		res := r.Aggregate(talk)
		if res.Data[0].Label != "label1" {
			t.Error("bad first label")
		}
	})
	t.Run("zero reports", func(t *testing.T) {
		const talk = "talkKey/something"
		now := time.Now()
		r := NewRepository()
		r.changeClock(func() time.Time {
			return now
		})

		_ = r.Label(Label{TalkName: talk, Name: "label1"})
		now = now.Add(time.Second)
		_ = r.Vote(Vote{TalkName: talk, VoterId: "voter1", Value: 10})

		res := r.Aggregate(talk)
		if res.Data[0].Pos != 1 {
			t.Error("incorrect calculation")
		}
	})
	t.Run("duplicates", func(t *testing.T) {
		const talk = "talkKey/something"
		r := NewRepository()
		_ = r.Label(Label{TalkName: talk, Name: "label1"})
		_ = r.Label(Label{TalkName: talk, Name: "label1"})

		res := r.Aggregate(talk)
		if len(res.Data) != 1 {
			t.Error("too much data")
		}
	})
	t.Run("wrong grouping", func(t *testing.T) {
		const talk = "talkKey/something"
		now := time.Now()
		r := NewRepository()
		r.changeClock(func() time.Time { return now })

		_ = r.Label(Label{TalkName: talk, Name: "label1"})
		now = now.Add(time.Second)
		_ = r.Label(Label{TalkName: talk, Name: "label2"})
		now = now.Add(time.Second)
		_ = r.Vote(Vote{TalkName: talk, VoterId: "voter1", Value: 10})

		res := r.Aggregate(talk)
		if len(res.Data) != 2 {
			t.Fatal("unexpected groups count")
		}
		if res.Data[1].Pos != 1 {
			t.Log(res)
			t.Error("wrong grouping")
		}
	})
}
