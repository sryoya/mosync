package mosync

import (
	"sync"
	"sync/atomic"
)

// NTimes is an object that will perform at most N times action.
//
// A NTimes must not be copied after first use.
type NTimes struct {
	// remaining indicates how many more times the action can be performed.
	remaining uint32
	m         sync.Mutex
}

func New(n uint32) *NTimes {
	return &NTimes{remaining: n}
}

// Do calls the function f if and only if Do is being called for the
// specified time.
func (o *NTimes) Do(f func()) {
	if atomic.LoadUint32(&o.remaining) != 0 {
		o.doSlow(f)
	}
}

func (o *NTimes) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if atomic.LoadUint32(&o.remaining) != 0 {
		defer atomic.AddUint32(&o.remaining, ^uint32(0))
		f()
	}
}
