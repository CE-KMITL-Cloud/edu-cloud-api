package main

import (
	"fmt"
	"log"

	// libvirt_connect "github.com/edu-cloud-api/libvirt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {
	// Setup .env file
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	// Get value from .env
	// value, ok := viper.Get("test").(string)
	// if !ok {
	// 	log.Fatal("Invalid type assertion")
	// }

	// libvirt_connect.Connect_tcp()

	app := fiber.New()

	// Get /health-check
	app.Get("/", func(c *fiber.Ctx) error {
		msg := "✋ Healthy"
		return c.SendString(msg)
	})

	// Create /api group
	api := app.Group("/api")

	// GET /api/register
	api.Get("/*", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("✋ %s", c.Params("*"))
		return c.SendString(msg) // => ✋ register
	})

	log.Fatal(app.Listen(":3001"))
}
