package main

import (
	"testing"
)

func TestGetElement_ValidIndex(t *testing.T) {
	arr := []int{10, 20, 30, 40, 50}

	tests := []struct {
		index    int
		expected int
	}{
		{0, 10},
		{1, 20},
		{2, 30},
		{4, 50},
	}

	for _, test := range tests {
		val, err := getElement(arr, test.index)
		if err != nil {
			t.Errorf("expected no error for index %d, got %v", test.index, err)
		}
		if val != test.expected {
			t.Errorf("expected value %d for index %d, got %d", test.expected, test.index, val)
		}
	}
}

func TestGetElement_NegativeIndex(t *testing.T) {
	arr := []int{1, 2, 3}
	_, err := getElement(arr, -1)
	if err == nil {
		t.Error("expected error for negative index, got nil")
	}
}

func TestGetElement_IndexOutOfBounds(t *testing.T) {
	arr := []int{1, 2, 3}
	_, err := getElement(arr, 3)
	if err == nil {
		t.Error("expected error for out-of-bounds index, got nil")
	}
}

func TestGetElement_EmptySlice(t *testing.T) {
	arr := []int{}
	_, err := getElement(arr, 0)
	if err == nil {
		t.Error("expected error for empty slice, got nil")
	}
}
