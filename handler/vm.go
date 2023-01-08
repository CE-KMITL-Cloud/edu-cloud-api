// Package handler - handling context
package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/internal"
	"github.com/edu-cloud-api/model"
	"github.com/gofiber/fiber/v2"
)

// GetVM - Getting specific VM's info
// GET /api2/json/nodes/{node}/qemu/{vmid}/status/current
/*
	using Query Params
	@username : account's username
	@node : node's name
	@vmid : VM's ID
*/
func GetVM(c *fiber.Ctx) error {
	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Getting request's query params
	username := c.Query("username")
	node := c.Query("node")
	vmid := c.Query("vmid")

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
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s", node, vmid)
	urlStr := u.String()

	// Getting VM's info
	info, err := internal.GetVMs(urlStr, username, cookies)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": info})
}

// GetVMs - Getting specific VM's info
// GET /api2/json/nodes/{node}/qemu
/*
	using Query Params
	@username : account's username
	@node : node's name
*/
func GetVMs(c *fiber.Ctx) error {
	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Getting request's query params
	username := c.Query("username")
	node := c.Query("node")

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
	info, err := internal.GetVMs(urlStr, username, cookies)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": info})
}

// DeleteVM - Deleting specific VM
// DELETE /api2/json/nodes/{node}/qemu/{vmid}
/*
	using Query Params
	@username : account's username
	@node : node's name
	@vmid : VM's ID
*/
// TODO : need to check VM status until it delete completely
func DeleteVM(c *fiber.Ctx) error {
	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Getting request's query params
	// username := c.Query("username")
	node := c.Query("node")
	vmid := c.Query("vmid")

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
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s", node, vmid)
	urlStr := u.String()

	// Getting VM's info
	info, err := internal.DeleteVM(urlStr, cookies)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": info})
}

// CloneVM - Cloning specific VM
// POST /api2/json/nodes/{node}/qemu/{vmid}}/clone
/*
	using Query Params
	@node : node's name
	@vmid : VM's ID

	using Request's Body
	@newid : new VM's ID
	@name : VM's name
*/
// TODO : need to check new cloned VM status until it ready to use
func CloneVM(c *fiber.Ctx) error {
	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Getting request's body
	cloneBody := new(model.CloneBody)
	if err := c.BodyParser(cloneBody); err != nil {
		return err
	}

	// Getting data from query & Mapping values
	node := c.Query("node")
	vmid := c.Query("vmid")

	data := url.Values{}
	data.Set("newid", fmt.Sprint(cloneBody.NewID))
	data.Set("name", cloneBody.Name)
	log.Println(data)

	// Construct URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/clone", node, vmid)
	urlStr := u.String()

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

	// Getting info
	info, err := internal.CloneVM(urlStr, data, cookies)
	if err != nil {
		return err
	}
	log.Println(info)
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Cloning Success"})
}
