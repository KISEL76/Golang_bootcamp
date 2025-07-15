package receiver

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(dsn string) *Database {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[DB] Failed to connect: %v", err)
	}
	log.Println("[DB] Connected")

	if err := db.AutoMigrate(&Anomaly{}); err != nil {
		log.Fatalf("[DB] AutoMigrate failed: %v", err)
	}

	return &Database{DB: db}
}

// SaveAnomaly сохраняет запись
func (d *Database) SaveAnomaly(anomaly *Anomaly) {
	if err := d.DB.Create(anomaly).Error; err != nil {
		log.Printf("[DB] Failed to save anomaly: %v", err)
	}
}
