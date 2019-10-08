package worker

import (
	"sync"
	"time"
)

// Worker : class to execute a periodic task
type Worker struct {
	sync.Mutex
	PeriodMs  int
	cond      *sync.Cond
	waitgroup sync.WaitGroup
	quit      chan struct{}
	running   bool
	work      func() bool //stops if false
}

// NewWorker : constructor
func NewWorker() *Worker {
	r := &Worker{
		running: false,
	}
	r.cond = sync.NewCond(r)
	return r
}

// Start : start worker returns success
func (d *Worker) Start(periodMs int, f func() bool) bool {
	ret := false
	if !d.running {
		d.PeriodMs = periodMs
		d.work = f
		go d.run()
		d.waitgroup.Add(1)
		ret = true
	}
	return ret
}

// Stop : stops worker returns success
func (d *Worker) Stop() bool {
	ret := false
	if d.running {
		close(d.quit)
		d.waitgroup.Wait()
		ret = true
	}
	return ret
}

func (d *Worker) run() {
	ticker := time.NewTicker(time.Duration(d.PeriodMs) * time.Millisecond)
	d.quit = make(chan struct{})
	d.running = true
	for d.running {
		select {
		case <-ticker.C:
			// do stuff
			if !d.work() {
				ticker.Stop()
				d.running = false
			}
		case <-d.quit:
			ticker.Stop()
			d.running = false
		}
	}
	d.waitgroup.Done()
}
