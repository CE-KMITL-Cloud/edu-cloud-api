package main

import (
	"fmt"
	"log"

	libvirt_connection "github.com/edu-cloud-api/libvirt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// Get any item from .env
func getFromENV(item string) string {
	value, ok := viper.Get(item).(string)
	if !ok {
		log.Fatalf("Error while getting item : %s", value)
	}
	return value
}

func main() {
	// Setup .env file
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	// Collecting variables from .env
	AUTHNAME := getFromENV("AUTHNAME")
	PASSPHASE := getFromENV("PASSPHASE")

	// Start connection with libvirt
	conn := libvirt_connection.TCP_Connect(AUTHNAME, PASSPHASE)
	log.Println(conn)

	// Need to close connection after process done
	defer conn.Close()

	app := fiber.New()

	// Get /health-check
	app.Get("/", func(c *fiber.Ctx) error {
		msg := "✋ Healthy"
		return c.SendString(msg)
	})

	api := app.Group("/api")

	// GET /api/register
	api.Get("/*", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("✋ %s", c.Params("*"))
		return c.SendString(msg) // => ✋ register
	})

	log.Fatal(app.Listen(":3001"))
}
