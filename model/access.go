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
	Username string `json:"username"`
	Password string `json:"password"`
}

// Cookies - struct for parsing Cookies
type Cookies struct {
	Cookie              http.Cookie
	CSRFPreventionToken fiber.Cookie
}

// CreateUserBody - struct for create user body in proxmox
type CreateUserBody struct {
	UserID   string `json:"userid"`
	Password string `json:"password"`
	Groups   string `json:"groups"`
}

// UpdateUserBody - struct for update user body in proxmox
type UpdateUserBody struct {
	Enable string `json:"enable"`
	Groups string `json:"groups"`
}
