// Package handler - for handling any context
package handler

import "github.com/gofiber/fiber/v2"

// GetTicket - get ticket & CSRF prevention token from Proxmox
func GetTicket(c *fiber.Ctx) error {
	msg := "🎟️ Get Ticket"
	return c.SendString(msg)
}
