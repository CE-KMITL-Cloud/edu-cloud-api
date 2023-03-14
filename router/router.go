// Package router - routing paths
package router

import (
	"github.com/edu-cloud-api/handler"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes - setting up router
func SetupRoutes(app *fiber.App) {
	// Health Check
	app.Get("/", handler.Healthy)

	// Realm
	// realm := app.Group("/realm")
	// realm.Post("/sync", handler.RealmSync) // ! Not use

	// DB's User
	user := app.Group("user")
	user.Delete(":username/delete", handler.DeleteUserDB) // delete user, user's limit in DB
	// user.Put(":username/update")    // update user in DB

	// DB's Instance
	instance := app.Group("instance")
	instance.Delete(":vmid/delete", handler.DeleteInstanceDB)

	// Access
	access := app.Group("/access")
	access.Post("/ticket", handler.GetTicket)
	access.Post("/user/create", handler.CreateUser)             // create user in Proxmox
	access.Put("/user/:username/update", handler.UpdateUser)    // update user in Proxmox
	access.Delete("/user/:username/delete", handler.DeleteUser) // delete user in Proxmox

	// Node
	node := app.Group("/node")
	// node.Get(":node/vm/list", handler.GetVMListByNode) // * to be deprecated
	node.Get(":node/vm/:vmid", handler.GetVM)

	// VM
	vm := app.Group("/vm")
	vm.Get("/list", handler.GetVMList)
	vm.Get("/template/list", handler.GetTemplateList)

	vm.Post("/create", handler.CreateVM)
	vm.Delete("/destroy", handler.DeleteVM)
	vm.Post("/clone", handler.CloneVM)
	vm.Post("/template", handler.CreateTemplate)
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
	clusterNode := cluster.Group("node")
	clusterNode.Get("/:name", handler.GetNode)
}
