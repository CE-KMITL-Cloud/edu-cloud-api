package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	api := app.Group("/api") // /api

	app.Get("/", func(c *fiber.Ctx) error {
		msg := "✋ Healthy"
		return c.SendString(msg)
	}).Name("health-check")

	// GET /api/register
	api.Get("/*", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("✋ %s", c.Params("*"))
		return c.SendString(msg) // => ✋ register
	}).Name("api")

	log.Fatal(app.Listen(":3001"))
}
