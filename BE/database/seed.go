package database

import (
	"log"
	"time"

	models "backend/entity/db"
	hash "backend/util/cryptography"

	"gorm.io/gorm"
)

// SeedDatabase seeds the database with initial data
func SeedDatabase(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	// Seed Users
	if err := seedUsers(db); err != nil {
		return err
	}

	// Seed Sensors
	if err := seedSensors(db); err != nil {
		return err
	}

	// Seed Readings
	if err := seedReadings(db); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

func seedUsers(db *gorm.DB) error {
	log.Println("Seeding users...")

	// Define users with plain text passwords
	usersData := []struct {
		name     string
		password string
	}{
		{
			name:     "admin",
			password: "password123",
		},
		{
			name:     "John Doe",
			password: "password123",
		},
	}

	for _, userData := range usersData {
		var count int64
		db.Model(&models.User{}).Where("name = ?", userData.name).Count(&count)
		if count == 0 {
			user := models.User{
				Name:     userData.name,
				Password: hash.HashPassword(userData.password), // Hash password before saving
			}
			if err := db.Create(&user).Error; err != nil {
				log.Printf("Error seeding user %s: %v", userData.name, err)
				return err
			}
			log.Printf("Seeded user: %s", userData.name)
		}
	}

	return nil
}

func seedSensors(db *gorm.DB) error {
	log.Println("Seeding sensors...")

	sensors := []models.Sensor{
		{
			Name:   "Sensor 1",
			Area:   "Taman Belakang",
			Status: "active",
		},
	}

	for _, sensor := range sensors {
		var count int64
		db.Model(&models.Sensor{}).Where("name = ?", sensor.Name).Count(&count)
		if count == 0 {
			if err := db.Create(&sensor).Error; err != nil {
				log.Printf("Error seeding sensor %s: %v", sensor.Name, err)
				return err
			}
			log.Printf("Seeded sensor: %s", sensor.Name)
		}
	}

	return nil
}

func seedReadings(db *gorm.DB) error {
	log.Println("Seeding readings...")

	readings := []models.Readings{
		{
			SensorID:  1,
			Type:      "humidity",
			Humidity:  65.5,
			Timestamp: time.Now().Add(-24 * time.Hour),
			PumpOn:    false,
		},
	}

	for _, reading := range readings {
		var count int64
		db.Model(&models.Readings{}).
			Where("sensor_id = ? AND timestamp = ?", reading.SensorID, reading.Timestamp).
			Count(&count)
		if count == 0 {
			if err := db.Create(&reading).Error; err != nil {
				log.Printf("Error seeding reading for sensor %d: %v", reading.SensorID, err)
				return err
			}
			log.Printf("Seeded reading for sensor %d: %.1f%%", reading.SensorID, reading.Humidity)
		}
	}

	return nil
}
