package mosync_test

import (
	"sync"
	"testing"

	. "github.com/sryoya/mosync"
)

type nTimes int

func (n *nTimes) increment() {
	*n++
}

func TestNtimes(t *testing.T) {
	nt := new(nTimes)
	nTimes := New(3)
	const N = 5
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			nTimes.Do(
				nt.increment,
			)
			wg.Done()
		}()
	}
	wg.Wait()
	if *nt != 3 {
		t.Errorf("once failed outside run: %d is not 3", *nt)
	}
}

func TestNtimesPanic(t *testing.T) {
	nt := New(1)
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("NTimes.Do did not panic")
			}
		}()
		nt.Do(func() {
			panic("failed")
		})
	}()

	nt.Do(func() {
		t.Fatalf("NTimes.Do executed, exceeding the specified times")
	})
}

func BenchmarkNtimes(b *testing.B) {
	n := New(uint32(b.N) / 2)
	f := func() {}
	for i := 0; i < b.N; i++ {
		n.Do(f)
	}
}

func BenchmarkNtimesWithOnce(b *testing.B) {
	n := New(1)
	f := func() {}
	for i := 0; i < b.N; i++ {
		n.Do(f)
	}
}

func BenchmarkNtimesWithFullTime(b *testing.B) {
	n := New(uint32(b.N))
	f := func() {}
	for i := 0; i < b.N; i++ {
		n.Do(f)
	}
}
