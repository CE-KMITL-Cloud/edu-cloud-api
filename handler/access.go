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
	"github.com/edu-cloud-api/internal/access"
	"github.com/edu-cloud-api/model"
	"github.com/gofiber/fiber/v2"
)

// GetTicket - handler GetTicket function
// GET /api2/json/access/ticket
/*
	using Request's Body
	@username : account's username
	@password : account's password
*/
func GetTicket(c *fiber.Ctx) error {
	// Getting request's body
	body := new(model.Login)
	if err := c.BodyParser(body); err != nil {
		log.Println("Error: Could not parse body parser to getting ticket's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to getting ticket's body"})
	}

	// Mapping values
	data := url.Values{}
	data.Set("username", body.Username)
	data.Set("password", body.Password)
	data.Set("realm", "pve")

	// Getting Ticket
	log.Printf("Getting ticket from user : %s", body.Username)
	ticket, ticketErr := access.GetTicket(data)
	if ticketErr != nil {
		log.Println("Error: Could not get ticket :", ticketErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting ticket from user : %s due to %s", body.Username, ticketErr)})
	}

	// Set Cookie
	c.Cookie(&fiber.Cookie{
		Name:    config.AUTH_COOKIE,
		Value:   ticket.Token.Cookie,
		Expires: time.Now().Add(time.Hour * 12), // Set expire time to 4 hrs
	})

	// Set CSRF Prevention Token
	c.Cookie(&fiber.Cookie{
		Name:    config.CSRF_TOKEN,
		Value:   ticket.Token.CSRFPreventionToken,
		Expires: time.Now().Add(time.Hour * 12), // Set expire time to 4 hrs
	})

	response := model.CookiesResponse{
		PVEAuthToken:        ticket.Token.Cookie,
		CSRFPreventionToken: ticket.Token.CSRFPreventionToken,
	}

	log.Printf("Finished getting ticket by user : %s", body.Username)
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": response})
}

// CreateUser - Create new user in Proxmox
// POST /api2/json/access/users
/*
	using Request's Body
	@userid
	@groups
	@expire : set default to 4 years
*/
func CreateUser(c *fiber.Ctx) error {
	// Getting request's body
	body := new(model.CreateUserBody)
	if err := c.BodyParser(body); err != nil {
		log.Println("Error: Could not parse body parser to creating user's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to creating user's body"})
	}
	userid := fmt.Sprintf("%s%s", body.UserID, config.REALM)

	// Mapping values
	data := url.Values{}
	data.Set("userid", userid)
	data.Set("password", body.Password)
	data.Set("groups", body.Groups)

	// Creating User
	log.Printf("Creating user : %s", body.UserID)
	_, createErr := access.CreateUser(data)
	if createErr != nil {
		log.Println("Error: Could not create user :", createErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed creating user : %s due to %s", body.UserID, createErr)})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Creating user %s successfully", body.UserID)})
}

// UpdateUser - Update user detail in Proxmox
// PUT /api2/json/access/users/{userid}
/*
	using Params
	@userid

	using Request's Body
	@enable : set to '0' to disable account
	@groups
*/
func UpdateUser(c *fiber.Ctx) error {
	// Getting request's body
	body := new(model.UpdateUserBody)
	if err := c.BodyParser(body); err != nil {
		log.Println("Error: Could not parse body parser to updating user's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to updating user's body"})
	}
	username := c.Params("username")
	cookies := config.GetCookies(c)
	userid := fmt.Sprintf("%s%s", username, config.REALM)
	updateURL := config.GetURL(fmt.Sprintf("/api2/json/access/users/%s", userid))

	// Mapping values
	data := url.Values{}
	data.Set("enable", body.Enable)
	data.Set("groups", body.Groups)

	// Creating User
	log.Printf("Updating user : %s", username)
	_, updateErr := access.UpdateUser(updateURL, data, cookies)
	if updateErr != nil {
		log.Println("Error: Could not update user :", updateErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed updating user : %s due to %s", username, updateErr)})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Updating user %s successfully", username)})
}

// DeleteUser - Delete user in Proxmox
// DELETE /api2/json/access/users/{userid}
/*
	using Params
	@userid
*/
func DeleteUser(c *fiber.Ctx) error {
	// Getting params from URL
	username := c.Params("username")
	cookies := config.GetCookies(c)
	userid := fmt.Sprintf("%s%s", username, config.REALM)
	deleteURL := config.GetURL(fmt.Sprintf("/api2/json/access/users/%s", userid))

	// Getting user's group
	group, getGroupErr := database.GetUserGroup(username)
	if getGroupErr != nil {
		log.Println("Error: Could not get user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting user's group due to %s", getGroupErr)})
	}

	// Deleting User in Proxmox
	log.Printf("Deleting user : %s", username)
	_, deleteErr := access.DeleteUser(deleteURL, cookies)
	if deleteErr != nil {
		log.Println("Error: Could not delete user :", deleteErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed deleting user : %s due to %s", username, deleteErr)})
	}

	// Deleting User in DB
	log.Printf("Deleting user : %s", username)
	err := database.DeleteUserDB(username, group)
	if err != nil {
		log.Println("Error: Could not delete user in DB due to :", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed deleting user : %s due to %s", username, err)})
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
