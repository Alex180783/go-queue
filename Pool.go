package queue

import (
	"fmt"
	//	"log"
)

// ==[ Pool ]==

func NewPool(name string, id, size int) *Pool { // {{{
	return &Pool{
		id:          id,
		name:        name,
		size:        size,
		workers:     make([]WorkerInterface, 0),
		workerQueue: make(WorkerQueue, size),
	}
} // }}}

type Pool struct {
	id            int
	name          string
	size          int
	pending       uint
	workerFactory WorkerFactoryInterface
	workers       []WorkerInterface
	workerQueue   WorkerQueue

	index int
}

func (this *Pool) Id() int { // {{{
	return this.id
} // }}}

func (this *Pool) Name() string { // {{{
	return this.name
} // }}}

func (this *Pool) Size() int { // {{{
	return this.size
} // }}}

func (this *Pool) Pending() uint { // {{{
	return this.pending
} // }}}

func (this *Pool) Info() string { // {{{
	return fmt.Sprintf("%sPool:%d#%d", this.name, this.id, this.size)
} // }}}

func (this *Pool) SetWorkerFactory(f WorkerFactoryInterface) { // {{{
	this.workerFactory = f
} // }}}

func (this *Pool) WorkerFactory() WorkerFactoryInterface { // {{{
	if this.workerFactory == nil {
		this.workerFactory = NewWorkerFactory(this.name)
	}

	return this.workerFactory
} // }}}

func (this *Pool) Dispatch(job JobInterface) { // {{{
	this.pending++
	worker := <-this.workerQueue
	worker <- job
} // }}}

func (this *Pool) Run() { // {{{
	for i := 1; i <= this.size; i++ {
		worker := this.WorkerFactory().Create(this.id, this.workerQueue)
		this.workers = append(this.workers, worker)
		worker.Start()
	}
} // }}}

func (this *Pool) close() { // {{{
	close(this.workerQueue)

	if len(this.workers) > 0 {
		this.workers = this.workers[:0]
	}
} // }}}

func (this *Pool) Close() { // {{{
	defer this.close()

	// Stop workers
	for i, _ := range this.workers {
		this.workers[i].Close()
	}
} // }}}

func (this *Pool) Start() { // {{{
	go this.Run()
} // }}}

func (this *Pool) Stop() { // {{{
	go this.Close()
} // }}}
