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
	userLogin := new(model.Login)
	if err := c.BodyParser(userLogin); err != nil {
		log.Println("Error: Could not parse body parser to getting ticket's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to getting ticket's body"})
	}

	// Mapping values
	data := url.Values{}
	data.Set("username", userLogin.Username)
	data.Set("password", userLogin.Password)
	data.Set("realm", "IAM-CE")

	// Getting Ticket
	log.Printf("Getting ticket from user : %s", userLogin.Username)
	getTicketURL := config.GetURL("/api2/json/access/ticket")
	ticket, ticketErr := access.GetTicket(getTicketURL, data)
	if ticketErr != nil {
		log.Println("Error: Could not get ticket :", ticketErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting ticket from user : %s due to %s", userLogin.Username, ticketErr)})
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

	log.Printf("Finished getting ticket by user : %s", userLogin.Username)
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Getting ticket from user %s successfully", userLogin.Username)})
}

// RealmSync - Syncs users and/or groups from the configured LDAP
func RealmSync(c *fiber.Ctx) error {
	cookies := config.GetCookies(c)
	log.Println("Realm sync started ...")
	err := access.RealmSync(cookies)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed realm syncing due to %s", err)})
	}
	log.Println("Realm sync finished successfully")
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": "Realm sync finished successfully"})
}
