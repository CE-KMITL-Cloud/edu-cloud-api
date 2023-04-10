// Package handler - handling context
package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/database"
	"github.com/edu-cloud-api/internal/qemu"
	"github.com/edu-cloud-api/model"
	"github.com/gofiber/fiber/v2"
)

// StartVM - Start specific VM
// POST /api2/json/nodes/{node}/qemu/{vmid}/status/start
/*
	using Request's Body
	@node : node's name
	@vmid : VM's ID

	using Query
	@username : account's username
*/
func StartVM(c *fiber.Ctx) error {
	// Getting request's body
	startBody := new(model.StartBody)
	if err := c.BodyParser(startBody); err != nil {
		log.Println("Error: Could not parse body parser to start VM's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to start VM's body"})
	}
	vmid := fmt.Sprint(startBody.VMID)
	username := c.Query("username")
	owner, checkOwnerErr := database.CheckInstanceOwner(username, vmid)
	if checkOwnerErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to %s", vmid, checkOwnerErr)})
	}
	if !owner {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to user is not owner of VM", vmid)})
	}
	cookies := config.GetCookies(c)

	// Getting VM's info
	vmStatusURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", startBody.Node, vmid))
	vm, err := qemu.GetVM(vmStatusURL, cookies)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting detail from VMID: %s in %s due to %s", vmid, startBody.Node, err)})
	}

	// If target VM's status is not "stopped" then return
	if vm.Info.Status != "stopped" {
		log.Printf("Error: Could not start VMID : %s in %s due to VM hasn't been stopped", vmid, startBody.Node)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been stopped", vmid, startBody.Node)})
	}

	// Starting VM
	log.Printf("Starting VMID : %s in %s", vmid, startBody.Node)
	vmStartURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/start", startBody.Node, vmid))
	_, startErr := qemu.PowerManagement(vmStartURL, nil, cookies)
	if startErr != nil {
		log.Printf("Error: Could not start VMID : %s in %s : %s", vmid, startBody.Node, startErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed starting VMID : %s in %s due to %s", vmid, startBody.Node, startErr)})
	}

	// Waiting until starting process has been completed
	started := qemu.CheckStatus(startBody.Node, vmid, []string{"running"}, false, (5 * time.Minute), time.Second)
	if started {
		log.Printf("Finished starting VMID : %s in %s", vmid, startBody.Node)
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Target VMID: %s in %s has been started", vmid, startBody.Node)})
	}
	log.Printf("Error: Could not start VMID : %s in %s", vmid, startBody.Node)
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been started correctly", vmid, startBody.Node)})
}

// StopVM - Stop specific VM, pulling the power plug of a running computer and may damage the VM data
// POST /api2/json/nodes/{node}/qemu/{vmid}/status/stop
/*
	using Request's Body
	@node : node's name
	@vmid : VM's ID

	using Query
	@username : account's username
*/
func StopVM(c *fiber.Ctx) error {
	// Getting request's body
	stopBody := new(model.StopBody)
	if err := c.BodyParser(stopBody); err != nil {
		log.Println("Error: Could not parse body parser to stop VM's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to stop VM's body"})
	}
	vmid := fmt.Sprint(stopBody.VMID)
	username := c.Query("username")
	owner, checkOwnerErr := database.CheckInstanceOwner(username, vmid)
	if checkOwnerErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to %s", vmid, checkOwnerErr)})
	}
	if !owner {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to user is not owner of VM", vmid)})
	}
	cookies := config.GetCookies(c)

	// Getting VM's info
	vmStatusURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", stopBody.Node, vmid))
	vm, err := qemu.GetVM(vmStatusURL, cookies)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting detail from VMID: %s in %s due to %s", vmid, stopBody.Node, err)})
	}

	// If target VM's status is not "running" then return
	if vm.Info.Status != "running" {
		log.Printf("Error: Could not stop VMID : %s in %s due to VM hasn't been running", vmid, stopBody.Node)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been running", vmid, stopBody.Node)})
	}

	// Stopping VM
	log.Printf("Stopping VMID : %s in %s", vmid, stopBody.Node)
	vmStopURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/stop", stopBody.Node, vmid))
	_, stopErr := qemu.PowerManagement(vmStopURL, nil, cookies)
	if stopErr != nil {
		log.Printf("Error: Could not stop VMID : %s in %s : %s", vmid, stopBody.Node, stopErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed stopping VMID : %s in %s due to %s", vmid, stopBody.Node, stopErr)})
	}

	// Waiting until stopping process has been completed
	stopped := qemu.CheckStatus(stopBody.Node, vmid, []string{"stopped"}, false, (5 * time.Minute), time.Second)
	if stopped {
		log.Printf("Finished stopping VMID : %s in %s", vmid, stopBody.Node)
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Target VMID: %s in %s has been stopped", vmid, stopBody.Node)})
	}
	log.Printf("Error: Could not stop VMID : %s in %s", vmid, stopBody.Node)
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been stopped correctly", vmid, stopBody.Node)})
}

