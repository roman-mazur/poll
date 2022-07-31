package votes

import (
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

func (r *Repository) Add(v Vote) error {
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

	votes := r.votes.get(talkName)
	res.Votes = make([]voteAgg, len(votes))
	for i := range votes {
		v := &res.Votes[i]
		v.Pos = uint(votes[i].Value)
		v.Time = votes[i].Timestamp
	}

	labels := r.labels.get(talkName)
	res.Labels = make([]labelAgg, len(labels))
	for i := range labels {
		l := &res.Labels[i]
		l.Name = labels[i].Name
		l.Time = labels[i].Timestamp
	}

	return
}

type Aggregate struct {
	TalkName string     `json:"talk_name"`
	Votes    []voteAgg  `json:"votes"`
	Labels   []labelAgg `json:"labels"`
}

type voteAgg struct {
	Pos  uint      `json:"pos"`
	Neg  uint      `json:"neg"`
	Time time.Time `json:"time"`
}

type labelAgg struct {
	Name string    `json:"name"`
	Time time.Time `json:"time"`
}

type talkData interface {
	talkName() string
}

type mapStore[T talkData] struct {
	sync.Mutex
	m map[string][]T
}

func makeStore[T talkData]() mapStore[T] {
	return mapStore[T]{
		m: make(map[string][]T),
	}
}

func (ms *mapStore[T]) add(item T) {
	ms.Lock()
	defer ms.Unlock()
	key := item.talkName()
	ms.m[key] = append(ms.m[key], item)
}

func (ms *mapStore[T]) get(key string) []T {
	ms.Lock()
	defer ms.Unlock()
	coll := ms.m[key]
	// TODO: Process.
	res := make([]T, len(coll))
	copy(coll, res)
	return res
}
