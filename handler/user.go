// Package handler - handling context
package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/edu-cloud-api/database"
	"github.com/gofiber/fiber/v2"
)

// DeleteUserDB - Delete user in DB
/*
	using Params
	@username
*/
func DeleteUserDB(c *fiber.Ctx) error {
	// Getting params from URL
	username := c.Params("username")

	// Getting user's group
	group, getGroupErr := database.GetUserGroup(username)
	if getGroupErr != nil {
		log.Println("Error: Could not get user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting user's group due to %s", getGroupErr)})
	}

	// Deleting User
	log.Printf("Deleting user : %s", username)
	deleteErr := database.DeleteUser(username, group)
	if deleteErr != nil {
		log.Println("Error: Could not delete user in DB due to :", deleteErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed deleting user : %s due to %s", username, deleteErr)})
	}

	// Deleting instance limit
	log.Printf("Deleting user's instance limit : %s", username)
	deleteLimitErr := database.DeleteInstanceLimit(username)
	if deleteLimitErr != nil {
		log.Println("Error: Could not delete user's instance limit in DB due to :", deleteLimitErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed deleting user's instance limit : %s due to %s", username, deleteLimitErr)})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Deleting user %s successfully", username)})
}

// UpdateUserDB - Update user in DB
/*
	using Params
	@username

	using Request body
	@password
	@groups
	@name
	@status
	@create_time
	@expire_time
*/
// func UpdateUserDB(c *fiber.Ctx) error {
// 	// Getting params from URL
// 	username := c.Params("username")

// 	// get user's group

// 	// Creating User
// 	log.Printf("Deleting user : %s", username)
// 	_, deleteErr := database.EditUser(username)
// 	if deleteErr != nil {
// 		log.Println("Error: Could not delete user :", deleteErr)
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed deleting user : %s due to %s", username, deleteErr)})
// 	}
// 	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Deleting user %s successfully", username)})
// }
