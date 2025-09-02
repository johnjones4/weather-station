package store

import (
	"context"
	"main/core"
	"time"
)

type weatherLinkedListNode struct {
	value *core.Weather
	next  *weatherLinkedListNode
}

func (n *weatherLinkedListNode) size() int {
	return n._size(0)
}

func (n *weatherLinkedListNode) add(value *core.Weather) {
	if n.next == nil {
		n.next = &weatherLinkedListNode{
			value: value,
		}
	} else {
		n.next.add(value)
	}
}

func (n *weatherLinkedListNode) _size(d int) int {
	dd := d + 1
	if n.next != nil {
		return n._size(dd)
	}
	return dd
}

type MemoryStore struct {
	weather *weatherLinkedListNode
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		weather: nil,
	}
}

func (s *MemoryStore) Save(ctx context.Context, w *core.Weather) error {
	if s.weather == nil {
		s.weather = &weatherLinkedListNode{
			value: w,
		}
	} else {
		s.weather.add(w)
		if s.weather.size() > 100 {
			s.weather = s.weather.next
		}
	}
	return nil
}

func (s *MemoryStore) Get(ctx context.Context, start, end time.Time) ([]core.Weather, error) {
	out := make([]core.Weather, 0)
	next := s.weather
loop:
	for {
		if next == nil {
			break loop
		}
		if next.value.Timestamp.After(start) && next.value.Timestamp.Before(end) {
			out = append(out, *next.value)
		}
		next = next.next
	}
	return out, nil
}
