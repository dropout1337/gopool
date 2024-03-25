package gopool

import (
	"sync"
)

type ConcurrencyPool struct {
	maxThreads int
	wg         sync.WaitGroup
	available  chan struct{}
}

func New(maxThreads int) *ConcurrencyPool {
	return &ConcurrencyPool{
		maxThreads: maxThreads,
		available:  make(chan struct{}, maxThreads),
	}
}

func (cp *ConcurrencyPool) Wait() {
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
	newChan := make(chan struct{}, newSize)
	toPush := len(cp.available)

	if toPush > newSize {
		toPush = newSize
	}

	for i := 0; i < toPush; i++ {
		newChan <- struct{}{}
	}

	cp.maxThreads = newSize
	cp.available = newChan
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
