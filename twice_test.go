package mosync_test

import (
	"testing"

	. "github.com/sryoya/mosync"
)

type twice int

func (t *twice) increment() {
	*t++
}

func run(t *testing.T, twice *Twice, tw *twice, c chan bool) {
	twice.Do(func() { tw.increment() })
	if v := *tw; v != 1 && v != 2 {
		t.Errorf("twice failed inside run: %d is not 1 or 2", *tw)
	}
	c <- true
}

func TestTwice(t *testing.T) {
	tw := new(twice)
	twice := new(Twice)
	c := make(chan bool)
	const N = 10
	for i := 0; i < N; i++ {
		go run(t, twice, tw, c)
	}
	for i := 0; i < N; i++ {
		<-c
	}
	if *tw != 2 {
		t.Errorf("twice failed outside run: %d is not 2", *tw)
	}
}

func TestTwicePanic(t *testing.T) {
	var twice Twice
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("twice.Do did not panic")
			}
		}()
		twice.Do(func() {
			panic("failed")
		})
	}()

	// run twice
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("twice.Do did not panic")
			}
		}()
		twice.Do(func() {
			panic("failed")
		})
	}()

	twice.Do(func() {
		t.Fatalf("twice.Do called three times")
	})
}

func BenchmarkTwice(b *testing.B) {
	var twice Twice
	f := func() {}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			twice.Do(f)
			twice.Do(f)
		}
	})
}
