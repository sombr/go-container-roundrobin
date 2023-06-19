package roundrobin

import (
	"fmt"
)

type RoundRobinQueue[T any] struct {
	data   []T
	isFull bool
	start  int
	end    int
}

func NewRoundRobinQueue[T any](capacity int) *RoundRobinQueue[T] {
	return &RoundRobinQueue[T]{
		data:   make([]T, capacity),
		isFull: false,
		start:  0,
		end:    0,
	}
}

func (r *RoundRobinQueue[T]) String() string {
	return fmt.Sprintf(
		"[RRQ full:%v size:%d start:%d end:%d data:%v]",
		r.isFull,
		len(r.data),
		r.start,
		r.end,
		r.data)
}

func (r *RoundRobinQueue[T]) Push(elem T) error {
	if r.isFull {
		return fmt.Errorf("out of bounds push, container is full")
	}

	r.data[r.end] = elem
	r.end = (r.end + 1) % len(r.data)
	r.isFull = r.end == r.start

	return nil
}

func (r *RoundRobinQueue[T]) Pop() (T, error) {
	var res T
	if !r.isFull && r.start == r.end {
		return res, fmt.Errorf("empty queue")
	}

	res = r.data[r.start]
	r.start = (r.start + 1) % len(r.data)
	r.isFull = false

	return res, nil
}

func (r *RoundRobinQueue[T]) Size() int {
	res := r.end - r.start
	if res < 0 || (res == 0 && r.isFull) {
		res = len(r.data) - res
	}

	return res
}

func (r *RoundRobinQueue[T]) IsFull() bool {
	return r.isFull
}
