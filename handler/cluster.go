// Package handler - handling context
package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/internal/cluster"
	"github.com/gofiber/fiber/v2"
)

// GetStorageList - Getting RBD storage list
// GET /api2/json/cluster/resources
func GetStorageList(c *fiber.Ctx) error {
	// Getting RBD Storage list
	cookies := config.GetCookies(c)
	log.Println("Getting RBD Storage list")
	storageList, err := cluster.GetStorageList(cookies)
	if err != nil {
		log.Println("Error: from getting VM's list :", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting RBD Storage list due to %s", err)})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": storageList})
}
