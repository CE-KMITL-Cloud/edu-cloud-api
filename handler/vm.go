// Package handler - handling context
package handler

import (
	"net/url"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/internal"
	"github.com/edu-cloud-api/model"
	"github.com/gofiber/fiber/v2"
)

func GetVM(c *fiber.Ctx) error {
	// Get host's URL
	hostURL := config.GetFromENV("proxmoxHost")

	// Getting request's body
	body := new(model.VM)
	if err := c.BodyParser(body); err != nil {
		return err
	}

	// Construct URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = "/api2/json/access/ticket"
	urlStr := u.String()

	// Getting Ticket
	ticket, err := internal.GetVM(urlStr)
	if err != nil {
		return err
	}

	return c.JSON(ticket)
}
