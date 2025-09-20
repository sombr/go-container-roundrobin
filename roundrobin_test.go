package roundrobin

import (
	"fmt"
	"testing"
	"time"
)

func eqSlices[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for idx := 0; idx < len(a); idx++ {
		if a[idx] != b[idx] {
			return false
		}
	}

	return true
}

func TestToString(t *testing.T) {
	obj := NewRingQueue[int](10)
	expected := "[RRQ full:false size:10 start:0 end:0 data:[0 0 0 0 0 0 0 0 0 0]]"
	actual := fmt.Sprint(obj)

	if actual != expected {
		t.Fatalf("Mismatch, expected:%s, found:%s", expected, actual)
	}
}

func TestPushEnough(t *testing.T) {
	obj := NewRingQueue[int](10)
	for idx := 0; idx < 10; idx++ {
		err := obj.Push(idx)
		if err != nil {
			t.Fatalf("Unexpected error in adding an element with index %d", idx)
		}
	}

	expected := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	if !eqSlices(obj.data, expected) {
		t.Fatalf("Container data mismatch, expected:%v, found:%v", expected, obj.data)
	}

    if obj.Size() != 10 {
        t.Fatalf("Incorrect size reported, expected:%d, found:%d", 10, obj.Size())
    }
}

func TestPushOver(t *testing.T) {
	obj := NewRingQueue[int](10)
	for idx := 0; idx < 10; idx++ {
		err := obj.Push(idx)
		if err != nil {
			t.Fatalf("Unexpected error in adding an element with index %d", idx)
		}
	}

	err := obj.Push(100)
	if err == nil {
		t.Fatalf("Expected overflow error")
	}

	expected := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	if !eqSlices(obj.data, expected) {
		t.Fatalf("Container data mismatch, expected:%v, found:%v", expected, obj.data)
	}
}

func TestPushPop(t *testing.T) {
	obj := NewRingQueue[int](10)
	for idx := 0; idx < 8; idx++ {
		obj.Push(idx)
	}
	for idx := 0; idx < 5; idx++ {
		e, err := obj.Pop()
		if err != nil || e != idx {
			t.Fatalf("inconsistent behavior")
		}
	}
	for idx := 0; idx < 7; idx++ {
		obj.Push(100 + idx)
	}

	expected := []int{102, 103, 104, 105, 106, 5, 6, 7, 100, 101}

	if !eqSlices(obj.data, expected) {
		t.Fatalf("Container data mismatch, expected:%v, found:%v", expected, obj.data)
	}

	if obj.Size() != 10 {
		t.Fatalf("inconsistent size: %d", obj.Size())
	}

	for idx := 0; idx < 10; idx++ {
		e, _ := obj.Pop()
		if e != expected[(5+idx)%10] {
			t.Fatalf("inconsistent behavior")
		}
	}
}

func sim(capacity int) {
	ar := make([]int, capacity)
	size := 0

	start := time.Now()
	for n := 0; n < 1000000; n++ {
		if size >= len(ar) {
			copy(ar[0:], ar[1:])
			size--
		}

		ar[size] = n
		size++
	}

	fmt.Printf("%d took %v\n", capacity, time.Since(start).Seconds())
}

func simRR(capacity int) {
	rr := NewRingQueue[int](capacity)

	start := time.Now()
	for n := 0; n < 1000000; n++ {
		if rr.IsFull() {
			rr.Pop()
		}
		rr.Push(n)
	}

	fmt.Printf("%d took %v\n", capacity, time.Since(start).Seconds())
}

func TestPrimitiveAsymptoticPerformance(t *testing.T) {
	fmt.Println("Standard array")
	for idx := 7; idx < 14; idx++ {
		sim(1 << idx)
	}

	fmt.Println("RoundRobin (ring) queue")
	for idx := 7; idx < 14; idx++ {
		simRR(1 << idx)
	}
}

func BenchmarkRR(b *testing.B) {
	rr := NewRingQueue[int](1_000)

	for n := 0; n < b.N; n++ {
		if rr.IsFull() {
			rr.Pop()
		}
		rr.Push(n)
	}
}

func BenchmarkArray(b *testing.B) {
	var ar [1_000]int
	size := 0

	for n := 0; n < b.N; n++ {
		if size >= len(ar) {
			copy(ar[0:], ar[1:])
			size--
		}

		ar[size] = n
		size++
	}
}
