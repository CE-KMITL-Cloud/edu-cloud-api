// Package database - database's functions
package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type dsn struct {
	hostname string
	username string
	password string
	dbname   string
	port     string
}

// TableToMigrate - Migrate the schema for each table
type TableToMigrate struct {
	Name   string
	Schema interface{}
}

// DB - db's global variable
var DB *gorm.DB

// GetDSN - getting datasource
func GetDSN() string {
	datasource := dsn{
		hostname: config.GetFromENV("DB_HOST"),
		username: config.GetFromENV("DB_USER"),
		password: config.GetFromENV("DB_PASS"),
		dbname:   config.GetFromENV("DB_NAME"),
		port:     config.GetFromENV("DB_PORT"),
	}
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", datasource.hostname, datasource.username, datasource.password, datasource.dbname, datasource.port)
}

// Initialize - setting logger & running migrations
func Initialize() {
	// Set up the logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io.Writer
		logger.Config{
			SlowThreshold: time.Second,  // Show SQL queries slower than this threshold
			LogLevel:      logger.Error, // Log only error messages
			Colorful:      true,         // Use color output
		},
	)
	DB, _ = gorm.Open(postgres.Open(GetDSN()), &gorm.Config{
		Logger: newLogger,
	})
	RunMigrations()
}

// RunMigrations - running migrations function
func RunMigrations() {
	tablesToMigrate := []TableToMigrate{
		{"admin", &model.User{}},
		{"student", &model.User{}},
		{"faculty", &model.User{}},
		{"instance", &model.Instance{}},
		{"instance_limit", &model.InstanceLimit{}},
		{"pool", &model.Pool{}},
		{"sizing", &model.Sizing{}},
		// {"proxy", &Proxy{}},
		// {"proxy_key", &ProxyKey{}},
	}
	log.Println("Running migrations ...")
	for _, table := range tablesToMigrate {
		err := DB.Table(table.Name).AutoMigrate(table.Schema)
		if err != nil {
			panic(fmt.Sprintf("migration of %s table failed: %v", table.Name, err))
		}
	}
	fmt.Println("Successfully running migrations")
}
