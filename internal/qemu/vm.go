// Package qemu - QEMU functions
package qemu

import (
	"encoding/json"
	"log"
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

// GetVMListByNode - GET /api2/json/nodes/{node}/qemu
func GetVMListByNode(url string, cookies model.Cookies) (model.VMList, error) {
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

// GetVMList - Getting VM list
// GET /api2/json/cluster/resources
func GetVMList(cookies model.Cookies) ([]model.VMsInfo, error) {
	// TODO: should return only user's VM
	// user := model.User{}
	// database.DB.Db.Find(&user, "username = ?", username)
	// log.Println(user)

	log.Println("Getting VM list from cluster's resources ...")
	url := config.GetURL("/api2/json/cluster/resources")
	resources := model.VMsList{}
	body, err := config.SendRequestWithErr(http.MethodGet, url, nil, cookies)
	if err != nil {
		return []model.VMsInfo{}, err
	}
	if marshalErr := json.Unmarshal(body, &resources); marshalErr != nil {
		return []model.VMsInfo{}, marshalErr
	}

	// Filter recources to get only VM (Template not included)
	var vmIDList []string
	var vmList []model.VMsInfo
	for i := 0; i < len(resources.VMsList); i++ {
		if resources.VMsList[i].Type == "qemu" && resources.VMsList[i].Template == 0 {
			if !config.Contains(vmIDList, resources.VMsList[i].ID) {
				vmIDList = append(vmIDList, resources.VMsList[i].ID)
				vmList = append(vmList, resources.VMsList[i])
			}
		}
	}
	return vmList, nil
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

// GetTemplateList - Getting VM Template list
// GET /api2/json/cluster/resources
func GetTemplateList(cookies model.Cookies) ([]model.VMsInfo, error) {
	// TODO: should return only user's VM
	// user := model.User{}
	// database.DB.Db.Find(&user, "username = ?", username)
	// log.Println(user)

	log.Println("Getting VM Template from cluster's resources ...")
	url := config.GetURL("/api2/json/cluster/resources")
	resources := model.VMsList{}
	body, err := config.SendRequestWithErr(http.MethodGet, url, nil, cookies)
	if err != nil {
		return []model.VMsInfo{}, err
	}
	if marshalErr := json.Unmarshal(body, &resources); marshalErr != nil {
		return []model.VMsInfo{}, marshalErr
	}

	// Filter recources to get only VM Template
	var idList []string
	var templateList []model.VMsInfo
	for i := 0; i < len(resources.VMsList); i++ {
		if resources.VMsList[i].Type == "qemu" && resources.VMsList[i].Template == 1 {
			if !config.Contains(idList, resources.VMsList[i].ID) {
				idList = append(idList, resources.VMsList[i].ID)
				templateList = append(templateList, resources.VMsList[i])
			}
		}
	}
	log.Println(templateList)
	return templateList, nil
}

// PowerManagement - POST /api2/json/nodes/{node}/qemu/{vmid}/status/{action}
/*
	action : { start, stop, suspend, shutdown, resume, reset }
*/
func PowerManagement(url string, data url.Values, cookies model.Cookies) (model.VMResponse, error) {
	// TODO: should return only user's VM
	// user := model.User{}
	// database.DB.Db.Find(&user, "username = ?", username)
	// log.Println(user)

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

// EditVM - POST /api2/json/nodes/{node}/qemu/{vmid}/config
func EditVM(url string, data url.Values, cookies model.Cookies) (model.VMResponse, error) {
	// TODO: should return only user's VM
	// user := model.User{}
	// database.DB.Db.Find(&user, "username = ?", username)
	// log.Println(user)

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

// ResizeDisk - PUT /api2/json/nodes/{node}/qemu/{vmid}/resize
func ResizeDisk(url string, data url.Values, cookies model.Cookies) (model.VMResponse, error) {
	// TODO: should return only user's VM
	// user := model.User{}
	// database.DB.Db.Find(&user, "username = ?", username)
	// log.Println(user)

	response := model.VMResponse{}
	body, err := config.SendRequestWithErr(http.MethodPut, url, data, cookies)
	if err != nil {
		return response, err
	}
	if marshalErr := json.Unmarshal(body, &response); marshalErr != nil {
		return response, marshalErr
	}
	return response, nil
}
