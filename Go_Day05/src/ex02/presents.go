package presents

import (
	"container/heap"
	"fmt"
)

type Present struct {
	Value int
	Size  int
}

type PresentHeap []Present

func (h PresentHeap) Len() int { return len(h) }

func (h PresentHeap) Less(i, j int) bool {
	if h[i].Value != h[j].Value {
		return h[i].Value > h[j].Value
	}
	return h[i].Size < h[j].Size
}
func (h PresentHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *PresentHeap) Push(x any) {
	*h = append(*h, x.(Present))
}
func (h *PresentHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func GetNCoolestPresents(presents []Present, n int) ([]Present, error) {
	if n < 0 || n > len(presents) {
		return nil, fmt.Errorf("n must be from 0 to %d, got %d", len(presents), n)
	}

	h := PresentHeap(presents)
	heap.Init(&h)

	result := make([]Present, n)
	for i := 0; i < n; i++ {
		result[i] = heap.Pop(&h).(Present)
	}
	return result, nil
}
