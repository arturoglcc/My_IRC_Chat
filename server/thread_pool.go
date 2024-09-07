package main

import (
	"sync"
)

// Worker es una estructura que representa una gorutina en el pool.
type Worker struct {
	id        int
	taskQueue chan func()
	wg        *sync.WaitGroup
}

// NewWorker crea una nueva instancia de Worker.
func NewWorker(id int, wg *sync.WaitGroup) *Worker {
	worker := &Worker{
		id:        id,
		taskQueue: make(chan func()),
		wg:        wg,
	}
	go worker.start()
	return worker
}

// start comienza a escuchar tareas en la cola de tareas del Worker.
func (w *Worker) start() {
	for task := range w.taskQueue {
		task()
		w.wg.Done()
	}
}

// ThreadPool representa un conjunto de Workers.
type ThreadPool struct {
	workers []*Worker
	wg      sync.WaitGroup
	index   int
	mu      sync.Mutex
}

// NewThreadPool crea una nueva instancia de ThreadPool.
func NewThreadPool(numWorkers int) *ThreadPool {
	pool := &ThreadPool{}
	for i := 0; i < numWorkers; i++ {
		pool.workers = append(pool.workers, NewWorker(i, &pool.wg))
	}
	return pool
}

// Submit envía una tarea al pool de trabajadores.
func (tp *ThreadPool) Submit(task func()) {
	tp.wg.Add(1)

	tp.mu.Lock()
	worker := tp.workers[tp.index]
	tp.index = (tp.index + 1) % len(tp.workers)
	tp.mu.Unlock()

	worker.taskQueue <- task
}

// Wait espera a que todas las tareas se completen.
func (tp *ThreadPool) Wait() {
	tp.wg.Wait()
	for _, worker := range tp.workers {
		close(worker.taskQueue)
	}
}
