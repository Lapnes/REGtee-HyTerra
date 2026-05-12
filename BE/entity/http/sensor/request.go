package sensor

import "backend/util/pagination"

type GetSensor struct {
	Pgn pagination.Pagination `form:"pgn"`
}

type UpdateSensorStatus struct {
	Id     int    `json:"sensor_id"`
	Status string `json:"sensor_status"`
}
