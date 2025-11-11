package workerpool

import (
	"assignment/config"
	log "github.com/sirupsen/logrus"
	"sync"
)

type WorkerPool struct {
	wg         sync.WaitGroup
	workerSize int
	tasks      chan func()
	cfg        config.WorkerPool
}

func NewWorkerPool(workerPoolCfg config.WorkerPool) *WorkerPool {
	return &WorkerPool{
		workerSize: workerPoolCfg.Size,
		cfg:        workerPoolCfg,
	}
}

func (w *WorkerPool) StartWorkers() {
	w.tasks = make(chan func(), w.cfg.BufferSize)

	w.wg.Add(w.workerSize)

	for i := 0; i < w.workerSize; i++ {
		go func(id int) {
			defer w.wg.Done()

			for task := range w.tasks {
				log.Infof("processing in workerID:%d", id)

				task()
			}
		}(i)
	}
}

// if buffer is full it will create blocking handle it
func (w *WorkerPool) Submit(task func()) {
	w.tasks <- task
}

func (w *WorkerPool) Close() {
	close(w.tasks)
	w.wg.Wait()
}
