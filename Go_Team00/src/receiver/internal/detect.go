package receiver

import (
	"fmt"
	"math"
)

type Detector struct {
	count       int
	mean        float64
	m2          float64
	k           float64
	readyCount  int
	anomalyMode bool
}

func NewDetector(k float64, readyCount int) *Detector {
	return &Detector{
		k:          k,
		readyCount: readyCount,
	}
}

func (d *Detector) Process(value float64) {
	d.count++

	// Welford update
	delta := value - d.mean
	d.mean += delta / float64(d.count)
	delta2 := value - d.mean
	d.m2 += delta * delta2

	// включаем аномалии
	if !d.anomalyMode && d.count >= d.readyCount {
		d.anomalyMode = true
		fmt.Printf("[INFO] Detector is ready after %d samples.\n", d.count)
	}
}

// Std возвращает текущее стандартное отклонение
func (d *Detector) Std() float64 {
	if d.count < 2 {
		return 0
	}
	variance := d.m2 / float64(d.count-1)
	return math.Sqrt(variance)
}

// Mean возвращает текущее среднее
func (d *Detector) Mean() float64 {
	return d.mean
}

// IsReady возвращает true если уже можно искать аномалии
func (d *Detector) IsReady() bool {
	return d.anomalyMode
}

// IsAnomaly проверяет значение на аномалию
func (d *Detector) IsAnomaly(value float64) bool {
	std := d.Std()
	if std == 0 {
		return false
	}
	diff := math.Abs(value - d.mean)
	return diff > d.k*std
}
