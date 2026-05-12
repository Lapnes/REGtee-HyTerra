package database

import (
	variable "backend/constant"
	models "backend/entity/db"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbMaster *gorm.DB
var dbSlave *gorm.DB

type DB struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Sensor{},
		&models.Readings{},
	)
}

func dbInit(server string) *gorm.DB {
	var host string
	var port string
	var username string
	var password string
	var dbName string

	switch server {
	case variable.DB_MASTER:
		host = variable.DB_HOST
		port = variable.DB_PORT
		username = variable.DB_USER_MASTER
		password = variable.DB_PASSWORD_MASTER
		dbName = variable.DB_NAME
	case variable.DB_SLAVE:
		host = variable.DB_HOST
		port = variable.DB_PORT
		username = variable.DB_USER_SLAVE
		password = variable.DB_PASSWORD_SLAVE
		dbName = variable.DB_NAME
	}
	postgresCon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		host,
		port,
		username,
		password,
		dbName,
	)
	DB, err := gorm.Open(postgres.Open(postgresCon), &gorm.Config{
		QueryFields: true,
		NowFunc: func() time.Time {
			return time.Now()
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Cant't Connect to db")
	}

	// ✅ FIX: Hanya migrate di MASTER
	if server == variable.DB_MASTER {
		err = autoMigrate(DB)
		if err != nil {
			log.Fatal("Migrate failed")
		}
	}

	connConfDB, err := DB.DB()
	if err != nil {
		panic(fmt.Sprintf("failed to get sql db configuration: %v", err))
	}
	if connConfDB == nil {
		panic(fmt.Sprintf("failed connect: %v", err))
	}

	log.Printf("successfully established connection to %v", server)

	return DB
}

func DBMaster() *gorm.DB {
	if dbMaster == nil {
		dbMaster = dbInit(variable.DB_MASTER)
	}
	return dbMaster
}

func DBSlave() *gorm.DB {
	if dbSlave == nil {
		dbSlave = dbInit(variable.DB_SLAVE)
	}
	return dbSlave
}
