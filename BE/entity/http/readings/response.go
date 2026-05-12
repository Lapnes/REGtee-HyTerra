package readings

import (
	entity "backend/entity/db"
	"time"
)

type Readings struct {
	Items []Reading `json:"readings"`
}

type Reading struct {
	SensorID              uint      `json:"sensor_id"`
	Humidity              float32   `json:"humidity"`
	IntervalSiram         string    `json:"interval_siram"`
	WaktuSiramSelanjutnya time.Time `json:"waktu_siram_selanjutnya"`
	Timestamp             time.Time `json:"timestamp"`
	PumpOn                bool      `json:"pump_on"`
}

func (r *Readings) Build(in []entity.Readings) error {
	for _, val := range in {
		var item Reading
		item.build(val)
		r.Items = append(r.Items, item)
	}
	return nil
}

func (r *Reading) build(in entity.Readings) {
	r.SensorID = in.SensorID
	r.Humidity = in.Humidity
	r.Timestamp = in.Timestamp
	// ✅ FIX: Ambil PumpOn dari DB entity, bukan Status
	r.PumpOn = in.PumpOn
	// Sementara hardcode interval & next watering
	r.IntervalSiram = "30m"
	r.WaktuSiramSelanjutnya = in.Timestamp.Add(30 * time.Minute)
}
