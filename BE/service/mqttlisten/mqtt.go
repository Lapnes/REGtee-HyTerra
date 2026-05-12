package mqttlisten

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	entity "backend/entity/db"
	"backend/repository/readings"
	"backend/repository/sensor"
)

type Service struct {
	repoSensor  sensor.Repositorier
	repoReading readings.Repositorier
}

type Servicer interface {
	ProcessMQTTMessage(topic string, payload []byte)
}

func NewMQTTService(repoSensor sensor.Repositorier, repoReading readings.Repositorier) *Service {
	return &Service{
		repoSensor:  repoSensor,
		repoReading: repoReading,
	}
}

// ProcessMQTTMessage handles incoming MQTT messages from ESP32
func (svc *Service) ProcessMQTTMessage(topic string, payload []byte) {
	log.Printf("[MQTT] Received on %s: %s", topic, string(payload))

	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		log.Printf("[MQTT] Failed to parse JSON: %v", err)
		return
	}

	// Extract device name
	deviceName, ok := data["device"].(string)
	if !ok || deviceName == "" {
		log.Println("[MQTT] Missing device field")
		return
	}

	// Find sensor by name — kalau nggak ada, auto-create
	sensorEntity, err := svc.repoSensor.GetByName(deviceName)
	if err != nil {
		log.Printf("[MQTT] Sensor not found: %s, auto-creating...", deviceName)

		newSensor := &entity.Sensor{
			Name:   deviceName,
			Area:   "Default Area",
			Status: "active",
		}
		if err := svc.repoSensor.Create(newSensor); err != nil {
			log.Printf("[MQTT] Failed to create sensor: %v", err)
			return
		}

		// Re-fetch sensor yang baru dibuat
		sensorEntity, err = svc.repoSensor.GetByName(deviceName)
		if err != nil {
			log.Printf("[MQTT] Failed to fetch new sensor: %v", err)
			return
		}
		log.Printf("[MQTT] Created new sensor: ID=%d, Name=%s", sensorEntity.ID, sensorEntity.Name)
	}

	// Extract moisture/humidity
	humidity := float32(0)
	if moisture, ok := data["moisture_percent"].(float64); ok {
		humidity = float32(moisture)
	} else if rawADC, ok := data["raw_adc"].(float64); ok {
		adcDry := 2924.0
		adcWet := 1170.0
		if d, ok := data["adc_dry"].(float64); ok {
			adcDry = d
		}
		if w, ok := data["adc_wet"].(float64); ok {
			adcWet = w
		}
		humidity = float32((1 - (rawADC-adcWet)/(adcDry-adcWet)) * 100)
		if humidity < 0 {
			humidity = 0
		}
		if humidity > 100 {
			humidity = 100
		}
	}

	// Extract pump_active
	pumpOn := false
	if pa, ok := data["pump_active"].(bool); ok {
		pumpOn = pa
	}

	// Extract timestamp
	timestamp := time.Now()
	if ts, ok := data["timestamp"].(string); ok {
		if parsed, err := time.Parse("2006-01-02 15:04:05", ts); err == nil {
			timestamp = parsed
		}
	}

	// Create reading
	reading := &entity.Readings{
		SensorID:  sensorEntity.ID,
		Type:      "soil_moisture",
		Humidity:  humidity,
		Timestamp: timestamp,
		PumpOn:    pumpOn,
	}

	if err := svc.repoReading.Create(reading); err != nil {
		log.Printf("[MQTT] Failed to save reading: %v", err)
		return
	}

	log.Printf("[MQTT] Saved reading: sensor_id=%d, humidity=%.1f%%, pump_on=%v",
		reading.SensorID, reading.Humidity, reading.PumpOn)
}

func parseUint(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 32)
	return uint(v)
}

func parseFloat(s string) float32 {
	v, _ := strconv.ParseFloat(s, 32)
	return float32(v)
}
