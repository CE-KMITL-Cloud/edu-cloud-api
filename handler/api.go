// Package handler - handling context
package handler

import (
	"github.com/gofiber/fiber/v2"
)

// Healthy - check API status
func Healthy(c *fiber.Ctx) error {
	msg := "✋ Healthy"
	return c.SendString(msg)
}
