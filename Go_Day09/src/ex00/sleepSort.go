package main

import (
	"sync"
	"time"
)

func SleepSort(arr []int) <-chan int {
	out := make(chan int, len(arr))

	if len(arr) == 0 {
		close(out)
		return out
	}

	var wg sync.WaitGroup

	for _, n := range arr {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			time.Sleep(time.Duration(n) * 10 * time.Millisecond)
			out <- n
		}(n)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
