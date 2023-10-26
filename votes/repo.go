package votes

import (
	"sort"
	"sync"
	"time"
)

// Repository keeps the collected data about votes and labels.
type Repository struct {
	votes  mapStore[Vote, *Vote]
	labels mapStore[Label, *Label]
}

func NewRepository() *Repository {
	return &Repository{
		votes:  makeStore[Vote, *Vote](time.Now),
		labels: makeStore[Label, *Label](time.Now),
	}
}

func (r *Repository) changeClock(c clock) {
	r.votes.c = c
	r.labels.c = c
}

func (r *Repository) Vote(v Vote) error {
	r.votes.add(&v)
	return nil
}

func (r *Repository) Label(l Label) error {
	for _, el := range r.labels.get(l.TalkName) {
		if el.Name == l.Name {
			return nil
		}
	}
	r.labels.add(&l)
	return nil
}

// Aggregate returns the aggregated data for a specific talk that can be used to provide a summary of the poll during
// the presentation.
func (r *Repository) Aggregate(talkName string) (res Aggregate) {
	res.TalkName = talkName

	labels := r.labels.get(talkName)
	votes := r.votes.get(talkName)

	res.Data = make([]aggData, len(labels))

	type votesSet map[string]struct{}

	li, vi := 0, 0
	for li < len(labels) {
		label := labels[li]
		d := aggData{
			Label: label.Name,
			Time:  label.Timestamp,
		}

		pos, neg := make(votesSet), make(votesSet)
		for vi < len(votes) && (li == len(labels)-1 || votes[vi].Timestamp.Before(labels[li+1].Timestamp)) {
			vid := votes[vi].VoterId
			if votes[vi].Value <= 5 {
				delete(pos, vid)
				neg[vid] = struct{}{}
			} else {
				delete(neg, vid)
				pos[vid] = struct{}{}
			}
			vi++
		}
		d.Pos, d.Neg = uint(len(pos)), uint(len(neg))

		res.Data[li] = d
		li++
	}
	return
}

type Aggregate struct {
	TalkName string    `json:"talk_name"`
	Data     []aggData `json:"data"`
}

type aggData struct {
	Pos   uint      `json:"pos"`
	Neg   uint      `json:"neg"`
	Time  time.Time `json:"time"`
	Label string    `json:"label"`
}

type itemData interface {
	talkName() string
	t() time.Time
}

type touch[T itemData] interface {
	touch(c clock)
	*T
}

type clock func() time.Time

type mapStore[T itemData, PT touch[T]] struct {
	c clock

	mu sync.Mutex
	m  map[string]chrono[T]
}

func makeStore[T itemData, PT touch[T]](c clock) mapStore[T, PT] {
	return mapStore[T, PT]{
		m: make(map[string]chrono[T]),
		c: c,
	}
}

func (ms *mapStore[T, PT]) add(item PT) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	key := (*item).talkName()
	item.touch(ms.c)
	ms.m[key] = append(ms.m[key], *item)
}

func (ms *mapStore[T, PT]) get(key string) []T {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	items := ms.m[key]
	res := make(chrono[T], len(items))
	copy(res, items)
	sort.Sort(res)
	return res
}

type chrono[T itemData] []T

func (c chrono[T]) Len() int           { return len(c) }
func (c chrono[T]) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c chrono[T]) Less(i, j int) bool { return c[i].t().Before(c[j].t()) }
