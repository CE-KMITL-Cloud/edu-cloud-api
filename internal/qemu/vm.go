// Package qemu - QEMU functions
package qemu

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/model"
)

// GetVM - GET /api2/json/nodes/{node}/qemu/{vmid}/status/current
func GetVM(url string, cookies model.Cookies) (model.VM, error) {
	// TODO: should return only user's VM
	// user := model.User{}
	// database.DB.Db.Find(&user, "username = ?", username)
	// log.Println(user)

	info := model.VM{}
	body, err := config.SendRequestWithErr(http.MethodGet, url, nil, cookies)
	if err != nil {
		return info, err
	}
	if marshalErr := json.Unmarshal(body, &info); marshalErr != nil {
		return info, marshalErr
	}
	return info, nil
}

// GetVMList - GET /api2/json/nodes/{node}/qemu
func GetVMList(url string, cookies model.Cookies) (model.VMList, error) {
	// TODO: should return only user's VM
	// user := model.User{}
	// database.DB.Db.Find(&user, "username = ?", username)
	// log.Println(user)

	info := model.VMList{}
	body, err := config.SendRequestWithErr(http.MethodGet, url, nil, cookies)
	if err != nil {
		return info, err
	}
	if marshalErr := json.Unmarshal(body, &info); marshalErr != nil {
		return info, marshalErr
	}
	return info, nil
}

// CreateVM - POST /api2/json/nodes/{node}/qemu
func CreateVM(url string, data url.Values, cookies model.Cookies) (model.VMResponse, error) {
	response := model.VMResponse{}
	body, err := config.SendRequestWithErr(http.MethodPost, url, data, cookies)
	if err != nil {
		return response, err
	}
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

	response := model.VMResponse{}
	body, err := config.SendRequestWithErr(http.MethodDelete, url, nil, cookies)
	if err != nil {
		return response, err
	}
	if marshalErr := json.Unmarshal(body, &response); marshalErr != nil {
		return response, marshalErr
	}
	return response, nil
}

// CloneVM - POST /api2/json/nodes/{node}/qemu/{vmid}/clone
func CloneVM(url string, data url.Values, cookies model.Cookies) (model.VMResponse, error) {
	response := model.VMResponse{}
	body, err := config.SendRequestWithErr(http.MethodPost, url, data, cookies)
	if err != nil {
		return response, err
	}
	if marshalErr := json.Unmarshal(body, &response); marshalErr != nil {
		return response, marshalErr
	}
	return response, nil
}

// CreateTemplate - POST /api2/json/nodes/{node}/qemu/{vmid}/template
func CreateTemplate(url string, cookies model.Cookies) (model.VMResponse, error) {
	response := model.VMResponse{}
	body, err := config.SendRequestWithErr(http.MethodPost, url, nil, cookies)
	if err != nil {
		return response, err
	}
	if marshalErr := json.Unmarshal(body, &response); marshalErr != nil {
		return response, marshalErr
	}
	return response, nil
}

// PowerManagement - POST /api2/json/nodes/{node}/qemu/{vmid}/status/{action}
/*
	action : { start, stop, suspend, shutdown, resume, reboot, reset }
*/
func PowerManagement(url string, data url.Values, cookies model.Cookies) (model.VMResponse, error) {
	response := model.VMResponse{}
	body, err := config.SendRequestWithErr(http.MethodPost, url, data, cookies)
	if err != nil {
		return response, err
	}
	if marshalErr := json.Unmarshal(body, &response); marshalErr != nil {
		return response, marshalErr
	}
	return response, nil
}
