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

	// Realm
	realm := app.Group("/realm")
	realm.Post("/sync", handler.RealmSync)

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
	vm.Get("/template/list", handler.GetTemplateList)
	vm.Post("/edit", handler.EditVM)

	// VM Power Management
	status := vm.Group("/status")
	status.Post("/start", handler.StartVM)
	status.Post("/stop", handler.StopVM)
	status.Post("/shutdown", handler.ShutdownVM)
	status.Post("/suspend", handler.SuspendVM)
	status.Post("/resume", handler.ResumeVM)
	status.Post("/reset", handler.ResetVM)

	// Cluster
	cluster := app.Group("/cluster")

	// Storage
	storage := cluster.Group("storage")
	storage.Get("/list", handler.GetStorageList)

	// Node
	node := cluster.Group("node")
	node.Get("/:name", handler.GetNode)
}
