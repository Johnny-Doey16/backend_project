package worker

import "sync"

type WorkerPool struct {
	maxWorkers int
	jobQueue   chan Job
	waitGroup  *sync.WaitGroup
}

type Job func()

func NewWorkerPool(maxWorkers int) *WorkerPool {
	return &WorkerPool{
		maxWorkers: maxWorkers,
		jobQueue:   make(chan Job),
		waitGroup:  &sync.WaitGroup{},
	}
}

func (wp *WorkerPool) Run() {
	for i := 0; i < wp.maxWorkers; i++ {
		wp.waitGroup.Add(1)
		go func() {
			defer wp.waitGroup.Done()
			for job := range wp.jobQueue {
				job()
			}
		}()
	}
}

func (wp *WorkerPool) AddJob(job Job) {
	wp.jobQueue <- job
}

func (wp *WorkerPool) Shutdown() {
	close(wp.jobQueue)
	wp.waitGroup.Wait()
}
