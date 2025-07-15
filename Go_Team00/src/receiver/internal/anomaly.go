package receiver

import (
	"math"
)

func IsAnomaly(value, mean, std, k float64) bool {
	return math.Abs(value-mean) > k*std
}
