package readings

import (
	httpEntity "backend/entity/http/readings"
	"backend/repository/readings"
)

type Service struct {
	repoReadings readings.Repository
}

func NewService(repoReadings readings.Repository) *Service {
	return &Service{
		repoReadings: repoReadings,
	}
}

type Servicer interface {
	GetReadings(req httpEntity.GetReadings) (resp httpEntity.Readings, err error)
}

func (svc Service) GetReadings(req httpEntity.GetReadings) (resp httpEntity.Readings, err error) {
	readingEntities, err := svc.repoReadings.GetReadings(req.Pgn)
	if err != nil {
		return httpEntity.Readings{}, err
	}
	err = resp.Build(readingEntities)
	return resp, nil
}
