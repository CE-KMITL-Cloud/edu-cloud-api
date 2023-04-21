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
	if group == config.STUDENT || owner != sender {
		log.Println("Error: user's group is not allowed or not owner to get pools")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to get pools due to user's group is not allowed or not owner"})
	}
	pools, getPoolsErr := database.GetPoolsByOwner(owner)
	if getPoolsErr != nil {
		log.Printf("Error: getting pools by given owner : %s due to %s", owner, getPoolsErr)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed to getting pools due to %s", getPoolsErr)})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": pools})
}

// GetPoolsByMemberDB - Get pools that sender is member
/*
	using Query
	@username : sender
*/
func GetPoolsByMemberDB(c *fiber.Ctx) error {
	sender := c.Query("username")
	group, getGroupErr := database.GetUserGroup(sender)
	if getGroupErr != nil {
		log.Println("Error: while getting user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}
	if group != config.STUDENT {
		log.Println("Error: user's group is not allowed to get pools")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to get pools due to user's group is not allowed"})
	}
	pools, getPoolsErr := database.GetAllPoolsByMember(sender)
	if getPoolsErr != nil {
		log.Printf("Error: getting pools by given member : %s due to %s", sender, getPoolsErr)
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
	pool, getPoolErr := database.GetPoolByCode(code, owner)
	if getPoolErr != nil {
		log.Printf("Error: getting pool by given owner : %s, code : %s due to %s", owner, code, getPoolErr)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed to getting pool due to %s", getPoolErr)})
	}
	isOwner := database.IsPoolOwner(code, owner, sender, group)
	isMember := database.IsPoolMember(code, owner, sender)
	if isMember || group == config.ADMIN || isOwner {
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": pool})
	}
	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to getting pool due to user is not member or owner"})
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
	// check duplicate pool
	pools, _ := database.GetAllPools()
	for _, pool := range pools {
		if pool.Code == createBody.Code && pool.Owner == createBody.Owner {
			log.Printf("Error: found pool code : %s, owner : %s exists", createBody.Code, createBody.Owner)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed to create pool due to found pool code : %s, owner : %s exists", createBody.Code, createBody.Owner)})
		}
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
	isOwner := database.IsPoolOwner(code, owner, sender, group)
	if isOwner || group == config.ADMIN {
		deletePoolErr := database.DeletePool(code, owner)
		if deletePoolErr != nil {
			log.Printf("Error: deleting pool by given owner : %s, code : %s due to %s", owner, code, deletePoolErr)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed to deleting pool due to %s", deletePoolErr)})
		}
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Target pool code : %s, owner : %s has been deleted", code, owner)})
	}
	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to deleting pool due to user is not owner"})
}

// GetRemainStudents - Getting remain students who not in given pool
/*
	using Params
	@username : pool owner
	@code : course code

	using Query
	@username : sender
*/
func GetRemainStudents(c *fiber.Ctx) error {
	sender := c.Query("username")
	owner := c.Params("username")
	code := c.Params("code")
	group, getGroupErr := database.GetUserGroup(sender)
	if getGroupErr != nil {
		log.Println("Error: while getting user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}
	if group == config.ADMIN || sender == owner {
		students, getStudentErr := database.GetAllStudentsUsername()
		if getStudentErr != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting student list due to %s", getStudentErr)})
		}
		pool, getPoolErr := database.GetPoolByCode(code, owner)
		if getPoolErr != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting pool from given code, owner due to %s", getPoolErr)})
		}
		members := config.FilterList(students, pool.Member)
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": members})
	}
	log.Println("Error: user's group is not allowed to get pools")
	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to get pools due to user's group is not allowed"})
}

// AddMembersPoolDB - Add members to specific pool
/*
	using Request Body
	@members : adding members

	using Params
	@username : pool owner
	@code : course code

	using Query
	@username : sender
*/
func AddMembersPoolDB(c *fiber.Ctx) error {
	addMembersBody := new(model.AddPoolMemberBody)
	if err := c.BodyParser(addMembersBody); err != nil {
		log.Println("Error: Could not parse body parser to add pool's members body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to add pool's members body"})
	}
	students, getStudentErr := database.GetAllStudentsUsername()
	if getStudentErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting student list due to %s", getStudentErr)})
	}
	for _, student := range addMembersBody.Member {
		if !config.Contains(students, student) {
			log.Printf("Error: username: %s in adding list is not exist", student)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed to add pool's username: %s due to user not found", student)})
		}
	}
	sender := c.Query("username")
	owner := c.Params("username")
	code := c.Params("code")
	group, getGroupErr := database.GetUserGroup(sender)
	if getGroupErr != nil {
		log.Println("Error: while getting user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}
	if group == config.ADMIN || sender == owner {
		pool, getPoolErr := database.GetPoolByCode(code, owner)
		if getPoolErr != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting pool from given code, owner due to %s", getPoolErr)})
		}
		for _, member := range addMembersBody.Member {
			if config.Contains(pool.Member, member) {
				log.Printf("Error: Found duplicate username: %s in given pool", member)
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed to add pool's username: %s due to found duplicate user", member)})
			}
		}
		// update pool member in DB
		addMembersBody.Member = append(addMembersBody.Member, pool.Member...)
		updateErr := database.AddPoolMembers(pool.Code, pool.Owner, addMembersBody.Member)
		if updateErr != nil {
			log.Printf("Error: updating member of pool code : %s, owner : %s due to %s", pool.Code, pool.Owner, updateErr)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Internal server error", "message": fmt.Sprintf("Failed updating member of pool code : %s, owner : %s due to %s", pool.Code, pool.Owner, updateErr)})
		}
		log.Printf("Successfully added members : %v to pool code : %s, owner : %s", addMembersBody.Member, code, owner)
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Added new members : %v in pool code : %s, owner : %s successfully", addMembersBody.Member, code, owner)})
	}
	log.Println("Error: user's group is not allowed to get pools")
	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to get pools due to user's group is not allowed"})
}

// AddInstancesPoolDB - Add instance to specific pool
/*
	using Request Body
	@vmid : adding vmid

	using Params
	@username : pool owner
	@code : course code

	using Query
	@username : sender
*/
func AddInstancesPoolDB(c *fiber.Ctx) error {
	addInstanceBody := new(model.AddPoolInstanceBody)
	if err := c.BodyParser(addInstanceBody); err != nil {
		log.Println("Error: Could not parse body parser to add pool's instance body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to add pool's instance body"})
	}
	sender := c.Query("username")
	owner := c.Params("username")
	code := c.Params("code")
	group, getGroupErr := database.GetUserGroup(sender)
	if getGroupErr != nil {
		log.Println("Error: while getting user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}
	if group == config.ADMIN || sender == owner {
		// Check that user is owner of given VM
		instanceTemplateOwner, _ := database.CheckInstanceTemplateOwner(sender, addInstanceBody.VMID)
		if !instanceTemplateOwner && group != config.ADMIN {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed adding VMID : %s due to VM is not template or user is not owner", addInstanceBody.VMID)})
		}
		pool, getPoolErr := database.GetPoolByCode(code, owner)
		if getPoolErr != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting pool from given code, owner due to %s", getPoolErr)})
		}
		template, _ := database.IsInstanceTemplate(addInstanceBody.VMID)
		if !template {
			log.Printf("Error: Found duplicate VMID: %s in given pool", addInstanceBody.VMID)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed to add pool's instance ID: %s due to instance is not template", addInstanceBody.VMID)})
		}
		if !config.Contains(pool.VMID, addInstanceBody.VMID) && template {
			// update pool member in DB
			pool.VMID = append(pool.VMID, addInstanceBody.VMID)
			updateErr := database.AddPoolInstances(pool.Code, pool.Owner, pool.VMID)
			if updateErr != nil {
				log.Printf("Error: updating instances of pool code : %s, owner : %s due to %s", pool.Code, pool.Owner, updateErr)
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Internal server error", "message": fmt.Sprintf("Failed updating instances of pool code : %s, owner : %s due to %s", pool.Code, pool.Owner, updateErr)})
			}
			log.Printf("Successfully added template ID : %s to pool code : %s, owner : %s", addInstanceBody.VMID, code, owner)
			return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Added new instances : %v in pool code : %s, owner : %s successfully", addInstanceBody.VMID, code, owner)})
		}
		log.Printf("Error: Found duplicate VMID: %s in given pool", addInstanceBody.VMID)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed to add pool's instance ID: %s due to found duplicate ID", addInstanceBody.VMID)})
	}
	log.Println("Error: user's group is not allowed to get pools")
	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to get pools due to user's group is not allowed"})
}
