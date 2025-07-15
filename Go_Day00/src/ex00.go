package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	var numbers []int
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Введите целые числа построчно. Завершить ввод нужно комбинацией клавиш \"Ctrl + D\".")
	for scanner.Scan() {
		line := scanner.Text()
		number, err := strconv.Atoi(line)

		if err != nil {
			fmt.Println("Ошибка: Неверный формат числа. Попробуйте снова.")
			continue
		} else if number > 100000 || number < -100000 {
			fmt.Println("Ошибка: Число выходит за границы [-100000; 100000]. Попробуйте снова.")
			continue
		}
		numbers = append(numbers, number)
	}

	if len(numbers) == 0 {
		fmt.Println("Вы не ввели ни одного числа.")
		return
	}

	CalculateAndPrintMetrics(numbers)
}

func CalculateAndPrintMetrics(numbers []int) {
	sort.Ints(numbers)
	PrintMetrics(numbers)
}

func FindMedian(numbers []int) float32 {
	length := len(numbers)

	if length%2 == 0 {
		return (float32(numbers[length/2-1]) + float32(numbers[length/2])) / 2
	} else {
		return float32(numbers[length/2])
	}
}

func FindMean(numbers []int) float32 {
	var sum int

	for _, number := range numbers {
		sum += number
	}
	return float32(sum) / float32(len(numbers))
}

func FindMode(numbers []int) int {
	frequency := make(map[int]int)

	for _, number := range numbers {
		frequency[number]++
	}

	mostFrequentNum := 0
	maxFreq := 0

	for num, freq := range frequency {
		if freq > maxFreq {
			maxFreq = freq
			mostFrequentNum = num
		}
	}

	return mostFrequentNum
}

func FindSD(numbers []int, mean float64) float64 {
	var sumSquares float64

	for _, number := range numbers {
		diff := float64(number) - mean
		sumSquares += diff * diff
	}
	return math.Sqrt(sumSquares / float64(len(numbers)))
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func PrintMetrics(numbers []int) {
	mean := FindMean(numbers)

	fmt.Printf("Введите метрики для вывода (mean, median, mode, sd) через запятую: ")
	input := readInput()
	metricsToPrint := strings.Split(input, ",")

	if strings.TrimSpace(input) == "" {
		median := FindMedian(numbers)
		mode := FindMode(numbers)
		SD := FindSD(numbers, float64(mean))

		fmt.Printf("Mean: %.2f\n", mean)
		fmt.Printf("Median: %.2f\n", median)
		fmt.Printf("Mode: %d\n", mode)
		fmt.Printf("SD: %.2f\n", SD)
		return
	}

	for _, metric := range metricsToPrint {
		metric = strings.TrimSpace(strings.ToLower(metric))

		switch metric {
		case "sd":
			SD := FindSD(numbers, float64(mean))
			fmt.Printf("SD: %.2f\n", SD)
		case "mean":
			fmt.Printf("Mean: %.2f\n", mean)
		case "mode":
			mode := FindMode(numbers)
			fmt.Printf("Mode: %d\n", mode)
		case "median":
			median := FindMedian(numbers)
			fmt.Printf("Median: %.2f\n", median)
		default:
			if metric != "" {
				fmt.Printf("Ошибка: Неизвестная метрика '%s'. Доступные метрики: sd, mean, mode, median.\n", metric)
			}
		}
	}
}
