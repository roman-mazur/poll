package votes

import (
	"sort"
	"sync"
	"time"
)

// Repository keeps the collected data about votes and labels.
type Repository struct {
	votes  mapStore[Vote]
	labels mapStore[Label]
}

func NewRepository() *Repository {
	return &Repository{
		votes:  makeStore[Vote](),
		labels: makeStore[Label](),
	}
}

func (r *Repository) Vote(v Vote) error {

	r.votes.add(v)
	return nil
}

func (r *Repository) Label(l Label) error {
	r.labels.add(l)
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
	for li < len(labels) && vi < len(votes) {
		label := labels[li]
		d := aggData{
			Label: label.Name,
			Time:  label.Timestamp,
		}

		pos, neg := make(votesSet), make(votesSet)
		for vi < len(votes) && (li == len(labels)-1 || votes[vi].Timestamp.Before(labels[li+1].Timestamp)) {
			vid := votes[vi].VoterId
			if votes[vi].Value == 0 {
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

type talkData interface {
	talkName() string
}

type storeItem interface {
	talkData
	touch
}

type mapStore[T storeItem] struct {
	sync.Mutex
	m map[string]chrono[T]
}

type touch interface {
	touch()
	t() time.Time
}

func makeStore[T storeItem]() mapStore[T] {
	return mapStore[T]{
		m: make(map[string]chrono[T]),
	}
}

func (ms *mapStore[T]) add(item T) {
	ms.Lock()
	defer ms.Unlock()
	key := item.talkName()
	item.touch()
	ms.m[key] = append(ms.m[key], item)
}

func (ms *mapStore[T]) get(key string) []T {
	ms.Lock()
	defer ms.Unlock()
	items := ms.m[key]
	res := make(chrono[T], len(items))
	copy(res, items)
	sort.Sort(res)
	return res
}

type chrono[T touch] []T

func (c chrono[T]) Len() int           { return len(c) }
func (c chrono[T]) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c chrono[T]) Less(i, j int) bool { return c[i].t().Before(c[j].t()) }
