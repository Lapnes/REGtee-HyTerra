package sensor

import (
	httpEntity "backend/entity/http/sensor"
	"backend/repository/sensor"
	"backend/service/mqttpub"
)

type Service struct {
	repoSensor sensor.Repositorier
	pubService mqttpub.CommandService
}

func NewService(
	repoSensor sensor.Repositorier,
	pubService mqttpub.CommandService,
) *Service {
	return &Service{
		repoSensor: repoSensor,
		pubService: pubService,
	}
}

type Servicer interface {
	GetSensor(req httpEntity.GetSensor) (resp httpEntity.Esps, err error)
	UpdateSensorStatus(req httpEntity.UpdateSensorStatus) error
}

func (svc Service) GetSensor(req httpEntity.GetSensor) (resp httpEntity.Esps, err error) {
	entities, err := svc.repoSensor.GetSensors(req.Pgn)
	if err != nil {
		return httpEntity.Esps{}, err
	}
	if err = resp.Build(entities); err != nil {
		return httpEntity.Esps{}, err
	}
	return resp, nil
}

func (svc Service) UpdateSensorStatus(req httpEntity.UpdateSensorStatus) error {
	if err := svc.repoSensor.UpdateSensorStatus(req.Id, req.Status); err != nil {
		return err
	}
	switch req.Status {
	case "active":
		return svc.pubService.SetPump("1", true)
	case "inactive":
		return svc.pubService.SetPump("1", false)
	}
	return nil
}
