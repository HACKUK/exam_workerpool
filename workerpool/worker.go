package workerpool

import (
"log"
"sync/atomic"
"runtime"
//"net/http"
//"io/ioutil"
)

type Worker struct {
	Pool  	   *Dispatcher
	Index      int
	WorkerPool chan chan Job
	JobChannel chan Job
	Quit       chan bool
}

func NewWorker(workerPool chan chan Job, i int, p *Dispatcher) Worker {
	return Worker{
		Pool: p,
		Index : i,
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		Quit:       make(chan bool, 1)}
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				job.Handler(job.Input)
				atomic.AddUint64(&(w.Pool.Ops), 1)
				runtime.Gosched()
			case <-w.Quit:
				log.Printf("[workerpool], worker %d stoped.\n", w.Index)
				return
			}
		}
	}()
}
func (w Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}
