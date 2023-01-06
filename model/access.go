// Package model - structs
package model

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Ticket - struct for getting access token
type Ticket struct {
	Token Token `json:"data"`
}

// Token - struct inside Ticket's data field
type Token struct {
	Username            string `json:"username"`
	Cookie              string `json:"ticket"`
	CSRFPreventionToken string `json:"CSRFPreventionToken"`
}

// Login - struct for authentication to Proxmox
type Login struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

// Cookies - struct for parsing Cookies
type Cookies struct {
	Cookie              http.Cookie
	CSRFPreventionToken fiber.Cookie
}
