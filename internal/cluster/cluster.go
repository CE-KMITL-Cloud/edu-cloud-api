// Package cluster - Cluster functions
package cluster

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"regexp"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/model"
)

// AllocateNode - allocate which node is the best choice to have interaction with (e.g. cloning, creating)
// GET /api2/json/cluster/resources
// TODO : Parse memory, cpu, disk for checking that is there enough resource for create, clone or edit or not
func AllocateNode(spec model.VMSpec, cookies model.Cookies) ([]model.Resources, string, error) {
	// Getting all nodes in cluster
	log.Println("Getting cluster's resources ...")

	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Construct URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = "/api2/json/cluster/resources"
	urlStr := u.String()

	// check status of the vm
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		log.Println(err)
	}

	// Getting cookie
	req.AddCookie(&cookies.Cookie)
	req.Header.Add(config.CSRF_TOKEN, cookies.CSRFPreventionToken.Value)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	// If not http.StatusOK then log error
	if resp.StatusCode != http.StatusOK {
		log.Println("Error: with status", resp.Status)
		return []model.Resources{}, "", errors.New(resp.Status)
	}

	// Parsing response
	node := model.Node{}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Println(readErr)
	}
	// log.Println(string(body))

	// Unmarshal body to struct
	if marshalErr := json.Unmarshal(body, &node); marshalErr != nil {
		log.Println(marshalErr)
	}

	var nodeList []model.Resources

	// Regex and return only worker nodes
	for i := 0; i < len(node.Resources); i++ {
		r, _ := regexp.Compile(config.WorkerNode) // match node which start with work-{number}
		if node.Resources[i].Type == "node" && r.MatchString(node.Resources[i].Node) {
			nodeList = append(nodeList, node.Resources[i])
		}
	}
	// log.Println(nodeList)

	// Compare all of them which node is the best {mem, cpu}
	var selectedNode model.Resources
	var maxFreeMemory, maxFreeMemoryPercent uint64
	var maxFreeDisk uint64
	var maxFreeCPU, maxFreeCPUPercent float64 = -math.MaxFloat64, -math.MaxFloat64
	for _, node := range nodeList {
		freeMemory, freeCPU, freeDisk := node.MaxMem-node.Mem, node.MaxCPU-node.CPU, node.MaxDisk-node.Disk
		freeMemoryPercent := (freeMemory / node.MaxMem) * 100
		freeCPUPercent := (freeCPU / node.MaxCPU) * 100
		log.Printf("node: %s, free mem: %d, free cpu: %f, free disk: %d", node.Node, freeMemory, freeCPU, freeDisk)
		if freeMemoryPercent > maxFreeMemoryPercent || (freeMemoryPercent == maxFreeMemoryPercent && freeCPUPercent > maxFreeCPUPercent) {
			maxFreeMemoryPercent, maxFreeCPUPercent = freeMemoryPercent, freeCPUPercent
			maxFreeMemory, maxFreeCPU, maxFreeDisk = freeMemory, freeCPU, freeDisk
			selectedNode = node
		}
	}

	log.Printf("selected node: %s, free mem: %d, free cpu: %f, free disk: %d", selectedNode.Node, maxFreeMemory/config.Gigabyte, maxFreeCPU, maxFreeDisk/config.Gigabyte)
	log.Printf("vm spec mem: %d, cpu: %f, disk: %d", spec.Memory/config.Gigabyte, spec.CPU, spec.Disk/config.Gigabyte)
	if maxFreeMemory > spec.Memory && maxFreeCPU > spec.CPU && maxFreeDisk > spec.Disk {
		log.Printf("Return selected node : %s", selectedNode.Node)
		return nodeList, selectedNode.Node, nil
	}
	log.Printf("Selected node : %s have no enough free space", selectedNode.Node)
	return nodeList, selectedNode.Node, errors.New("error: Node have no enough free space")
}
