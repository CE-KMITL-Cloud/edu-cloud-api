// Package database - database's functions
package database

import (
	"fmt"
	"log"

	"github.com/edu-cloud-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Dbinstance - db's struct
type Dbinstance struct {
	Db *gorm.DB
}

// DB - db's global variable
var DB Dbinstance

// ConnectDb - create connection to db
func ConnectDb() {
	items := []string{"DB_HOST", "DB_USER", "DB_PASS", "DB_NAME", "DB_PORT"}
	dbItems := config.GetListFromENV(items)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbItems[0], dbItems[1], dbItems[2], dbItems[3], dbItems[4])

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	log.Println("Connected to Database successfully.")
	db.Logger = logger.Default.LogMode(logger.Info)
	// log.Println("running migrations")
	// db.AutoMigrate(&model.User{})

	DB = Dbinstance{
		Db: db,
	}
}
