package roundrobin

import (
	"fmt"
)

type RingQueue[T any] struct {
	data   []T
	isFull bool
	start  int
	end    int
}

func NewRingQueue[T any](capacity int) *RingQueue[T] {
	return &RingQueue[T]{
		data:   make([]T, capacity),
		isFull: false,
		start:  0,
		end:    0,
	}
}

func (r *RingQueue[T]) String() string {
	return fmt.Sprintf(
		"[RRQ full:%v size:%d start:%d end:%d data:%v]",
		r.isFull,
		len(r.data),
		r.start,
		r.end,
		r.data)
}

func (r *RingQueue[T]) Push(elem T) error {
	if r.isFull {
		return fmt.Errorf("out of bounds push, container is full")
	}

	r.data[r.end] = elem
	r.end = (r.end + 1) % len(r.data)
	r.isFull = r.end == r.start

	return nil
}

func (r *RingQueue[T]) Pop() (T, error) {
	var res T
	if !r.isFull && r.start == r.end {
		return res, fmt.Errorf("empty queue")
	}

	res = r.data[r.start]
	r.start = (r.start + 1) % len(r.data)
	r.isFull = false

	return res, nil
}

func (r *RingQueue[T]) Size() int {
	res := r.end - r.start
	if res < 0 || (res == 0 && r.isFull) {
		res = len(r.data) - res
	}

	return res
}

func (r *RingQueue[T]) IsFull() bool {
	return r.isFull
}
