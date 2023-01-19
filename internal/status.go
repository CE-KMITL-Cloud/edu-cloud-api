// Package internal - internal functions
package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/model"
)

// StatusVM - for checking status of any process from specific VM
func StatusVM(node, vmid string, statuses []string, wg *sync.WaitGroup, cookies model.Cookies) {
	defer wg.Done()
	log.Println("Checking status...")

	// Timeout - Default set to 30 mins
	timeoutCh := time.After(30 * time.Minute)

	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Construct URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", node, vmid)
	urlStr := u.String()

	for {
		select {
		case <-timeoutCh:
			log.Println("Timeout reached, Task not finished")
			return
		default:
			// check status of the vm
			client := &http.Client{}
			req, err := http.NewRequest(http.MethodGet, urlStr, nil)
			if err != nil {
				log.Println(err)
			}

			// Getting cookie
			req.AddCookie(&cookies.Cookie)
			req.Header.Add("CSRFPreventionToken", cookies.CSRFPreventionToken.Value)

			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
			}
			defer resp.Body.Close()

			// If not 200 OK then log error
			if resp.StatusCode != 200 {
				log.Println("error: with status", resp.Status)
			}

			// Parsing response
			info := model.VM{}
			body, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				log.Println(readErr)
			}
			log.Println(string(body))

			// Unmarshal body to struct
			if marshalErr := json.Unmarshal(body, &info); marshalErr != nil {
				log.Println(marshalErr)
			}

			// logging status of target VM
			log.Printf("Status of %s in %s : %s", vmid, node, info.Info.Status)

			// if lock field is null => unlocked
			if info.Info.Lock == "" {
				log.Printf("VMID : %s from %s has been unlocked", vmid, node)
				return
			}

			// incase status is in successful status list
			if config.Contains(statuses, info.Info.Status) {
				log.Printf("Break status : %s", info.Info.Status)
				return
			}

			// Default setting : check every 15 sec
			time.Sleep(15 * time.Second)
		}
	}
}

// DeleteCompletely - for assuring that status of target VM has been deleted
func DeleteCompletely(node, vmid string, wg *sync.WaitGroup, cookies model.Cookies) {
	defer wg.Done()
	log.Println("Checking delete status...")

	// Timeout - Default set to 10 mins
	timeoutCh := time.After(10 * time.Minute)

	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Construct URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", node, vmid)
	urlStr := u.String()

	for {
		select {
		case <-timeoutCh:
			log.Println("Timeout reached, Task not finished")
			return
		default:
			// check status of the vm
			client := &http.Client{}
			req, err := http.NewRequest(http.MethodGet, urlStr, nil)
			if err != nil {
				log.Println(err)
			}

			// Getting cookie
			req.AddCookie(&cookies.Cookie)
			req.Header.Add("CSRFPreventionToken", cookies.CSRFPreventionToken.Value)

			resp, sendErr := client.Do(req)
			if sendErr != nil {
				log.Println(sendErr)
			}
			defer resp.Body.Close()

			// If not 200 OK then log error
			if resp.StatusCode != 200 {
				// TODO : Another work around on this?
				if resp.StatusCode == 500 {
					log.Printf("VMID : %s from %s is missing, Assume that VM has been deleted", vmid, node)
					return
				}
				log.Println("error: with status", resp.Status)
			}

			// Parsing response
			info := model.VM{}
			body, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				log.Println(readErr)
			}
			log.Println(string(body))

			// Unmarshal body to struct
			if marshalErr := json.Unmarshal(body, &info); marshalErr != nil {
				log.Println(marshalErr)
			}

			// logging status of target VM
			log.Printf("Status of %s in %s : %s", vmid, node, info.Info.Status)

			// if status field is "deleted" return
			if info.Info.Status == "deleted" {
				log.Printf("VMID : %s from %s has been deleted", vmid, node)
				return
			}

			// Default setting : check every 5 sec
			time.Sleep(5 * time.Second)
		}
	}
}

// TemplateCompletely - for assuring that status of target VM has been templated
func TemplateCompletely(node, vmid string, statuses []string, wg *sync.WaitGroup, cookies model.Cookies) {
	defer wg.Done()
	log.Println("Checking template status...")

	// Timeout - Default set to 10 mins
	timeoutCh := time.After(10 * time.Minute)

	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Construct URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", node, vmid)
	urlStr := u.String()

	for {
		select {
		case <-timeoutCh:
			log.Println("Timeout reached, Task not finished")
			return
		default:
			// check status of the vm
			client := &http.Client{}
			req, err := http.NewRequest(http.MethodGet, urlStr, nil)
			if err != nil {
				log.Println(err)
			}

			// Getting cookie
			req.AddCookie(&cookies.Cookie)
			req.Header.Add("CSRFPreventionToken", cookies.CSRFPreventionToken.Value)

			resp, sendErr := client.Do(req)
			if sendErr != nil {
				log.Println(sendErr)
			}
			defer resp.Body.Close()

			// If not 200 OK then log error
			if resp.StatusCode != 200 {
				// TODO : Another work around on this?
				if resp.StatusCode == 500 {
					log.Printf("Error when templating VMID : %s from %s, Assume that VM has been templated", vmid, node)
					return
				}
				log.Println("error: with status", resp.Status)
			}

			// Parsing response
			info := model.VMTemplate{}
			body, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				log.Println(readErr)
			}
			log.Println(string(body))

			// Unmarshal body to struct
			if marshalErr := json.Unmarshal(body, &info); marshalErr != nil {
				log.Println(marshalErr)
			}

			// logging status of target VM
			log.Printf("Status of %s in %s : %s", vmid, node, info.Info.Status)

			// If template = 1 : true -> templated completely
			if info.Info.Template == 1 {
				log.Printf("VMID : %s from %s has been templated", vmid, node)
				return
			}

			// incase status is in successful status list
			if config.Contains(statuses, info.Info.Status) {
				log.Printf("Break status : %s", info.Info.Status)
				return
			}

			// Default setting : check every 5 sec
			time.Sleep(5 * time.Second)
		}
	}
}
