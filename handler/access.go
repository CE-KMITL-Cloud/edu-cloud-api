// Package handler - handling context
package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/edu-cloud-api/config"
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
		Expires: time.Now().Add(time.Hour * 4), // Set expire time to 4 hrs
	})

	// Set CSRF Prevention Token
	c.Cookie(&fiber.Cookie{
		Name:    config.CSRF_TOKEN,
		Value:   ticket.Token.CSRFPreventionToken,
		Expires: time.Now().Add(time.Hour * 4), // Set expire time to 4 hrs
	})

	log.Printf("Finished getting ticket by user : %s", body.Username)
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Getting ticket from user %s successfully", body.Username)})
}

// CreateUser - Create new user
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

// UpdateUser - Update user detail
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

// DeleteUser - Delete user
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

	// Creating User
	log.Printf("Deleting user : %s", username)
	_, deleteErr := access.DeleteUser(deleteURL, cookies)
	if deleteErr != nil {
		log.Println("Error: Could not delete user :", deleteErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed deleting user : %s due to %s", username, deleteErr)})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Deleting user %s successfully", username)})
}

// RealmSync - Syncs users and/or groups from the configured LDAP
// ! Not use
// func RealmSync(c *fiber.Ctx) error {
// 	cookies := config.GetCookies(c)
// 	log.Println("Realm sync started ...")
// 	err := access.RealmSync(cookies)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed realm syncing due to %s", err)})
// 	}
// 	log.Println("Realm sync finished successfully")
// 	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": "Realm sync finished successfully"})
// }
