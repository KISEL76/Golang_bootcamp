package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Тест успешного получения данных с тестового сервера
func TestCrawlWebSuccess(t *testing.T) {
	// Создаем мок-сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}))
	defer ts.Close()

	ctx := context.Background()
	urls := make(chan string, 1)
	urls <- ts.URL
	close(urls)

	results := crawlWeb(ctx, urls)

	var resultCount int
	for body := range results {
		resultCount++
		if !strings.Contains(body, "hello world") {
			t.Errorf("Expected body to contain 'hello world', got: %s", body)
		}
	}

	if resultCount != 1 {
		t.Errorf("Expected 1 result, got %d", resultCount)
	}
}

// Тест обработки нескольких URL
func TestCrawlWebMultipleURLs(t *testing.T) {
	const expectedResponses = 5

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))
	defer ts.Close()

	ctx := context.Background()
	urls := make(chan string, expectedResponses)

	for i := 0; i < expectedResponses; i++ {
		urls <- ts.URL
	}
	close(urls)

	results := crawlWeb(ctx, urls)

	count := 0
	for res := range results {
		count++
		if res != "OK" {
			t.Errorf("Expected 'OK', got %s", res)
		}
	}

	if count != expectedResponses {
		t.Errorf("Expected %d results, got %d", expectedResponses, count)
	}
}

// Тест отмены через context.WithCancel
func TestCrawlWebCancel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // Эмуляция долгого ответа
		w.Write([]byte("delayed"))
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	urls := make(chan string, 3)

	for i := 0; i < 3; i++ {
		urls <- ts.URL
	}
	close(urls)

	// Отменяем контекст почти сразу
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	results := crawlWeb(ctx, urls)

	count := 0
	for range results {
		count++
	}

	if count >= 3 {
		t.Errorf("Expected fewer than 3 results due to cancellation, got %d", count)
	}
}
