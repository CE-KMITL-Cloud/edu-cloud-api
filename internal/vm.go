// Package internal - internal functions
package internal

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/edu-cloud-api/model"
)

// GetVM - GET /api2/json/nodes/{node}/qemu/{vmid}/status/current
func GetVM(url, username string, cookies model.Cookies) (model.VM, error) {
	// TODO: should return only user's VM
	// user := model.User{}
	// database.DB.Db.Find(&user, "username = ?", username)
	// log.Println(user)

	// Return objects using string map due to returned object has many use-cases
	info := model.VM{}

	// Construct new request
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return info, err
	}

	// Getting cookie
	req.AddCookie(&cookies.Cookie)
	req.Header.Add("CSRFPreventionToken", cookies.CSRFPreventionToken.Value)

	// GET request
	resp, sendErr := client.Do(req)
	if sendErr != nil {
		return info, sendErr
	}
	defer resp.Body.Close()

	// If not 200 OK then log error
	if resp.StatusCode != 200 {
		log.Println("error: with status", resp.Status)
		return info, errors.New(resp.Status)
	}

	// We Read the response body on the line below.
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return info, readErr
	}
	log.Println(string(body))

	// Unmarshal body to struct
	if marshalErr := json.Unmarshal(body, &info); marshalErr != nil {
		return info, marshalErr
	}
	return info, nil
}

// GetVMList - GET /api2/json/nodes/{node}/qemu
func GetVMList(url, username string, cookies model.Cookies) (model.VMList, error) {
	// TODO: should return only user's VM
	// user := model.User{}
	// database.DB.Db.Find(&user, "username = ?", username)
	// log.Println(user)

	// Return objects
	info := model.VMList{}

	// Construct new request
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return info, err
	}

	// Getting cookie
	req.AddCookie(&cookies.Cookie)
	req.Header.Add("CSRFPreventionToken", cookies.CSRFPreventionToken.Value)

	// GET request
	resp, sendErr := client.Do(req)
	if sendErr != nil {
		return info, sendErr
	}
	defer resp.Body.Close()

	// If not 200 OK then log error
	if resp.StatusCode != 200 {
		log.Println("error: with status", resp.Status)
		return info, errors.New(resp.Status)
	}

	// We Read the response body on the line below.
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return info, readErr
	}

	// Unmarshal body to struct
	if marshalErr := json.Unmarshal(body, &info); marshalErr != nil {
		return info, marshalErr
	}
	return info, nil
}

// CreateVM - POST /api2/json/nodes/{node}/qemu
func CreateVM(url string, data url.Values, cookies model.Cookies) (model.VMResponse, error) {
	// Return objects
	response := model.VMResponse{}

	// Construct new request
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return response, err
	}

	// Getting cookie
	req.AddCookie(&cookies.Cookie)
	req.Header.Add("CSRFPreventionToken", cookies.CSRFPreventionToken.Value)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// POST request
	resp, sendErr := client.Do(req)
	if sendErr != nil {
		return response, sendErr
	}
	defer resp.Body.Close()

	// If not 200 OK then log error
	if resp.StatusCode != 200 {
		log.Println("error: with status", resp.Status)
		return response, errors.New(resp.Status)
	}

	// Read byte from body
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return response, readErr
	}

	// Unmarshal body to struct
	if marshalErr := json.Unmarshal(body, &response); marshalErr != nil {
		return response, marshalErr
	}
	return response, nil
}

// DeleteVM - DELETE /api2/json/nodes/{node}/qemu/{vmid}
func DeleteVM(url string, cookies model.Cookies) (model.VMResponse, error) {
	// TODO: should able to delete only user's VM
	// user := model.User{}
	// database.DB.Db.Find(&user, "username = ?", username)
	// log.Println(user)

	// Return objects
	response := model.VMResponse{}

	// Construct new request
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return response, err
	}

	// Getting cookie
	req.AddCookie(&cookies.Cookie)
	req.Header.Add("CSRFPreventionToken", cookies.CSRFPreventionToken.Value)

	// DELETE request
	resp, sendErr := client.Do(req)
	if sendErr != nil {
		return response, sendErr
	}
	defer resp.Body.Close()

	// If not 200 OK then log error
	if resp.StatusCode != 200 {
		log.Println("error: with status", resp.Status)
		return response, errors.New(resp.Status)
	}

	// Read byte from body
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return response, readErr
	}

	// Unmarshal body to struct
	if marshalErr := json.Unmarshal(body, &response); marshalErr != nil {
		return response, marshalErr
	}
	return response, nil
}

// CloneVM - POST /api2/json/nodes/{node}/qemu/{vmid}/clone
func CloneVM(url string, data url.Values, cookies model.Cookies) (model.VMResponse, error) {
	// Return objects
	response := model.VMResponse{}

	// Construct new request
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return response, err
	}

	// Getting cookie
	req.AddCookie(&cookies.Cookie)
	req.Header.Add("CSRFPreventionToken", cookies.CSRFPreventionToken.Value)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// POST request
	resp, sendErr := client.Do(req)
	if sendErr != nil {
		return response, sendErr
	}
	defer resp.Body.Close()

	// If not 200 OK then log error
	if resp.StatusCode != 200 {
		log.Println("error: with status", resp.Status)
		return response, errors.New(resp.Status)
	}

	// Read byte from body
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return response, readErr
	}

	// Unmarshal body to struct
	if marshalErr := json.Unmarshal(body, &response); marshalErr != nil {
		return response, marshalErr
	}
	return response, nil
}

// CreateTemplate - POST /api2/json/nodes/{node}/qemu/{vmid}/template
func CreateTemplate(url string, cookies model.Cookies) (model.VMResponse, error) {
	// Return objects
	response := model.VMResponse{}

	// Construct new request
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return response, err
	}

	// Getting cookie
	req.AddCookie(&cookies.Cookie)
	req.Header.Add("CSRFPreventionToken", cookies.CSRFPreventionToken.Value)

	// POST request
	resp, sendErr := client.Do(req)
	if sendErr != nil {
		return response, sendErr
	}
	defer resp.Body.Close()

	// If not 200 OK then log error
	if resp.StatusCode != 200 {
		log.Println("error: with status", resp.Status)
		return response, errors.New(resp.Status)
	}

	// Read byte from body
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return response, readErr
	}

	// Unmarshal body to struct
	if marshalErr := json.Unmarshal(body, &response); marshalErr != nil {
		return response, marshalErr
	}
	return response, nil
}
