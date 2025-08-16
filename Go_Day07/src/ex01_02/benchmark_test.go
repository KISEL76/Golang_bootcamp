package main

import (
	"testing"
)

func BenchmarkMinCoins2(b *testing.B) {
	coins := []int{1, 2, 5, 10, 25, 50}
	val := 637
	for i := 0; i < b.N; i++ {
		minCoins2(val, coins)
	}
}

func BenchmarkMinCoins(b *testing.B) {
	coins := []int{1, 2, 5, 10, 25, 50}
	val := 637
	for i := 0; i < b.N; i++ {
		minCoins(val, coins)
	}
}
