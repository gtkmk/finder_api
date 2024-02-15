package worker

import (
	"sync"

	"github.com/gtkmk/finder_api/core/port"
)

type Worker struct {
	id int
}

type WorkerPool struct {
	NumWorkers int
	Workers    []*Worker
	TaskChan   chan port.WorkerTask
	ErrorChan  chan error
	Wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers int) port.WorkerInterface {
	return &WorkerPool{
		NumWorkers: numWorkers,
		TaskChan:   make(chan port.WorkerTask),
		ErrorChan:  make(chan error),
	}
}

func (pool *WorkerPool) StartWorkers() {
	for i := 0; i < pool.NumWorkers; i++ {
		worker := &Worker{id: i}
		pool.Workers = append(pool.Workers, worker)
		go worker.start(pool.TaskChan, pool.ErrorChan)
	}
}

func (pool *WorkerPool) HandleErrorsAsync(externalErrorChanel chan error) {
	go func() {
		for err := range pool.ErrorChan {
			externalErrorChanel <- err
			pool.Wg.Done()
		}
	}()
}

func (worker *Worker) start(taskChan <-chan port.WorkerTask, errorChan chan error) {
	for task := range taskChan {
		errorChan <- task.ExecuteWorker()
	}
}

func (pool *WorkerPool) AddTask(task port.WorkerTask) {
	pool.Wg.Add(1)
	pool.TaskChan <- task
}

func (pool *WorkerPool) Wait() {
	pool.Wg.Wait()
	close(pool.TaskChan)
	close(pool.ErrorChan)
}