// ShutdownVM - This is similar to pressing the power button on a physical machine.
// This will send an ACPI event for the guest OS, which should then proceed to a clean shutdown.
// POST /api2/json/nodes/{node}/qemu/{vmid}/status/shutdown
/*
	using Request's Body
	@node : node's name
	@vmid : VM's ID

	using Query
	@username : account's username
*/
func ShutdownVM(c *fiber.Ctx) error {
	// Getting request's body
	shutdownBody := new(model.ShutdownBody)
	if err := c.BodyParser(shutdownBody); err != nil {
		log.Println("Error: Could not parse body parser to shut down VM's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to shut down VM's body"})
	}
	vmid := fmt.Sprint(shutdownBody.VMID)

	// Construct payload
	data := url.Values{}
	data.Set("forceStop", "1") // ! Fixed to set "1" for waiting until VM stopped
	username := c.Query("username")
	owner, checkOwnerErr := database.CheckInstanceOwner(username, vmid)
	if checkOwnerErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to %s", vmid, checkOwnerErr)})
	}
	if !owner {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to user is not owner of VM", vmid)})
	}
	cookies := config.GetCookies(c)

	// Getting VM's info
	vmStatusURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", shutdownBody.Node, vmid))
	vm, err := qemu.GetVM(vmStatusURL, cookies)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting detail from VMID: %s in %s due to %s", vmid, shutdownBody.Node, err)})
	}

	// If target VM's status is not "running" then return
	if vm.Info.Status != "running" {
		log.Printf("Error: Could not stop VMID : %s in %s due to VM hasn't been running", vmid, shutdownBody.Node)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been running", vmid, shutdownBody.Node)})
	}

	// Shutting down VM
	log.Printf("Shutting down VMID : %s in %s", vmid, shutdownBody.Node)
	vmShutdownURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/shutdown", shutdownBody.Node, vmid))
	_, shutdownErr := qemu.PowerManagement(vmShutdownURL, data, cookies)
	if shutdownErr != nil {
		log.Printf("Error: Could not shut down VMID : %s in %s : %s", vmid, shutdownBody.Node, shutdownErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed shutting down VMID : %s in %s due to %s", vmid, shutdownBody.Node, shutdownErr)})
	}

	// Waiting until shutting down process has been completed
	shutdown := qemu.CheckStatus(shutdownBody.Node, vmid, []string{"stopped"}, false, (5 * time.Minute), (3 * time.Second))
	if shutdown {
		log.Printf("Finished shutting down VMID : %s in %s", vmid, shutdownBody.Node)
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Target VMID: %s in %s has been shut down", vmid, shutdownBody.Node)})
	}
	log.Printf("Error: Could not shut down VMID : %s in %s", vmid, shutdownBody.Node)
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been shut down correctly", vmid, shutdownBody.Node)})
}

// SuspendVM - Suspend specific VM
// POST /api2/json/nodes/{node}/qemu/{vmid}/status/suspend
/*
	using Request's Body
	@node : node's name
	@vmid : VM's ID

	using Query
	@username : account's username
*/
func SuspendVM(c *fiber.Ctx) error {
	// Getting request's body
	suspendBody := new(model.SuspendBody)
	if err := c.BodyParser(suspendBody); err != nil {
		log.Println("Error: Could not parse body parser to suspend VM's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to suspend VM's body"})
	}
	vmid := fmt.Sprint(suspendBody.VMID)
	username := c.Query("username")
	owner, checkOwnerErr := database.CheckInstanceOwner(username, vmid)
	if checkOwnerErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to %s", vmid, checkOwnerErr)})
	}
	if !owner {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to user is not owner of VM", vmid)})
	}
	cookies := config.GetCookies(c)

	// Getting VM's info
	vmStatusURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", suspendBody.Node, vmid))
	vm, err := qemu.GetVM(vmStatusURL, cookies)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting detail from VMID: %s in %s due to %s", vmid, suspendBody.Node, err)})
	}

	// If target VM's QMP Status is not "running" then return
	if vm.Info.QmpStatus != "running" {
		log.Printf("Error: Could not suspend VMID : %s in %s due to QMP Status of VM hasn't been running", vmid, suspendBody.Node)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Target VMID: %s in %s QMP Status hasn't been running", vmid, suspendBody.Node)})
	}

	// Suspending VM
	log.Printf("Suspending VMID : %s in %s", vmid, suspendBody.Node)
	vmSuspendURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/suspend", suspendBody.Node, vmid))
	_, suspendErr := qemu.PowerManagement(vmSuspendURL, nil, cookies)
	if suspendErr != nil {
		log.Printf("Error: Could not suspend VMID : %s in %s : %s", vmid, suspendBody.Node, suspendErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed suspending VMID : %s in %s due to %s", vmid, suspendBody.Node, suspendErr)})
	}

	// Waiting until suspending process has been completed
	suspended := qemu.CheckQmpStatus(suspendBody.Node, vmid, []string{"paused"}, false, (5 * time.Minute), time.Second)
	if suspended {
		log.Printf("Finished suspending VMID : %s in %s", vmid, suspendBody.Node)
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Target VMID: %s in %s has been suspended", vmid, suspendBody.Node)})
	}
	log.Printf("Error: Could not suspend VMID : %s in %s", vmid, suspendBody.Node)
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been suspended correctly", vmid, suspendBody.Node)})
}

