package mosync

import (
	"sync"
	"sync/atomic"
)

// Twice is an object that will perform at most two action.
//
// A Twice must not be copied after first use.
type Twice struct {
	// done indicates how many time the action has been performed.
	done uint32
	m    sync.Mutex
}

// Do calls the function f if and only if Do is being called for the
// first or second time for this instance of Once.
func (o *Twice) Do(f func()) {
	if atomic.LoadUint32(&o.done) < 2 {
		o.doSlow(f)
	}
}

func (o *Twice) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done < 2 {
		defer atomic.AddUint32(&o.done, 1)
		f()
	}
}
