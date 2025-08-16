package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

const maxWorkers = 8

func crawlWeb(ctx context.Context, urls <-chan string) <-chan string {
	results := make(chan string)
	semaphore := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	for url := range urls {
		select {
		case <-ctx.Done():
			continue
		default:
			semaphore <- struct{}{}
			wg.Add(1)

			go func(u string) {
				defer wg.Done()
				defer func() {
					<-semaphore
				}()

				req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
				if err != nil {
					fmt.Println("Ошибка создания запроса")
					return
				}

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					fmt.Println("HTTP ошибка: ", err)
					return
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("Ошибка чтения тела: ", err)
					return
				}

				select {
				case results <- string(body):
				case <-ctx.Done():
					return
				}
			}(url)
		}
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
