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

func (cp *ConcurrencyPool) Execute(fn func()) {
	cp.Wait()
	cp.wg.Add(1)

	go func() {
		defer cp.Done()
		fn()
	}()
}

func (cp *ConcurrencyPool) ResizePool(newSize int) {
	if newSize > cp.maxThreads {
		for i := cp.maxThreads; i < newSize; i++ {
			cp.available <- struct{}{}
		}
		cp.maxThreads = newSize
	} else if newSize < cp.maxThreads {
		diff := cp.maxThreads - newSize
		for i := 0; i < diff; i++ {
			<-cp.available
		}
		cp.maxThreads = newSize
	}
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
