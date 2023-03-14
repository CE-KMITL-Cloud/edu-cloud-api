// Package handler - handling context
package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/edu-cloud-api/database"
	"github.com/gofiber/fiber/v2"
)

// DeleteInstanceDB - Delete instance in DB
/*
	using Params
	@vmid

	using Query
	@username
*/
func DeleteInstanceDB(c *fiber.Ctx) error {
	vmid := c.Params("vmid")
	username := c.Query("username")

	// Check that user is owner of given VM
	instance, getInstanceErr := database.GetInstance(vmid)
	if getInstanceErr != nil {
		log.Printf("Error: Getting instance ID : %s from DB due to %s", vmid, getInstanceErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting instance ID : %s from DB due to %s", vmid, getInstanceErr)})
	}
	if instance.OwnerID != username {
		log.Printf("Error: deleting VMID : %s due to user is not owner of VM", vmid)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed deleting VMID : %s due to user is not owner of VM", vmid)})
	}

	// Delete instance
	log.Printf("Deleting instance : %s", vmid)
	deleteErr := database.DeleteInstance(vmid)
	if deleteErr != nil {
		log.Println("Error: Could not delete instance in DB due to :", deleteErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed deleting instance : %s due to %s", vmid, deleteErr)})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Deleting instance %s successfully", vmid)})
}
