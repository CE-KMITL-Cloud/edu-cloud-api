// Package handler - handling context
package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

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
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", node, vmid)
	urlStr := u.String()

	// Getting VM's info
	info, err := internal.GetVM(urlStr, cookies)
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
	// username := c.Query("username")
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
	info, err := internal.GetVMList(urlStr, cookies)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": info})
}

// CreateVM - Create VM on specific node
// POST /api2/json/nodes/{node}/qemu
/*
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
// TODO : Specific resource pool by add pool in request's body
func CreateVM(c *fiber.Ctx) error {
	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Getting request's body
	createBody := new(model.CreateBody)
	if err := c.BodyParser(createBody); err != nil {
		return err
	}
	vmid := fmt.Sprint(createBody.VMID)

	// TODO: Another work around here?
	// Construct payload
	data := url.Values{}
	data.Set("vmid", vmid)
	data.Set("name", createBody.Name)
	data.Set("memory", fmt.Sprint(createBody.Memory))
	data.Set("cores", fmt.Sprint(createBody.Cores))
	data.Set("sockets", fmt.Sprint(createBody.Sockets))
	data.Set("onboot", fmt.Sprint(createBody.Onboot))
	data.Set("scsi0", createBody.SCSI0)
	data.Set("cdrom", createBody.CDROM)
	data.Set("net0", createBody.Net0)
	data.Set("scsihw", createBody.SCSIHW)

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

	// Getting target node from node allocation
	target, nodeErr := internal.AllocateNode(cookies)
	if nodeErr != nil {
		return nodeErr
	}
	log.Println("create body :", data, "target node:", target)

	// Construct URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu", target)
	urlStr := u.String()

	// Getting info
	info, err := internal.CreateVM(urlStr, data, cookies)
	if err != nil {
		return err
	}
	log.Println(info)

	// Waiting until creating process has been complete
	var wg sync.WaitGroup
	wg.Add(1)
	go internal.StatusVM(target, vmid, []string{"created", "starting", "running"}, &wg, cookies)
	wg.Wait()
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": fmt.Sprintf("Creating new VMID: %s in %s successfully", vmid, target)})
}

// DeleteVM - Deleting specific VM
// DELETE /api2/json/nodes/{node}/qemu/{vmid}
/*
	using Query Params
	@username : account's username
	@node : node's name
	@vmid : VM's ID
*/
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

	// Construct Getting info URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", node, vmid)
	urlStr := u.String()

	// First check that target VM has been stopped
	vm, err := internal.GetVM(urlStr, cookies)
	if err != nil {
		return err
	}

	// If target VM's status is not "stopped" then return
	if vm.Info.Status != "stopped" {
		return c.Status(400).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been stopped", vmid, node)})
	}

	// Construct Deleting API URL
	deleteURL, _ := url.ParseRequestURI(hostURL)
	deleteURL.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s", node, vmid)
	deleteURLStr := deleteURL.String()

	// Deleting target VM
	_, deleteErr := internal.DeleteVM(deleteURLStr, cookies)
	if deleteErr != nil {
		return deleteErr
	}

	// Check that target VM has been deleted completely yet
	var wg sync.WaitGroup
	wg.Add(1)
	go internal.DeleteCompletely(node, vmid, &wg, cookies)
	wg.Wait()
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been deleted", vmid, node)})
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
// TODO : Able to clone only VM Template : check template == 1
func CloneVM(c *fiber.Ctx) error {
	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Getting request's body
	cloneBody := new(model.CloneBody)
	if err := c.BodyParser(cloneBody); err != nil {
		return err
	}
	newid := fmt.Sprint(cloneBody.NewID)

	// Getting data from query & Mapping values
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

	// Getting target node from node allocation
	target, nodeErr := internal.AllocateNode(cookies)
	if nodeErr != nil {
		return nodeErr
	}

	// Construct payload
	data := url.Values{}
	data.Set("newid", newid)
	data.Set("name", cloneBody.Name)
	data.Set("target", target)
	log.Println("clone body :", data)

	// Construct URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/clone", node, vmid)
	urlStr := u.String()

	// Getting info
	info, err := internal.CloneVM(urlStr, data, cookies)
	if err != nil {
		return err
	}
	log.Println(info)

	// Waiting until cloning process has been completed
	var wg sync.WaitGroup
	wg.Add(1)
	go internal.StatusVM(node, newid, []string{"created", "starting", "running"}, &wg, cookies)
	wg.Wait()
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": fmt.Sprintf("Cloning new VMID: %s to %s successfully", newid, target)})
}

// CreateTemplate - Templating specific VM
// POST /api2/json/nodes/{node}/qemu/{vmid}}/template
/*
	using Query Params
	@node : node's name
	@vmid : VM's ID
*/
func CreateTemplate(c *fiber.Ctx) error {
	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Getting data from query & Mapping values
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

	// Construct Getting info URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", node, vmid)
	urlStr := u.String()

	// First check that target VM has been stopped
	vm, err := internal.GetVM(urlStr, cookies)
	if err != nil {
		return err
	}

	// If target VM's status is not "stopped" then return
	if vm.Info.Status != "stopped" {
		return c.Status(400).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been stopped", vmid, node)})
	}

	// Construct URL
	templateURL, _ := url.ParseRequestURI(hostURL)
	templateURL.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/template", node, vmid)
	templateURLStr := templateURL.String()

	// Templating VM
	_, templateErr := internal.CreateTemplate(templateURLStr, cookies)
	if templateErr != nil {
		return templateErr
	}

	// Waiting until templating process has been completed
	var wg sync.WaitGroup
	wg.Add(1)
	go internal.TemplateCompletely(node, vmid, []string{"created", "existing"}, &wg, cookies)
	wg.Wait()
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": fmt.Sprintf("Target VMID: %s in %s has been templated", vmid, node)})
}
