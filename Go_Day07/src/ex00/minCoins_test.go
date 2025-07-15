package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestMinCoins(t *testing.T) {
	tests := []struct {
		name     string
		val      int
		coins    []int
		expected []int
	}{
		{"Basic greedy OK", 13, []int{1, 5, 10}, []int{10, 1, 1, 1}},
		{"Greedy fails", 6, []int{1, 3, 4}, []int{3, 3}}, // minCoins fails here
		{"Empty coins", 10, []int{}, []int{}},
		{"Duplicates and unsorted", 6, []int{4, 3, 1, 1}, []int{3, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := minCoins2(tt.val, tt.coins)
			if !equalUnordered(got, tt.expected) {
				t.Errorf("minCoins2(%v, %v) = %v; want %v", tt.val, tt.coins, got, tt.expected)
			}
		})
	}
}

func equalUnordered(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	aCopy := append([]int(nil), a...)
	bCopy := append([]int(nil), b...)
	sort.Ints(aCopy)
	sort.Ints(bCopy)
	return reflect.DeepEqual(aCopy, bCopy)
}
