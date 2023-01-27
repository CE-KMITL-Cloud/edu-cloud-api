// Package main ...
package main

import (
	"log"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
)

func main() {
	// database.ConnectDb()
	app := fiber.New()
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: config.GetFromENV("ENCRYPT_KEY"),
	}))
	router.SetupRoutes(app)
	log.Fatal(app.Listen(":3001"))
}
