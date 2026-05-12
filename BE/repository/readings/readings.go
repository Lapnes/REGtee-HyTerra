package readings

import (
	"backend/database"
	entity "backend/entity/db"
	"backend/util/pagination"

	"gorm.io/gorm"
)

type Repository struct {
	master *gorm.DB
	slave  *gorm.DB
}

func NewRepository(hosts database.DB) *Repository {
	return &Repository{
		master: hosts.Master,
		slave:  hosts.Slave,
	}
}

type Repositorier interface {
	Create(reading *entity.Readings) error
	GetReadings(pgn pagination.Pagination) (resp []entity.Readings, err error)
}

func (repo *Repository) GetReadings(pgn pagination.Pagination) (resp []entity.Readings, err error) {
	// ✅ FIX: ORDER BY timestamp DESC supaya data terbaru di atas
	err = repo.slave.Model(&entity.Readings{}).
		Order("timestamp DESC").
		Limit(pgn.Limit).
		Offset(pgn.Offset).
		Find(&resp).
		Error

	return
}

func (repo *Repository) Create(reading *entity.Readings) error {
	return repo.master.Create(reading).Error
}