// ResumeVM - Resume specific VM
// POST /api2/json/nodes/{node}/qemu/{vmid}/status/resume
/*
	using Request's Body
	@node : node's name
	@vmid : VM's ID

	using Query
	@username : account's username
*/
func ResumeVM(c *fiber.Ctx) error {
	// Getting request's body
	resumeBody := new(model.ResumeBody)
	if err := c.BodyParser(resumeBody); err != nil {
		log.Println("Error: Could not parse body parser to resume VM's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to resume VM's body"})
	}
	vmid := fmt.Sprint(resumeBody.VMID)
	username := c.Query("username")
	owner, checkOwnerErr := database.CheckInstanceOwner(username, vmid)
	if checkOwnerErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to %s", vmid, checkOwnerErr)})
	}
	if !owner {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to user is not owner of VM", vmid)})
	}
	cookies := config.GetCookies(c)

	// Getting VM's info
	vmStatusURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", resumeBody.Node, vmid))
	vm, err := qemu.GetVM(vmStatusURL, cookies)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting detail from VMID: %s in %s due to %s", vmid, resumeBody.Node, err)})
	}

	// If target VM's QMP Status is not "paused" then return
	if vm.Info.QmpStatus != "paused" {
		log.Printf("Error: Could not resume VMID : %s in %s due to QMP Status of VM hasn't been paused", vmid, resumeBody.Node)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Target VMID: %s in %s QMP Status hasn't been paused", vmid, resumeBody.Node)})
	}

	// Resuming VM
	log.Printf("Resuming VMID : %s in %s", vmid, resumeBody.Node)
	vmResumeURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/resume", resumeBody.Node, vmid))
	_, resumeErr := qemu.PowerManagement(vmResumeURL, nil, cookies)
	if resumeErr != nil {
		log.Printf("Error: Could not resume VMID : %s in %s : %s", vmid, resumeBody.Node, resumeErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed resuming VMID : %s in %s due to %s", vmid, resumeBody.Node, resumeErr)})
	}

	// Waiting until resuming process has been completed
	resumed := qemu.CheckQmpStatus(resumeBody.Node, vmid, []string{"running"}, false, (5 * time.Minute), time.Second)
	if resumed {
		log.Printf("Finished resuming VMID : %s in %s", vmid, resumeBody.Node)
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Target VMID: %s in %s has been resumed", vmid, resumeBody.Node)})
	}
	log.Printf("Error: Could not resume VMID : %s in %s", vmid, resumeBody.Node)
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been resumed correctly", vmid, resumeBody.Node)})
}

// ResetVM - Reset specific VM
// POST /api2/json/nodes/{node}/qemu/{vmid}/status/reset
/*
	using Request's Body
	@node : node's name
	@vmid : VM's ID

	using Query
	@username : account's username
*/
func ResetVM(c *fiber.Ctx) error {
	// Getting request's body
	resetBody := new(model.ResetBody)
	if err := c.BodyParser(resetBody); err != nil {
		log.Println("Error: Could not parse body parser to reset VM's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to reset VM's body"})
	}
	vmid := fmt.Sprint(resetBody.VMID)
	username := c.Query("username")
	owner, checkOwnerErr := database.CheckInstanceOwner(username, vmid)
	if checkOwnerErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to %s", vmid, checkOwnerErr)})
	}
	if !owner {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to user is not owner of VM", vmid)})
	}
	cookies := config.GetCookies(c)

	// Getting VM's info
	vmStatusURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", resetBody.Node, vmid))
	vm, err := qemu.GetVM(vmStatusURL, cookies)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting detail from VMID: %s in %s due to %s", vmid, resetBody.Node, err)})
	}

	// If target VM's status is not "running" then return
	if vm.Info.Status != "running" {
		log.Printf("Error: Could not reset VMID : %s in %s due to VM hasn't been running", vmid, resetBody.Node)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been running", vmid, resetBody.Node)})
	}

	// Resetting VM
	log.Printf("Resetting VMID : %s in %s", vmid, resetBody.Node)
	vmResetURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/reset", resetBody.Node, vmid))
	_, resetErr := qemu.PowerManagement(vmResetURL, nil, cookies)
	if resetErr != nil {
		log.Printf("Error: Could not reset VMID : %s in %s : %s", vmid, resetBody.Node, resetErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed resetting VMID : %s in %s due to %s", vmid, resetBody.Node, resetErr)})
	}

	// Waiting until resetting process has been completed
	reset := qemu.CheckStatus(resetBody.Node, vmid, []string{"running"}, false, (5 * time.Minute), time.Second)
	if reset {
		log.Printf("Finished resetting VMID : %s in %s", vmid, resetBody.Node)
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Target VMID: %s in %s has been reset", vmid, resetBody.Node)})
	}
	log.Printf("Error: Could not reset VMID : %s in %s", vmid, resetBody.Node)
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been reset correctly", vmid, resetBody.Node)})
}
