package knapsack

type Present struct {
	Value int
	Size  int
}

func grabPresents(presents []Present, capacity int) []Present {
	n := len(presents)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, capacity+1)
	}

	for i := 1; i <= n; i++ {
		for w := 0; w <= capacity; w++ {
			if presents[i-1].Size > w {
				dp[i][w] = dp[i-1][w]
			} else {
				dp[i][w] = max(
					dp[i-1][w],
					dp[i-1][w-presents[i-1].Size]+presents[i-1].Value,
				)
			}
		}
	}

	res := []Present{}
	w := capacity
	for i := n; i > 0; i-- {
		if dp[i][w] != dp[i-1][w] {
			res = append(res, presents[i-1])
			w -= presents[i-1].Size
		}
	}

	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
