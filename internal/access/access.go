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
func GetTicket(url string, data url.Values) (model.Ticket, error) {
	// Return objects
	ticket := model.Ticket{}
	body, err := config.SendRequestWithoutCookie(http.MethodPost, url, data)
	if err != nil {
		return ticket, err
	}
	// Unmarshal body to struct
	if marshalErr := json.Unmarshal(body, &ticket); marshalErr != nil {
		return ticket, marshalErr
	}
	return ticket, nil
}
