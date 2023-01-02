// Package main ...
package main

import (
	"log"

	"github.com/edu-cloud-api/database"
	"github.com/edu-cloud-api/router"
	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDb()
	app := fiber.New()

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":3001"))
}
