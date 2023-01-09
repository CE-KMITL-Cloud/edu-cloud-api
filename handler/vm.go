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
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", node, vmid)
	urlStr := u.String()

	// Getting VM's info
	info, err := internal.GetVM(urlStr, username, cookies)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": info})
}

// GetVMList - Getting specific VM's info
// GET /api2/json/nodes/{node}/qemu
/*
	using Query Params
	@username : account's username
	@node : node's name
*/
func GetVMList(c *fiber.Ctx) error {
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
	info, err := internal.GetVMList(urlStr, username, cookies)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": info})
}

// CreateVM - Create VM on specific node
// POST /api2/json/nodes/{node}/qemu
/*
	using Query Params
	@node : node's name

	using Request's Body
	@vmid : VM's ID
	@name : VM's name
	@memory : e.g. 1024 (MB)
	@cores : e.g. 2 (cores)
	@sockets : e.g. 2 (sockets of cpu)
	@onboot : {0, 1}
	@scsi0 : "ceph-vm:32"
	@cdrom : "cephfs:iso/ubuntu-20.04.4-live-server-amd64.iso"
	@net0 : "virtio,bridge=vmbr0,firewall=1"
	@scsihw : "virtio-scsi-single"
*/
func CreateVM(c *fiber.Ctx) error {
	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Getting request's body
	createBody := new(model.CreateBody)
	if err := c.BodyParser(createBody); err != nil {
		return err
	}

	// Getting data from query & Mapping values
	node := c.Query("node")

	data := url.Values{}
	data.Set("vmid", fmt.Sprint(createBody.VMID))
	data.Set("name", createBody.Name)
	data.Set("memory", fmt.Sprint(createBody.Memory))
	data.Set("cores", fmt.Sprint(createBody.Cores))
	data.Set("sockets", fmt.Sprint(createBody.Sockets))
	data.Set("onboot", fmt.Sprint(createBody.Onboot))
	data.Set("scsi0", createBody.SCSI0)
	data.Set("cdrom", createBody.CDROM)
	data.Set("net0", createBody.Net0)
	data.Set("scsihw", createBody.SCSIHW)
	log.Println("create body :", data)

	// Construct URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu", node)
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
	info, err := internal.CreateVM(urlStr, data, cookies)
	if err != nil {
		return err
	}
	log.Println(info)
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Creating Success"})
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
// TODO : need to check new cloned VM status until it ready to use, able to clone only VM Template
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
	log.Println("clone body :", data)

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

// CreateTemplate - Templating specific VM
// POST /api2/json/nodes/{node}/qemu/{vmid}}/template
/*
	using Query Params
	@node : node's name
	@vmid : VM's ID
*/
// TODO : make sure that VM has stopped before start templating
func CreateTemplate(c *fiber.Ctx) error {
	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Getting data from query & Mapping values
	node := c.Query("node")
	vmid := c.Query("vmid")

	// Construct URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/template", node, vmid)
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
	info, err := internal.CreateTemplate(urlStr, cookies)
	if err != nil {
		return err
	}
	log.Println(info)
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Templating Success"})
}
