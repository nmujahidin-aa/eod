package utils

import "context"

type WorkerPool interface {
	Run()
	AddTask(task func(workerId int32))
	Stop()
}

type workerPool struct {
	ctx            context.Context
	stopChan       chan bool
	numberOfWorker int32
	queueTaskChan  chan func(workerId int32)
}

func NewWorkerPool(ctx context.Context, numberOfWorker int32) WorkerPool {
	return &workerPool{
		ctx:            ctx,
		numberOfWorker: numberOfWorker,
		queueTaskChan:  make(chan func(workerId int32), numberOfWorker),
		stopChan:       make(chan bool),
	}
}
func (wp *workerPool) Run() {
	var i int32
	for i = 1; i <= wp.numberOfWorker; i++ {
		go func(ctx context.Context, workerId int32) {
			for {
				select {
				case task := <-wp.queueTaskChan:
					task(workerId)
				case <-ctx.Done():
					return
				case <-wp.stopChan:
					return
				}
			}
		}(wp.ctx, i)
	}
}
func (wp *workerPool) Stop() {
	wp.stopChan <- true
}
func (wp *workerPool) AddTask(task func(workerId int32)) {
	wp.queueTaskChan <- task
}
