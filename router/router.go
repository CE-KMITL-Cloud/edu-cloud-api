// Package router - for routing
package router

import (
	"github.com/edu-cloud-api/handler"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes - setting up router
func SetupRoutes(app *fiber.App) {
	// Health-Check
	app.Get("/", handler.Healthy)

	// Access
	access := app.Group("/access")
	access.Post("/ticket", handler.GetTicket)
}
