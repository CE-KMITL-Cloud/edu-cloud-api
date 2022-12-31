// Package handler - interface for using with router
package handler

import (
	"log"
	"net/url"
	"time"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/internal"
	"github.com/edu-cloud-api/model"
	"github.com/gofiber/fiber/v2"
)

// GetTicket - interface for GetTicket func
func GetTicket(c *fiber.Ctx) error {
	// Get host's URL
	hostURL := config.GetFromENV("proxmoxHost")

	// Construct URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = "/api2/json/access/ticket"
	urlStr := u.String()

	// Get body parser
	userLogin := new(model.Login)
	if err := c.BodyParser(userLogin); err != nil {
		return err
	}

	// Mapping values
	data := url.Values{}
	data.Set("username", userLogin.Username)
	data.Set("password", userLogin.Password)

	ticket, err := internal.GetTicket(urlStr, data)
	if err != nil {
		return err
	}

	// Set Cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "PVEAuthCookie"
	cookie.Value = ticket.Token.Cookie
	cookie.Domain = u.Hostname()
	cookie.Expires = time.Now().Add(time.Hour)
	c.Cookie(cookie)

	log.Println(cookie)

	return c.JSON(ticket)
}
