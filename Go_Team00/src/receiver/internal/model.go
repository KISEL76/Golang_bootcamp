package receiver

import (
	"time"
)

type Anomaly struct {
	ID        uint `gorm:"primaryKey"`
	SessionID string
	Frequency float64
	Timestamp int64
	CreatedAt time.Time
}
