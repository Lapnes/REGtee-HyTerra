package user

import (
	"backend/database"
	entity "backend/entity/db"

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
	GetUserById(id int) (resp entity.User, err error)
	GetUserByName(name string) (resp entity.User, err error)
}

func (repo *Repository) GetUserById(id int) (resp entity.User, err error) {
	err = repo.slave.Model(&entity.User{}).Where("id = ?", id).First(&resp).Error
	return
}

func (repo *Repository) GetUserByName(name string) (resp entity.User, err error) {
	err = repo.slave.Model(&entity.User{}).Where("name = ?", name).First(&resp).Error
	return
}
