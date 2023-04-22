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

// GetUserDB - Get users from given username
/*
	using Params
	@username

	using Query
	@username : sender
*/
func GetUserDB(c *fiber.Ctx) error {
	username := c.Params("username")
	sender := c.Query("username")

	// Checking sender's role
	userGroup, getGroupErr := database.GetUserGroup(sender)
	if getGroupErr != nil {
		log.Println("Error: while getting user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}
	if userGroup != config.ADMIN && sender != username {
		log.Println("Error: user's group is not allowed to create user")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to create user due to user's group is not allowed"})
	}

	// Getting user's group
	group, getGroupErr := database.GetUserGroup(username)
	if getGroupErr != nil {
		log.Println("Error: Could not get user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting user's group due to %s", getGroupErr)})
	}

	user, getUserErr := database.GetUser(username, group)
	if getUserErr != nil {
		log.Printf("Error: Could not get user %s from group %s due to : %s", username, group, getUserErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting user %s due to %s", username, getUserErr)})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": user})
}

// GetUsersDB - Get users from given group
/*
	using Params
	@group

	using Query
	@username
*/
func GetUsersDB(c *fiber.Ctx) error {
	group := c.Params("group")
	username := c.Query("username")

	// Checking user's role
	userGroup, getGroupErr := database.GetUserGroup(username)
	if getGroupErr != nil {
		log.Println("Error: while getting user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}
	if userGroup != config.ADMIN {
		log.Println("Error: user's group is not allowed to get users from given group")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to get users due to user's group is not allowed"})
	}

	users, getUsersErr := database.GetAllUsersByGroup(group)
	if getUsersErr != nil {
		log.Printf("Error: Could not get users from given group %s due to : %s", group, getUsersErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting users from given group due to %s", getUsersErr)})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": users})
}

// GetUsersDB - Get users from given group
/*
	using Params
	@group

	using Query
	@username
*/
func GetStudentsDB(c *fiber.Ctx) error {
	username := c.Query("username")

	// Checking user's role
	userGroup, getGroupErr := database.GetUserGroup(username)
	if getGroupErr != nil {
		log.Println("Error: while getting all students due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}
	if userGroup == config.STUDENT {
		log.Println("Error: user's group is not allowed to get all students")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to get all students due to user's group is not allowed"})
	}

	users, getUsersErr := database.GetAllUsersByGroup("student")
	if getUsersErr != nil {
		log.Printf("Error: Could not get all students due to : %s", getUsersErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting all students due to %s", getUsersErr)})
	}
	var students []string
	for _, student := range users {
		students = append(students, student.Username)
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": students})
}

// CreateUserDB - Create user in DB
/*
	using Query
	@username : sender

	using Request body
	@username
	@password
	@name
	@status
*/
func CreateUserDB(c *fiber.Ctx) error {
	// Getting params from URL
	username := c.Query("username")

	// Getting request's body
	body := new(model.CreateUserDB)
	if err := c.BodyParser(body); err != nil {
		log.Println("Error: Could not parse body parser to create user's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to create user's body"})
	}

	// Checking user's role
	userGroup, getGroupErr := database.GetUserGroup(username)
	if getGroupErr != nil {
		log.Println("Error: while getting user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}
	if userGroup != config.ADMIN {
		log.Println("Error: user's group is not allowed to create user")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to create user due to user's group is not allowed"})
	}

	// Checking duplicate username
	usernames, getUsersErr := database.GetUsers()
	if getUsersErr != nil {
		log.Printf("Error: Could not get user's username list due to : %s", getUsersErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting user's username list due to %s", getUsersErr)})
	}
	if config.Contains(usernames, body.Username) {
		log.Printf("Error: Could not create user %s username list due to duplicated username", body.Username)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed creating user %s due to duplicated username", body.Username)})
	}

	// Creating User, User's limit
	log.Printf("Creating user : %s", body.Username)
	_, createErr := database.CreateUserDB(body)
	if createErr != nil {
		log.Printf("Error: Could not create user %s in DB due to : %s", body.Username, createErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed creating user %s due to %s", body.Username, createErr)})
	}
	if createLimitErr := database.CreateInstanceLimit(body.Username, body.Group); createLimitErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed creating user %s's limit due to %s", body.Username, createLimitErr)})
	}
	log.Printf("Finished creating user : %s", body.Username)
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Creating user %s successfully", body.Username)})
}

// DeleteUserDB - Delete user in DB
/*
	using Params
	@username

	using Query
	@username : sender
*/
func DeleteUserDB(c *fiber.Ctx) error {
	// Getting params from URL
	username := c.Params("username")
	sender := c.Query("username")

	// Checking sender's role
	userGroup, getGroupErr := database.GetUserGroup(sender)
	if getGroupErr != nil {
		log.Println("Error: while getting user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}
	if userGroup != config.ADMIN {
		log.Println("Error: user's group is not allowed to create user")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to create user due to user's group is not allowed"})
	}

	// Getting user's group
	group, getGroupErr := database.GetUserGroup(username)
	if getGroupErr != nil {
		log.Println("Error: Could not get user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting user's group due to %s", getGroupErr)})
	}

	// Deleting User
	log.Printf("Deleting user : %s", username)
	deleteErr := database.DeleteUserDB(username, group)
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
	log.Printf("Finished deleting user's instance limit : %s", username)
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Deleting user %s successfully", username)})
}

// UpdateUserDB - Update user in DB
/*
	using Params
	@username

	using Query
	@username : sender

	using Request body
	@password
	@name
	@status
	@expire_time
*/
// todo : check why could not change status from true -> false
func UpdateUserDB(c *fiber.Ctx) error {
	// Getting params from URL
	username := c.Params("username")
	sender := c.Query("username")

	// Checking sender's role
	userGroup, getGroupErr := database.GetUserGroup(sender)
	if getGroupErr != nil {
		log.Println("Error: while getting user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}
	if userGroup != config.ADMIN {
		log.Println("Error: user's group is not allowed to update user")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to update user due to user's group is not allowed"})
	}

	// Getting request's body
	body := new(model.EditUserDB)
	if err := c.BodyParser(body); err != nil {
		log.Println("Error: Could not parse body parser to edit user's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to edit user's body"})
	}

	// get user's group
	group, getGroupErr := database.GetUserGroup(username)
	if getGroupErr != nil {
		log.Println("Error: Could not get user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting user's group due to %s", getGroupErr)})
	}

	// Editing User
	log.Printf("Editing user : %s", username)
	editErr := database.EditUser(username, group, body)
	if editErr != nil {
		log.Printf("Error: Could not edit user %s in DB due to : %s", username, editErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed editing user : %s due to %s", username, editErr)})
	}
	log.Printf("Finished editing user : %s", username)
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Editing user %s successfully", username)})
}

// GetUserLimitDB - Get user's limit from given username
/*
	using Params
	@username

	using Query
	@username : sender
*/
func GetUserLimitDB(c *fiber.Ctx) error {
	username := c.Params("username")
	sender := c.Query("username")

	// Checking sender's role
	userGroup, getGroupErr := database.GetUserGroup(sender)
	if getGroupErr != nil {
		log.Println("Error: while getting user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}
	if userGroup != config.ADMIN && sender != username {
		log.Println("Error: user's group is not allowed to create user")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to create user due to user's group is not allowed"})
	}

	// Getting user's group
	group, getGroupErr := database.GetUserGroup(username)
	if getGroupErr != nil {
		log.Println("Error: Could not get user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting user's group due to %s", getGroupErr)})
	}

	limit, getUserLimitErr := database.GetInstanceLimit(username)
	if getUserLimitErr != nil {
		log.Printf("Error: Could not get user %s from group %s due to : %s", username, group, getUserLimitErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting user %s due to %s", username, getUserLimitErr)})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": limit})
}

// UpdateUserLimitDB - Update user's limit in DB
/*
	using Params
	@username

	using Query
	@username : sender

	using Request body
	@max_cpu
	@max_ram
	@max_disk
	@max_instance
*/
func UpdateUserLimitDB(c *fiber.Ctx) error {
	// Getting params from URL
	username := c.Params("username")
	sender := c.Query("username")

	// Checking sender's role
	userGroup, getGroupErr := database.GetUserGroup(sender)
	if getGroupErr != nil {
		log.Println("Error: while getting user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}
	if userGroup != config.ADMIN {
		log.Println("Error: user's group is not allowed to edit user's limit")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to edit user's limit due to user's group is not allowed"})
	}

	// Getting request's body
	body := new(model.EditInstanceLimit)
	if err := c.BodyParser(body); err != nil {
		log.Println("Error: Could not parse body parser to edit user's limit body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to edit user's limit body"})
	}

	// Editing User
	log.Printf("Editing user's limit : %s", username)
	editErr := database.EditInstanceLimit(username, body)
	if editErr != nil {
		log.Printf("Error: Could not edit user %s's limit in DB due to : %s", username, editErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": editErr})
	}
	log.Printf("Finished editing user's limit : %s", username)
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Editing user %s's limit successfully", username)})
}
