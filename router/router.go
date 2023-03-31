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

	// DB's User
	user := app.Group("user")
	user.Get("/group/:group", handler.GetUsersDB)
	user.Post("/create", handler.CreateUserDB) // create user, user's limit in DB
	user.Get(":username", handler.GetUserDB)
	user.Delete(":username/delete", handler.DeleteUserDB) // delete user, user's limit in DB
	user.Put(":username/update", handler.UpdateUserDB)

	// user's limit
	user.Get(":username/limit", handler.GetUserLimitDB)
	user.Put(":username/limit/update", handler.UpdateUserLimitDB)

	// Pool
	pool := app.Group("pool")
	pool.Get("/owner/:username", handler.GetPoolsDB)
	pool.Get(":code/owner/:username", handler.GetPoolDB)
	pool.Post("/create", handler.CreatePoolDB)
	pool.Delete(":code/owner/:username", handler.DeletePoolDB)
	pool.Get(":code/owner/:username/members/remain", handler.GetRemainStudents)
	pool.Post(":code/owner/:username/members/add", handler.AddMembersPoolDB)
	pool.Post(":code/owner/:username/instances/add", handler.AddInstancesPoolDB)

	// Proxmox's Access
	access := app.Group("/access")
	access.Post("/ticket", handler.GetTicket)
	access.Post("/user/create", handler.CreateUser) // create user in Proxmox
	access.Put("/user/:username/update", handler.UpdateUser)
	access.Delete("/user/:username/delete", handler.DeleteUser)

	// Node
	node := app.Group("/node")
	// node.Get(":node/vm/list", handler.GetVMListByNode) // ! to be deprecated
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

	// VNC
	vm.Post("/vncproxy", handler.GetVncTicket)

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
	storage.Get("/iso/list", handler.GetISOList)

	// Node
	clusterNode := cluster.Group("node")
	clusterNode.Get("/:name", handler.GetNode)
}
