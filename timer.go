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
	c := make(chan int)
	rt.StopCh <- c
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
		case <-t.C:
			select {
			case rt.Ch <- int(time.Now().Unix()):
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
		Ch:     make(chan int),
		MsMin:  time.Duration(msMin.Nanoseconds()),
		MsMax:  time.Duration(msMax.Nanoseconds()),
		StopCh: make(chan int),
	}
	go rt.runTick()
	return rt
}

func (rt *Randomtick) nextInterval() time.Duration {
	interval := time.Duration(rand.Int63n(int64(rt.MsMax-rt.MsMin))) + rt.MsMin
	return time.Duration(interval) * time.Nanosecond
}
