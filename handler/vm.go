// Package handler - handling context
package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/internal"
	"github.com/edu-cloud-api/model"
	"github.com/gofiber/fiber/v2"
)

// TODO : Using middleware ?

// GetVM - Getting specific VM's info
// ! need to re-think about how we can get specific vm's info
// GET /api2/json/nodes/{node}/qemu/{vmid}
/*
	using Query Params
	@username : account's username
	@node : node's name
	@vmid : VM's ID
*/
// func GetVM(c *fiber.Ctx) error {
// 	// Get host's URL
// 	hostURL := config.GetFromENV("PROXMOX_HOST")

// 	// Getting request's query params
// 	username := c.Query("username")
// 	node := c.Query("nodes")
// 	vmid := c.Query("vmid")

// 	// Set Header
// 	user := model.User{}
// 	database.DB.Db.Where("username = ?", username).Find(&user)

// 	// Construct URL
// 	u, _ := url.ParseRequestURI(hostURL)
// 	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s", node, vmid)
// 	urlStr := u.String()

// 	// Getting VM's info
// 	// info, err := internal.GetVM(urlStr)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	internal.GetVM(urlStr, user)

// 	return c.JSON("Get VM")
// }

// GetVMs - Getting specific VM's info
// GET /api2/json/nodes/{node}/qemu
/*
	using Query Params
	@username : account's username
	@node : node's name
	@vmid : VM's ID
*/
func GetVMs(c *fiber.Ctx) error {
	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Getting request's query params
	username := c.Query("username")
	node := c.Query("nodes")

	// Getting Cookie, CSRF Token
	cookies := model.Cookies{
		Cookie: http.Cookie{
			Name:  "PVEAuthCookie",
			Value: c.Cookies("PVEAuthCookie"),
		},
		CSRFPreventionToken: fiber.Cookie{
			Name:  "CSRFPreventionToken",
			Value: c.Cookies("CSRFPreventionToken"),
		},
	}

	// Construct URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu", node)
	urlStr := u.String()

	// Getting VM's info
	// info, err := internal.GetVM(urlStr)
	// if err != nil {
	// 	return err
	// }
	internal.GetVMs(urlStr, username, cookies)

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Got VMs"})
}
