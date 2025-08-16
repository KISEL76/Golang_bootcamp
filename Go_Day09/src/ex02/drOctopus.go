package main

import "sync"

func multiplex(channels ...chan interface{}) chan interface{} {
	out := make(chan interface{})
	var wg sync.WaitGroup

	copy := func(ch <-chan interface{}) {
		defer wg.Done()
		for val := range ch {
			out <- val
		}
	}

	wg.Add(len(channels))
	for _, ch := range channels {
		go copy(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
