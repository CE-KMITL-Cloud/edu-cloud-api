// Package access - Access functions
package access

import (
	"encoding/json"
	"log"
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

// RealmSync - Syncs users and/or groups from the configured LDAP
// POST /api2/json/access/domains/{realm}/sync
func RealmSync(cookies model.Cookies) error {
	data := url.Values{}
	data.Set("scope", "both")
	url := config.GetURL("/api2/json/access/domains/IAM-CE/sync")
	info, err := config.SendRequestWithErr(http.MethodPost, url, data, cookies)
	if err != nil {
		return err
	}
	log.Println(string(info))
	return nil
}
