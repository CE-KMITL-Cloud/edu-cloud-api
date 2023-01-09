// Package router - routing paths
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

	// VM
	vm := app.Group("/vm")
	vm.Get("/info", handler.GetVM)
	vm.Get("/list", handler.GetVMList)
	vm.Post("/create", handler.CreateVM)
	vm.Delete("/destroy", handler.DeleteVM)
	vm.Post("/clone", handler.CloneVM)
	vm.Post("/template", handler.CreateTemplate)
}
