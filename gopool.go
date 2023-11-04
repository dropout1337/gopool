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
	available := make(chan struct{}, maxThreads)
	for i := 0; i < maxThreads; i++ {
		available <- struct{}{}
	}

	return &ConcurrencyPool{
		maxThreads: maxThreads,
		available:  available,
	}
}

func (cp *ConcurrencyPool) Wait() {
	<-cp.available
	cp.wg.Add(1)
}

func (cp *ConcurrencyPool) Done() {
	cp.available <- struct{}{}
	cp.wg.Done()
}

func (cp *ConcurrencyPool) WaitUntilDone() {
	cp.wg.Wait()
}

func (cp *ConcurrencyPool) Execute(fn func()) {
	cp.Wait()

	go func() {
		defer cp.Done()
		fn()
	}()
}

func (cp *ConcurrencyPool) ResizePool(newSize int) error {
	cp.wg.Wait() // Ensure all goroutines have finished

	currentSize := len(cp.available)
	if newSize > currentSize {
		for i := currentSize; i < newSize; i++ {
			cp.available <- struct{}{}
		}
	} else if newSize < currentSize {
		return fmt.Errorf("cannot resize pool to a smaller size than the number of active goroutines")
	}

	cp.maxThreads = newSize
	return nil
}

func (cp *ConcurrencyPool) IsAvailable() bool {
	return len(cp.available) > 0
}
