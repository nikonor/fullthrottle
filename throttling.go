package fullthrottle

import (
	"sync/atomic"
	"time"
)

// TODO: раcписать доку
type Throttling struct {
	BeginDelay   time.Duration
	MaxDelay     time.Duration
	IncDelay     time.Duration
	MultDelay    time.Duration
	MaxCount     int64
	NextDelay    func(Throttling) time.Duration
	count        int64
	currentDelay time.Duration
	DoneChan     chan struct{}
}

func (o *Throttling) Throttling() chan bool {
	ch := make(chan bool)
	atomic.AddInt64(&o.count, 1)

	go func(ch chan bool) {
		if o.MaxCount != 0 && o.MaxCount <= o.count {
			ch <- true
		}
		t := time.NewTicker(o.getCurrentDelay())
		for {
			select {
			case <-o.DoneChan:
				ch <- true
				return
			case <-t.C:
				o.setNextCurrentDelay()
				close(ch)
				return
			}
		}
	}(ch)

	return ch
}

func (o *Throttling) getCurrentDelay() time.Duration {
	switch {
	case o.currentDelay == 0 && o.NextDelay != nil:
		o.currentDelay = o.NextDelay(*o)
	case o.currentDelay != 0:
		return o.currentDelay
	case o.BeginDelay != 0:
		o.currentDelay = o.BeginDelay
	default:
		o.currentDelay = time.Second
	}
	return o.currentDelay
}

func (o *Throttling) setNextCurrentDelay() {
	switch {
	case o.NextDelay != nil:
		o.currentDelay = o.NextDelay(*o)
	case o.IncDelay > 0:
		o.currentDelay += o.IncDelay
	case o.MultDelay > 0:
		o.currentDelay *= o.MultDelay
	}

	// проверяем, что не превысили максимальную задержку
	if o.MaxDelay > 0 && o.currentDelay >= o.MaxDelay {
		o.currentDelay = o.MaxDelay
	}
}
