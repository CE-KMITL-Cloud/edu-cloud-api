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

// GetNode - Getting node information from given name
// GET /api2/json/cluster/resources
func GetNode(c *fiber.Ctx) error {
	cookies := config.GetCookies(c)
	name := c.Params("name")
	log.Println("Getting Node from given name")
	nodeInfo, err := cluster.GetNode(name, cookies)
	if err != nil {
		log.Println("Error: from getting node info :", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting node info due to %s", err)})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": nodeInfo})
}

// GetStorageList - Getting RBD storage list
// GET /api2/json/cluster/resources
func GetStorageList(c *fiber.Ctx) error {
	cookies := config.GetCookies(c)
	log.Println("Getting RBD Storage list")
	storageList, err := cluster.GetStorageList(cookies)
	if err != nil {
		log.Println("Error: from getting Storage list :", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting RBD Storage list due to %s", err)})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": storageList})
}

// GetISOList - Getting cephfs storage's iso file list
// GET /api2/json/nodes/{node}/storage/{storage}/content
func GetISOList(c *fiber.Ctx) error {
	cookies := config.GetCookies(c)
	log.Println("Getting ISO file list")
	ISOList, err := cluster.GetISOList(cookies)
	if err != nil {
		log.Println("Error: from getting ISO file list :", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting ISO file list due to %s", err)})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": ISOList})
}
