# go-container-roundrobin
Ring (Circular or RoundRobin) queue-like conainer with fixed memory footprint

Read full post here: https://www.sergetoro.com/golang-round-robin-queue-from-scratch

## What a Ring Queue is and Why use it
So, what is a Ring Queue and why it's useful:
1. It uses a static / fixed array buffer (no costly allocations)
1. It does not shift elements on element removal (no costly memcopy)
1. It can start anywhere in the buffer and loop over the end

![img](https://www.sergetoro.com/content/images/2023/06/RQ.svg)

## RingQueue Performance
Do you like benchmarks? I love them, even though many cases they are non-exhaustive, relatively synthetic and might give a skewed view of reality :)
So, let's see how our implementation performs against a simple Go array.
Remember, what we're looking for is an array slow down due to excessive copying.

```go
func BenchmarkRR(b *testing.B) {
	rr := NewRingQueue[int](100_000)

	for n := 0; n < b.N; n++ {
		if rr.IsFull() {
			rr.Pop()
		}
		rr.Push(n)
	}
}

func BenchmarkArray(b *testing.B) {
	var ar [100_000]int
	size := 0

	for n := 0; n < b.N; n++ {
		if size >= 100_000 {
			copy(ar[0:], ar[1:])
			size--
		}

		ar[size] = n
		size++
	}
}
```

Which one do you thing would be the fastest? 😉 Let's run it!

```bash
go test -bench=. -benchmem
```

Here's results on my Intel NUC Ubuntu Linux machine:

```bash
goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i7-10710U CPU @ 1.10GHz
BenchmarkRR-12          59798865                18.50 ns/op
BenchmarkArray-12        1000000             17321.00 ns/op
PASS
```

If I run the same test with just 10_000 queue size instead of 100_000:

```bash
goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i7-10710U CPU @ 1.10GHz
BenchmarkRR-12          57184480                18.65 ns/op
BenchmarkArray-12        1000000              1068 ns/op
PASS
```

Let's plot the relation of run times and queue lengths:
![img](https://www.sergetoro.com/content/images/2023/06/image.svg)
