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

// StartVM - Start specific VM
// POST /api2/json/nodes/{node}/qemu/{vmid}/status/start
/*
	using Query Params
	@node : node's name
	@vmid : VM's ID
*/
func StartVM(c *fiber.Ctx) error {
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
		log.Printf("Error: Could not start VMID : %s in %s due to VM hasn't been stopped", vmid, node)
		return c.Status(400).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been stopped", vmid, node)})
	}

	// Construct URL
	startURL, _ := url.ParseRequestURI(hostURL)
	startURL.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/start", node, vmid)
	startURLStr := startURL.String()

	// Templating VM
	_, startErr := internal.PowerManagement(startURLStr, cookies)
	if startErr != nil {
		log.Printf("Error: Could not start VMID : %s in %s : %s", vmid, node, startErr)
		return startErr
	}

	// Waiting until templating process has been completed
	started := internal.StatusVM(node, vmid, []string{"running"}, cookies)
	if started {
		log.Printf("Finished starting VMID : %s in %s", vmid, node)
		return c.Status(200).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Target VMID: %s in %s has been started", vmid, node)})
	}
	log.Printf("Error: Could not start VMID : %s in %s", vmid, node)
	return c.Status(500).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been started correctly", vmid, node)})
}

// StopVM - Stop specific VM
// POST /api2/json/nodes/{node}/qemu/{vmid}/status/stop
/*
	using Query Params
	@node : node's name
	@vmid : VM's ID
*/
func StopVM(c *fiber.Ctx) error {
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
	if vm.Info.Status != "running" {
		log.Printf("Error: Could not stop VMID : %s in %s due to VM hasn't been running", vmid, node)
		return c.Status(400).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been running", vmid, node)})
	}

	// Construct URL
	startURL, _ := url.ParseRequestURI(hostURL)
	startURL.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/stop", node, vmid)
	startURLStr := startURL.String()

	// Templating VM
	_, startErr := internal.PowerManagement(startURLStr, cookies)
	if startErr != nil {
		log.Printf("Error: Could not stop VMID : %s in %s : %s", vmid, node, startErr)
		return startErr
	}

	// Waiting until templating process has been completed
	started := internal.StatusVM(node, vmid, []string{"stopped"}, cookies)
	if started {
		log.Printf("Finished stopping VMID : %s in %s", vmid, node)
		return c.Status(200).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Target VMID: %s in %s has been stopped", vmid, node)})
	}
	log.Printf("Error: Could not stop VMID : %s in %s", vmid, node)
	return c.Status(500).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been stopped correctly", vmid, node)})
}
