// Package qemu - QEMU functions
package qemu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/model"
)

// CheckStatus - for checking status of any process from given VM
func CheckStatus(node, vmid string, statuses []string, lock bool, timeout, sleepTime time.Duration) bool {
	log.Printf("Checking VM status on %s in %s ...", vmid, node)
	timeoutCh := time.After(timeout)
	for {
		select {
		case <-timeoutCh:
			log.Println("Timeout reached, Task not finished")
			return false
		default:
			vmStatusURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", node, vmid))
			resp, err := config.SendRequest(http.MethodGet, vmStatusURL, nil)
			if err != nil {
				log.Println(err)
			}
			defer resp.Body.Close()

			// if resp.StatusCode != http.StatusInternalServerError && resp.StatusCode != http.StatusOK {
			// 	log.Println("Error: with status", resp.Status)
			// }

			vm := model.VM{}
			body, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				log.Println(readErr)
			}
			// log.Println(string(body))
			if marshalErr := json.Unmarshal(body, &vm); marshalErr != nil {
				log.Println(marshalErr)
			}
			log.Printf("Status of %s in %s : %s", vmid, node, vm.Info.Status)
			if resp.StatusCode == http.StatusInternalServerError {
				log.Println("Error: with status", resp.Status, "or could not found VM")
			}
			// Check lock field in response. If lock field is null => unlocked
			if lock {
				log.Println("Enter checking lock")
				if vm.Info.Lock == "" || config.Contains(statuses, vm.Info.Status) {
					log.Printf("VMID : %s from %s has been unlocked or break with finished status", vmid, node)
					return true
				}
			}
			// incase status is in successful status list
			if config.Contains(statuses, vm.Info.Status) {
				log.Printf("Break status : %s", vm.Info.Status)
				return true
			}
			time.Sleep(sleepTime)
		}
	}
}

// CheckQmpStatus - for checking QMP Status of any process from specific VM
func CheckQmpStatus(node, vmid string, statuses []string, lock bool, timeout, sleepTime time.Duration) bool {
	log.Println("Checking QMP Status ...")
	timeoutCh := time.After(timeout)
	for {
		select {
		case <-timeoutCh:
			log.Println("Timeout reached, Task not finished")
			return false
		default:
			vmStatusURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", node, vmid))
			resp, err := config.SendRequest(http.MethodGet, vmStatusURL, nil)
			if err != nil {
				log.Println(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusInternalServerError && resp.StatusCode != http.StatusOK {
				log.Println("Error: with status", resp.Status)
			}

			vm := model.VM{}
			body, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				log.Println(readErr)
			}
			if marshalErr := json.Unmarshal(body, &vm); marshalErr != nil {
				log.Println(marshalErr)
			}
			log.Printf("Status of %s in %s : %s", vmid, node, vm.Info.QmpStatus)

			if resp.StatusCode == http.StatusInternalServerError {
				log.Println("Error: with status", resp.Status, "or could not found VM")
			}
			// Check lock field in response. If lock field is null => unlocked
			if lock {
				if vm.Info.Lock == "" && config.Contains(statuses, vm.Info.QmpStatus) {
					log.Printf("VMID : %s from %s has been unlocked, break with QMP Status : %s", vmid, node, vm.Info.QmpStatus)
					return true
				}
			}
			// incase status is in successful QMP Status list
			if config.Contains(statuses, vm.Info.QmpStatus) {
				log.Printf("Break QMP Status : %s", vm.Info.QmpStatus)
				return true
			}
			time.Sleep(sleepTime)
		}
	}
}

// DeleteCompletely - for assuring that status of target VM has been deleted
func DeleteCompletely(node, vmid string) bool {
	log.Println("Checking delete status ...")

	// Timeout - Default set to 1 min
	timeoutCh := time.After(time.Minute)
	for {
		select {
		case <-timeoutCh:
			log.Println("Timeout reached, Task not finished")
			return false
		default:
			vmStatusURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", node, vmid))
			resp, err := config.SendRequest(http.MethodGet, vmStatusURL, nil)
			if err != nil {
				log.Println(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				if resp.StatusCode == http.StatusInternalServerError {
					log.Printf("VMID : %s from %s is missing, Assume that VM has been deleted", vmid, node)
					return true
				}
				log.Println("Error: with status", resp.Status)
			}

			vm := model.VM{}
			body, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				log.Println(readErr)
			}
			if marshalErr := json.Unmarshal(body, &vm); marshalErr != nil {
				log.Println(marshalErr)
			}
			log.Printf("Status of %s in %s : %s", vmid, node, vm.Info.Status)

			// if status field is "deleted" return
			if vm.Info.Status == "deleted" {
				log.Printf("VMID : %s from %s has been deleted", vmid, node)
				return true
			}
			time.Sleep(time.Second)
		}
	}
}

// TemplateCompletely - for assuring that status of target VM has been templated
func TemplateCompletely(node, vmid string, statuses []string) bool {
	log.Println("Checking template status ...")

	// Timeout - Default set to 1 min
	timeoutCh := time.After(time.Minute)

	for {
		select {
		case <-timeoutCh:
			log.Println("Timeout reached, Task not finished")
			return false
		default:
			vmStatusURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", node, vmid))
			resp, err := config.SendRequest(http.MethodGet, vmStatusURL, nil)
			if err != nil {
				log.Println(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				if resp.StatusCode == http.StatusInternalServerError {
					log.Printf("Error when templating VMID : %s from %s, Assume that VM has been templated", vmid, node)
					return true
				}
				log.Println("Error: with status", resp.Status)
			}

			template := model.VMTemplate{}
			body, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				log.Println(readErr)
			}
			if marshalErr := json.Unmarshal(body, &template); marshalErr != nil {
				log.Println(marshalErr)
			}
			log.Printf("Status of %s in %s : %s", vmid, node, template.Info.Status)

			// If template = 1 : true -> templated completely
			if template.Info.Template == 1 {
				log.Printf("VMID : %s from %s has been templated", vmid, node)
				return true
			}

			// incase status is in successful status list
			if config.Contains(statuses, template.Info.Status) {
				log.Printf("Break status : %s", template.Info.Status)
				return true
			}
			time.Sleep(time.Second)
		}
	}
}

// IsTemplate - Checking VM template from VMID
func IsTemplate(node, vmid string) bool {
	log.Println("Checking VM template from VMID ...")

	// Timeout - Default set to 1 min
	timeoutCh := time.After(time.Minute)
	select {
	case <-timeoutCh:
		log.Println("Timeout reached, Task not finished")
		return false
	default:
		vmStatusURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", node, vmid))
		resp, err := config.SendRequest(http.MethodGet, vmStatusURL, nil)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Println("Error: with status", resp.Status)
			return false
		}

		template := model.VMTemplate{}
		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Println(readErr)
		}
		if marshalErr := json.Unmarshal(body, &template); marshalErr != nil {
			log.Println(marshalErr)
		}

		// If template = 1 : true -> templated completely
		if template.Info.Template == 1 {
			log.Printf("VMID : %s from %s has been templated", vmid, node)
			return true
		}
		time.Sleep(time.Second)
	}
	return false
}
