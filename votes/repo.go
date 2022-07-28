package votes

import "sync"

type talkData interface {
	TalkName() string
}

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
	key := item.TalkName()
	ms.m[key] = append(ms.m[key], item)
}
