package main

import ()

func minCoins2(val int, coins []int) []int {
	if len(coins) == 0 {
		return []int{}
	}

	uniq := make(map[int]bool)
	clean := make([]int, 0)
	for _, c := range coins {
		if !uniq[c] {
			uniq[c] = true
			clean = append(clean, c)
		}
	}

	dp := make([]int, val+1)
	from := make([]int, val+1)
	for i := 1; i <= val; i++ {
		dp[i] = 1<<31 - 1
		for _, coin := range clean {
			if i >= coin && dp[i-coin]+1 < dp[i] {
				dp[i] = dp[i-coin] + 1
				from[i] = coin
			}
		}
	}

	if dp[val] == 1<<31-1 {
		return []int{}
	}

	res := []int{}
	for val > 0 {
		res = append(res, from[val])
		val -= from[val]
	}
	return res
}

func main() {
	ans := minCoins2(с)
}
