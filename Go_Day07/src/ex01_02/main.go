// Пакет ex01 реализует два алгоритма размена монет.
//
// Реализовано:
//   - minCoins: жадный алгоритм (быстрый, но не всегда оптимальный).
//   - minCoins2: алгоритм динамического программирования, который всегда находит минимальное число монет.
package main

// minCoins реализует жадный алгоритм размена монет:
// он всегда берёт самую большую подходящую монету.
//
// Быстрый, но не гарантирует минимальное количество монет.
func minCoins(val int, coins []int) []int {
	res := make([]int, 0)
	i := len(coins) - 1
	for i >= 0 {
		for val >= coins[i] {
			val -= coins[i]
			res = append(res, coins[i])
		}
		i -= 1
	}
	return res
}

// minCoins2 возвращает минимальное количество монет, сумма которых равна val,
// используя алгоритм динамического программирования (ДП).
//
// Отличия от жадного подхода:
//   - сначала удаляются дубликаты номиналов;
//   - используется ДП-массив dp[i], где i — сумма, а значение — минимальное количество монет;
//   - массив from позволяет восстановить, какие монеты использовались.
//
// Чтобы сгенерировать документацию:
//  1. Установи golds: go install go101.org/golds@latest
//  2. Выполни команду: golds -dir=docs -output=docs
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
