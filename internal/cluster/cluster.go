// Package cluster - Cluster functions
package cluster

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"regexp"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/model"
)

// AllocateNode - allocate which node is the best choice to have interaction with (e.g. cloning, creating)
// GET /api2/json/cluster/resources
func AllocateNode(spec model.VMSpec, cookies model.Cookies) ([]model.Node, string, error) {
	log.Println("Getting nodes from cluster's resources ...")

	url := config.GetURL("/api2/json/cluster/resources")
	nodeResource := model.NodeResource{}
	body, err := config.SendRequestWithErr(http.MethodGet, url, nil, cookies)
	if err != nil {
		return []model.Node{}, "", err
	}
	if marshalErr := json.Unmarshal(body, &nodeResource); marshalErr != nil {
		return []model.Node{}, "", marshalErr
	}

	// Regex and return only worker nodes
	var nodeList []model.Node
	for i := 0; i < len(nodeResource.Nodes); i++ {
		r, _ := regexp.Compile(config.WorkerNode) // match node which start with work-{number}
		if nodeResource.Nodes[i].Type == "node" && r.MatchString(nodeResource.Nodes[i].Node) {
			nodeList = append(nodeList, nodeResource.Nodes[i])
		}
	}

	// Compare all of them which node is the best {mem, cpu}
	var selectedNode model.Node
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

// GetStorageList - Getting RBD storage list
// GET /api2/json/cluster/resources
// func GetStorageList(cookies model.Cookies) ([]model.Storage, error) {
// 	log.Println("Getting storages from cluster's resources ...")

// 	url := config.GetURL("/api2/json/cluster/resources")
// 	storageResource := model.StorageResource{}
// 	body, err := config.SendRequestWithErr(http.MethodGet, url, nil, cookies)
// 	if err != nil {
// 		return []model.Storage{}, err
// 	}
// 	if marshalErr := json.Unmarshal(body, &storageResource); marshalErr != nil {
// 		return []model.Storage{}, marshalErr
// 	}

// 	// Regex and return only worker nodes
// 	var storageList []model.Storage
// 	for i := 0; i < len(storageResource.Storages); i++ {
// 		if storageResource.Storages[i].Type == "storage" && storageResource.Storages[i].PluginType == "rbd" {
// 			storageList = append(storageList, storageResource.Storages[i])
// 		}
// 	}
// 	log.Println(storageList)

// 	return storageList, errors.New("error: Node have no enough free space")
// }
