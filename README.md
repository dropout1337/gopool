# gopool

**gopool** is a Go package for managing goroutines and simplifying concurrency in your Go applications.

## Features

- Efficient goroutine pooling for managing parallel tasks.
- Dynamic resizing to adapt to varying workloads.
- Control the maximum number of concurrent goroutines.
- Simplify parallel task execution.

## Installation

To use **gopool** in your Go project, you can simply run:

```shell
go get github.com/dropout1337/gopool
```

## Example
```go
package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	pool := New(15)

	for i := 0; i < 25; i++ {
		pool.Wait()

		go func(i int) {
			defer pool.Done()

			log.Println("Task started:", i)
			time.Sleep(time.Millisecond * 100)
		}(i)
	}

	pool.WaitUntilDone()
	fmt.Println("Resizing pool to 1 threads")
	pool.ResizePool(1)

	for i := 0; i < 10; i++ {
		pool.Wait()

		go func(i int) {
			defer pool.Done()

			log.Println("Task started:", i)
			time.Sleep(time.Millisecond * 100)
		}(i)
	}

	pool.WaitUntilDone()

	if pool.GetMaxThreads() != 1 {
		log.Fatalf("Expected maxThreads to be 1, got %d", pool.GetMaxThreads())
	}

	if pool.GetCurrentThreadCount() != 0 {
		log.Fatalf("Expected current thread count to be 0, got %d", pool.GetCurrentThreadCount())
	}

	if !pool.IsAvailable() {
		log.Fatalf("Expected pool to be available after resizing")
	}
}
```
