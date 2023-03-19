// Package access - Access functions
package access

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/model"
)

// GetTicket - get cookie & CSRF prevention token from Proxmox
func GetTicket(data url.Values) (model.Ticket, error) {
	ticket := model.Ticket{}
	url := config.GetURL("/api2/json/access/ticket")
	body, err := config.SendRequestWithoutCookie(http.MethodPost, url, data)
	if err != nil {
		return ticket, err
	}
	if marshalErr := json.Unmarshal(body, &ticket); marshalErr != nil {
		return ticket, marshalErr
	}
	return ticket, nil
}

// CreateUser - create user and set group in proxmox
func CreateUser(data url.Values) (string, error) {
	url := config.GetURL("/api2/json/access/users")
	body, err := config.SendRequestUsingToken(http.MethodPost, url, data)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// UpdateUser - update user in proxmox
func UpdateUser(url string, data url.Values, cookies model.Cookies) (string, error) {
	body, err := config.SendRequestWithErr(http.MethodPut, url, data, cookies)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// DeleteUser - delete user in proxmox
func DeleteUser(url string, cookies model.Cookies) (string, error) {
	body, err := config.SendRequestWithErr(http.MethodDelete, url, nil, cookies)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
