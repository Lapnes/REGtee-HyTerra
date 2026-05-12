package entity

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	Password  string    `json:"-" gorm:"size:255;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type Sensor struct {
	ID        uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string     `json:"name" gorm:"size:100;not null"`
	Area      string     `json:"location" gorm:"size:200"`
	Status    string     `json:"status" gorm:"size:20;default:'active'"`
	Readings  []Readings `json:"readings,omitempty" gorm:"foreignKey:SensorID;references:ID"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (Sensor) TableName() string {
	return "sensors"
}

type Readings struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	SensorID  uint      `json:"sensor_id" gorm:"not null;index"`
	Type      string    `json:"type" gorm:"size:50"`
	Sensor    Sensor    `json:"sensor,omitempty" gorm:"foreignKey:SensorID;references:ID"`
	Humidity  float32   `json:"humidity" gorm:"type:decimal(10,2)"`
	Timestamp time.Time `json:"timestamp" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// ✅ FIX: Status di-rename jadi PumpOn, default false
	PumpOn bool `json:"pump_on" gorm:"default:false"`
}

func (Readings) TableName() string {
	return "readings"
}
