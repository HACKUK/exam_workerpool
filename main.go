package main

import (
	"exam_workerpool/jobqueue"
	"exam_workerpool/workerpool"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"time"
)

func SetupCPU() {
	num := runtime.NumCPU()
	runtime.GOMAXPROCS(num)
}

func main() {

	//SetupCPU()
	//log.SetFlags(log.Lshortfile | log.LstdFlags)
	defer time.Sleep(1 * time.Second)
	//创建工作池
	jobQueue := make(workerpool.JobQueue, 300)
	d := workerpool.NewDispatcher(10)
	d.Run(jobQueue)
	defer d.Stop()

	go func(jobs workerpool.JobQueue) {
		for i := 0; i < 100000000; i++ {
			jobs <- workerpool.Job{Input: i,
				Handler: func(i interface{}) {
					//fmt.Printf("f1 pop %d from JobQueue\n", i.(int))
					for i := 0; i < 100; i++ {
					}
				},
			}
		}
	}(jobQueue)

	go func(jobs workerpool.JobQueue) {
		for i := 0; i < 10; i++ {
			jobs <- workerpool.Job{
				Input: i,
				Handler: func(i interface{}) {
					//fmt.Printf("f2 pop %d from JobQueue\n", i.(int))
				},
			}
		}
	}(jobQueue)

	//创建任务队列
	worker := jobqueue.NewWorker(10)
	worker.Start()
	defer worker.Stop()
	go func(w jobqueue.Worker) {
		for i := 1; i <= 1; i++ {
			job := jobqueue.Job{
				Input: i,
				Handler: func(i interface{}) {
					fmt.Printf("pop %d from jobqueue\n", i.(int))
				},
			}
			w.Push(job)
		}
	}(worker)

	go func(d *workerpool.Dispatcher) {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				ops := atomic.LoadUint64(&(d.Ops))
				log.Printf("hanled %dw jobs\n", ops/10000)
				if ops >= 100000000 {
					d.Stop()
					log.Printf("stop workerpool, hanled %dw jobs\n", atomic.LoadUint64(&(d.Ops))/10000)
					return
				}
			}
		}
	}(d)
	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, os.Kill)
ConsumerLoop:
	for {
		select {
		case <-signals:
			break ConsumerLoop
		}
	}
}
