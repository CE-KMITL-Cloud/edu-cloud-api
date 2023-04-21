// Package qemu - QEMU functions
package qemu

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/model"
)

// GetVM - GET /api2/json/nodes/{node}/qemu/{vmid}/status/current
func GetVM(url string, cookies model.Cookies) (model.VM, error) {
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

// GetVMUsingToken - GET /api2/json/nodes/{node}/qemu/{vmid}/status/current but using api token
func GetVMUsingToken(url string) (model.VM, error) {
	info := model.VM{}
	body, err := config.SendRequestUsingToken(http.MethodGet, url, nil)
	if err != nil {
		return info, err
	}
	if marshalErr := json.Unmarshal(body, &info); marshalErr != nil {
		return info, marshalErr
	}
	return info, nil
}

// GetVMConfig - GET /api2/json/nodes/{node}/qemu/{vmid}/config
func GetVMConfig(url string, cookies model.Cookies) (model.VMConfig, error) {
	info := model.VMConfig{}
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

// DeleteVMUsingToken - DELETE /api2/json/nodes/{node}/qemu/{vmid} but using api token
func DeleteVMUsingToken(url string) error {
	_, err := config.SendRequestUsingToken(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	return nil
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
	// log.Println(templateList)
	return templateList, nil
}

// PowerManagement - POST /api2/json/nodes/{node}/qemu/{vmid}/status/{action}
/*
	action : { start, stop, suspend, shutdown, resume, reset }
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

// PowerManagementUsingToken - POST /api2/json/nodes/{node}/qemu/{vmid}/status/{action} but using api token
/*
	action : { start, stop, suspend, shutdown, resume, reset }
*/
func PowerManagementUsingToken(url string, data url.Values) (model.VMResponse, error) {
	response := model.VMResponse{}
	body, err := config.SendRequestUsingToken(http.MethodPost, url, data)
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

// VncProxy - POST /api2/json/nodes/{node}/qemu/{vmid}/vncproxy
func VncProxy(url string, data url.Values, cookies model.Cookies) (model.VncProxy, error) {
	response := model.VncProxy{}
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

// RegenerateCloudinit - PUT /api2/json/nodes/{node}/qemu/{vmid}/cloudinit
// ! to be deprecated
func RegenerateCloudinit(url string, cookies model.Cookies) (model.VMResponse, error) {
	response := model.VMResponse{}
	body, err := config.SendRequestWithErr(http.MethodPut, url, nil, cookies)
	if err != nil {
		return response, err
	}
	if marshalErr := json.Unmarshal(body, &response); marshalErr != nil {
		return response, marshalErr
	}
	return response, nil
}

// GetVMID - Getting VMID for creating, cloning from cluster's resources
func GetVMID(cookies model.Cookies) (string, error) {
	log.Println("Getting last VMID from cluster's resources ...")
	url := config.GetURL("/api2/json/cluster/resources")
	var min, max uint64
	resources := model.VMsList{}
	idList := []uint64{}
	body, err := config.SendRequestWithErr(http.MethodGet, url, nil, cookies)
	if err != nil {
		return "", err
	}
	if marshalErr := json.Unmarshal(body, &resources); marshalErr != nil {
		return "", marshalErr
	}

	// Filter recources to get only VM Template
	for i := 0; i < len(resources.VMsList); i++ {
		if resources.VMsList[i].Type == "qemu" {
			idList = append(idList, resources.VMsList[i].VMID)
			if resources.VMsList[i].VMID > max {
				max = resources.VMsList[i].VMID
			}
		}
	}

	// Find min value that not in list
	present := make([]bool, max+1)
	for _, num := range idList {
		present[num] = true
	}
	// fixed VMID must more than 100
	for j := uint64(100); j < uint64(len(present)); j++ {
		if !present[j] {
			min = j
			break
		}
	}
	log.Println("min vmid :", min)
	return fmt.Sprint(min), nil
}
