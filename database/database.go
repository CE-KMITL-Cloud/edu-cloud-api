// Package database - database's functions
package database

import (
	"fmt"
	"log"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

// connectDb
func ConnectDb() {
	items := []string{"dbHost", "dbUser", "dbPass", "dbName", "dbPort"}
	dbItems := config.GetListFromENV(items)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbItems[0], dbItems[1], dbItems[2], dbItems[3], dbItems[4])

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("running migrations")
	db.AutoMigrate(&model.User{})

	DB = Dbinstance{
		Db: db,
	}
}
