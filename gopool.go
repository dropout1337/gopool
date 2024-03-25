package main

import (
	"sync"
)

type ConcurrencyPool struct {
	maxThreads int
	wg         sync.WaitGroup
	mu         sync.Mutex
	available  chan struct{}
}

func New(maxThreads int) *ConcurrencyPool {
	return &ConcurrencyPool{
		maxThreads: maxThreads,
		available:  make(chan struct{}, maxThreads),
	}
}

func (cp *ConcurrencyPool) Wait() {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	cp.available <- struct{}{}
	cp.wg.Add(1)
}

func (cp *ConcurrencyPool) Done() {
	<-cp.available
	cp.wg.Done()
}

func (cp *ConcurrencyPool) WaitUntilDone() {
	cp.wg.Wait()
}

func (cp *ConcurrencyPool) ResizePool(newSize int) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	cp.wg.Wait()

	cp.maxThreads = newSize
	cp.available = make(chan struct{}, newSize)
}

func (cp *ConcurrencyPool) SetMaxThreads(newMaxThreads int) {
	cp.maxThreads = newMaxThreads
}

func (cp *ConcurrencyPool) GetCurrentThreadCount() int {
	return len(cp.available)
}

func (cp *ConcurrencyPool) IsAvailable() bool {
	return cp.GetCurrentThreadCount() < cp.maxThreads
}

func (cp *ConcurrencyPool) GetMaxThreads() int {
	return cp.maxThreads
}
