// Package model - structs
package model

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

// Login - struct for authentication
type Login struct {
	Username string `form:"username"`
	Password string `form:"password"`
}
