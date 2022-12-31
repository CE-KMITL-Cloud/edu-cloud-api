// Package model - struct
package model

// Ticket - struct for getting access token
type Ticket struct {
	Token Token `json:"data"`
}

// Token - struct inside Ticket's data field
type Token struct {
	Cookie    string `json:"ticket"`
	CSRFToken string `json:"CSRFPreventionToken"`
}

// Login - struct to authenticate
type Login struct {
	Username string `form:"username"`
	Password string `form:"password"`
}
