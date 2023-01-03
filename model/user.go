// Package model - structs
package model

// DB's table
// User - struct for user's info
type User struct {
	Username            string `json:"username"`
	Password            string `json:"password"`
	GroupName           string `json:"group_name"`
	Email               string `json:"email"`
	Name                string `json:"name"`
	Surname             string `json:"surname"`
	TelMobile           string `json:"tel_mobile"`
	Status              string `json:"status"`
	CSRFPreventionToken string `json:"csrf_token"`
	CreateTime          string `json:"create_time"`
}
