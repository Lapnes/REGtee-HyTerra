package sensor

import entity "backend/entity/db"

type Esp struct {
	Id     int    `json:"sensor_id"`
	Name   string `json:"sensor_name"`
	Area   string `json:"sensor_area"`
	Status string `json:"sensor_status"`
}

type Esps struct {
	Esps []Esp `json:"sensor_lists"`
}

func (r *Esps) Build(in []entity.Sensor) error {
	for _, val := range in {
		temp := Esp{}
		temp.build(val)
		// ✅ FIX: Hapus impossible condition (build() nggak return error)
		r.Esps = append(r.Esps, temp)
	}
	return nil
}

func (r *Esp) build(in entity.Sensor) {
	r.Id = int(in.ID)
	r.Name = in.Name
	r.Area = in.Area
	r.Status = in.Status
}
