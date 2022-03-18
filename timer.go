package tools

import (
	"math/rand"
	"sync"
	"time"
)

// RandomTick is similar to time.Ticker but ticks at random intervals between the min and max duration
type Randomtick struct {
	MsMin  time.Duration
	MsMax  time.Duration
	Ch     chan int
	StopCh chan int
	Times  int
	sync.Mutex
}

// Stop terminates the ticker goroutine and closes the C channel.
func (rt *Randomtick) Stop() {
	c := make(chan struct{})
	rt.stopCh <- c
	<-c
}

func (rt *Randomtick) runTick() {
	defer close(rt.Ch)
	t := time.NewTimer(rt.nextInterval())
	for {
		// either a stop signal or a timeout
		select {
		case c := <-rt.StopCh:
			t.Stop()
			close(c)
			return
		case <-t.Ch:
			select {
			case rt.Ch <- time.Now():
				t.Stop()
				t = time.NewTimer(rt.nextInterval())
			default:
				// there could be noone receiving...
			}
		}
	}
}

// NewUltimateRandomTick returns a pointer to an initialized instance of the
// RandomTicker. Min and max are durations of the shortest and longest allowed
func NewUltimateRandomTick(msMin, msMax time.Duration) *Randomtick {
	rt := &Randomtick{
		Ch:     make(chan int, 10),
		MsMin:  msMin.Nanoseconds(),
		MsMax:  msMax.Nanoseconds(),
		StopCh: make(chan int, 1),
	}
	go rt.runTick()
	return rt
}

func (rt *RandomTicker) nextInterval() time.Duration {
	interval := rand.Int63n(rt.MsMax-rt.MsMin) + rt.MsMin
	return time.Duration(interval) * time.Nanosecond
}
