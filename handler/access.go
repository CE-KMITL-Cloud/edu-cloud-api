// Package handler - handling context
package handler

import (
	"net/url"
	"time"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/internal"
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
	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Construct URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = "/api2/json/access/ticket"
	urlStr := u.String()

	// Getting request's body
	userLogin := new(model.Login)
	if err := c.BodyParser(userLogin); err != nil {
		return err
	}

	// Mapping values
	data := url.Values{}
	data.Set("username", userLogin.Username)
	data.Set("password", userLogin.Password)

	// Getting Ticket
	ticket, err := internal.GetTicket(urlStr, data)
	if err != nil {
		return err
	}

	// Set Cookie
	c.Cookie(&fiber.Cookie{
		Name:    "PVEAuthCookie",
		Value:   ticket.Token.Cookie,
		Expires: time.Now().Add(time.Hour * 4),
	})

	// Set CSRF Prevention Token
	c.Cookie(&fiber.Cookie{
		Name:    "CSRFPreventionToken",
		Value:   ticket.Token.CSRFPreventionToken,
		Expires: time.Now().Add(time.Hour * 4),
	})

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Got Ticket"})
}
