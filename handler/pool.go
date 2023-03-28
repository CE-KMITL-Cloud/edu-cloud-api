// Package handler - handling context
package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/database"
	"github.com/edu-cloud-api/model"
	"github.com/gofiber/fiber/v2"
)

// GetPoolsDB - Get pools from given owner
/*
	using Params
	@username

	using Query
	@username : sender
*/
func GetPoolsDB(c *fiber.Ctx) error {
	owner := c.Params("username")
	sender := c.Query("username")
	group, getGroupErr := database.GetUserGroup(sender)
	if getGroupErr != nil {
		log.Println("Error: while getting user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}
	if group != config.ADMIN {
		log.Println("Error: user's group is not allowed to get pools")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to get pools due to user's group is not allowed"})
	}
	pools, getPoolsErr := database.GetPoolsByOwner(owner)
	if getPoolsErr != nil {
		log.Printf("Error: getting pools by given owner : %s due to %s", owner, getPoolsErr)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed to getting pools due to %s", getPoolsErr)})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": pools})
}

// GetPoolDB - Get pool from given course code, owner
/*
	using Params
	@username
	@code

	using Query
	@username : sender
*/
func GetPoolDB(c *fiber.Ctx) error {
	owner := c.Params("username")
	code := c.Params("code")
	sender := c.Query("username")
	group, getGroupErr := database.GetUserGroup(sender)
	if getGroupErr != nil {
		log.Println("Error: while getting user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}
	isMember, isMemberErr := database.IsPoolMember(code, owner, sender)
	if isMemberErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to checking pool member due to %s", isMemberErr)})
	}
	if isMember || group == config.ADMIN {
		pool, getPoolErr := database.GetPoolByCode(code, owner)
		if getPoolErr != nil {
			log.Printf("Error: getting pool by given owner : %s, code : %s due to %s", owner, code, getPoolErr)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed to getting pool due to %s", getPoolErr)})
		}
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": pool})
	}
	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to getting pool due to user is not member"})
}

// CreatePoolDB - Create pool
/*
	using Request body
	@owner
	@code
	@name

	using Query
	@username : sender
*/
func CreatePoolDB(c *fiber.Ctx) error {
	createBody := new(model.CreatePoolBody)
	if err := c.BodyParser(createBody); err != nil {
		log.Println("Error: Could not parse body parser to create pool's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to create pool's body"})
	}
	// Check owner's role
	ownerGroup, getOwnerGroupErr := database.GetUserGroup(createBody.Owner)
	if getOwnerGroupErr != nil {
		log.Println("Error: while getting owner's group due to :", getOwnerGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting owner's group due to %s", getOwnerGroupErr)})
	}
	// Check sender's role
	sender := c.Query("username")
	senderGroup, getSenderGroupErr := database.GetUserGroup(sender)
	if getSenderGroupErr != nil {
		log.Println("Error: while getting sender's group due to :", getSenderGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting sender's group due to %s", getSenderGroupErr)})
	}
	// only faculty and admin
	if senderGroup != config.STUDENT && ownerGroup != config.STUDENT {
		// sender is faculty role but create for the other
		if senderGroup == config.FACULTY && sender != createBody.Owner {
			log.Println("Error: faculty role is able to create pool only for their own")
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to create pool due to user's group is not allowed to create for other"})
		}
		// Create pool in DB
		pool, createPoolErr := database.CreatePool(createBody)
		if createPoolErr != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed to creating pool due to %s", createPoolErr)})
		}
		log.Printf("Finished creating pool : %s, owner : %s", createBody.Name, createBody.Owner)
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": pool})
	}
	log.Println("Error: user's group is not allowed to create pool")
	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to create pool due to user's group is not allowed"})
}

// DeletePoolDB - Delete pool from given course code, owner
/*
	using Params
	@username
	@code

	using Query
	@username : sender
*/
func DeletePoolDB(c *fiber.Ctx) error {
	owner := c.Params("username")
	code := c.Params("code")
	sender := c.Query("username")
	group, getGroupErr := database.GetUserGroup(sender)
	if getGroupErr != nil {
		log.Println("Error: while getting user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}
	if group == config.STUDENT {
		log.Println("Error: user's group is not allowed to get pools")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to get pools due to user's group is not allowed"})
	}
	isOwner, isOwnerErr := database.IsPoolOwner(code, owner, sender)
	if isOwnerErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to checking pool owner due to %s", isOwnerErr)})
	}
	if isOwner || group == config.ADMIN {
		deletePoolErr := database.DeletePool(code, owner)
		if deletePoolErr != nil {
			log.Printf("Error: deleting pool by given owner : %s, code : %s due to %s", owner, code, deletePoolErr)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed to deleting pool due to %s", deletePoolErr)})
		}
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Target pool code : %s, owner : %s hasn't been deleted", code, owner)})
	}
	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to deleting pool due to user is not owner"})
}
