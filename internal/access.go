// Package internal - internal functions
package internal

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/edu-cloud-api/model"
)

// GetTicket - get cookie & CSRF prevention token from Proxmox
func GetTicket(url string, data url.Values) (model.Ticket, error) {
	// Return objects
	ticket := model.Ticket{}

	// Construct new request
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// POST request
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	// If not http.StatusOK then log error
	if resp.StatusCode != 200 {
		log.Println("Error: with status", resp.Status)
	}

	// Read byte from body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	// Unmarshal body to struct
	if marshalErr := json.Unmarshal(body, &ticket); marshalErr != nil {
		return ticket, marshalErr
	}
	return ticket, nil
}
