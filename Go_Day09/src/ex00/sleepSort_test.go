package main

import (
	"reflect"
	"sort"
	"testing"
)

func readAll(ch <-chan int) []int {
	var result []int
	for v := range ch {
		result = append(result, v)
	}
	return result
}

func TestSleepSort_Basic(t *testing.T) {
	input := []int{3, 1, 2}
	expected := []int{1, 2, 3}

	out := SleepSort(input)
	result := readAll(out)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestSleepSort_Basic2(t *testing.T) {
	input := []int{30, 10, 20, 15, 50, 45}
	expected := []int{10, 15, 20, 30, 45, 50}

	out := SleepSort(input)
	result := readAll(out)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestSleepSort_AlreadySorted(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := []int{1, 2, 3, 4}

	out := SleepSort(input)
	result := readAll(out)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Already sorted: expected %v, got %v", expected, result)
	}
}

func TestSleepSort_ReverseOrder(t *testing.T) {
	input := []int{5, 4, 3, 2, 1}
	expected := []int{1, 2, 3, 4, 5}

	out := SleepSort(input)
	result := readAll(out)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Reverse order: expected %v, got %v", expected, result)
	}
}

func TestSleepSort_Duplicates(t *testing.T) {
	input := []int{2, 1, 2, 3, 1}
	expected := []int{1, 1, 2, 2, 3}

	out := SleepSort(input)
	result := readAll(out)

	sort.Ints(result)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Duplicates: expected %v, got %v", expected, result)
	}
}

func TestSleepSort_OneElement(t *testing.T) {
	input := []int{42}
	expected := []int{42}

	out := SleepSort(input)
	result := readAll(out)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Single element: expected %v, got %v", expected, result)
	}
}
