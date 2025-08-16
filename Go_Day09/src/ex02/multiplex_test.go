package main

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

func makeChan(values ...interface{}) chan interface{} {
	ch := make(chan interface{}, len(values))
	for _, v := range values {
		ch <- v
	}
	close(ch)
	return ch
}

func TestMultiplex_AllValuesReceived(t *testing.T) {
	ch1 := makeChan("a", "b")
	ch2 := makeChan(1, 2, 3)
	ch3 := makeChan(true, false)

	expected := []interface{}{"a", "b", 1, 2, 3, true, false}
	received := []interface{}{}

	out := multiplex(ch1, ch2, ch3)

	// Читаем все значения
	for val := range out {
		received = append(received, val)
	}

	// Сортируем для независимости от порядка
	sort.Slice(received, func(i, j int) bool { return toString(received[i]) < toString(received[j]) })
	sort.Slice(expected, func(i, j int) bool { return toString(expected[i]) < toString(expected[j]) })

	if !reflect.DeepEqual(received, expected) {
		t.Errorf("Expected %v, got %v", expected, received)
	}
}

// Переводим любые значения в строку для сортировки
func toString(v interface{}) string {
	return reflect.ValueOf(v).String()
}

func TestMultiplex_EmptyInput(t *testing.T) {
	out := multiplex()

	select {
	case _, ok := <-out:
		if ok {
			t.Errorf("Expected closed output channel for empty input")
		}
	case <-time.After(500 * time.Millisecond):
		t.Errorf("Timeout: output channel was not closed")
	}
}
