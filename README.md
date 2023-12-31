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

## Usage
Here's a quick example of how to use **gopool**:
```go
package main

import (
    "fmt"
    "github.com/dropout1337/gopool"
)

func main() {
    // Create a new gopool with a maximum of 5 threads
    pool := gopool.New(5)

    for i := 1; i <= 10; i++ {
        pool.Execute(func() {
            fmt.Printf("Task %d started\n", i)
            // Simulate some work
            fmt.Printf("Task %d completed\n", i)
        })
    }

    pool.WaitUntilDone() // Wait for all tasks to finish
}
```

