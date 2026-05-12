package sensor

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
	Create(sensor *entity.Sensor) error
	GetByName(name string) (entity.Sensor, error)
	GetByID(id uint) (entity.Sensor, error)
	GetAll() ([]entity.Sensor, error)
	UpdateSensorStatus(id int, status string) error
	GetSensors(pgn pagination.Pagination) ([]entity.Sensor, error)
}

func (r *Repository) Create(sensor *entity.Sensor) error {
	return r.master.Create(sensor).Error
}

func (r *Repository) GetByName(name string) (entity.Sensor, error) {
	var s entity.Sensor
	err := r.slave.Where("name = ?", name).First(&s).Error
	return s, err
}

func (r *Repository) GetByID(id uint) (entity.Sensor, error) {
	var s entity.Sensor
	err := r.slave.First(&s, id).Error
	return s, err
}

func (r *Repository) GetAll() ([]entity.Sensor, error) {
	var sensors []entity.Sensor
	err := r.slave.Find(&sensors).Error
	return sensors, err
}

func (repo *Repository) GetSensors(pgn pagination.Pagination) ([]entity.Sensor, error) {
	var resp []entity.Sensor
	err := repo.slave.
		Model(&entity.Sensor{}).
		Limit(pgn.Limit).
		Offset(pgn.Offset).
		Find(&resp).Error

	return resp, err
}

func (repo *Repository) UpdateSensorStatus(id int, status string) error {
	err := repo.master.
		Model(&entity.Sensor{}).
		Where("id = ?", id).
		Update("status", status).
		Error

	return err
}
